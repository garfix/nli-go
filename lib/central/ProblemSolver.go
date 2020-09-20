package central

import (
	"fmt"
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
)

// The problem solver takes a relation set and a set of bindings
// and returns a set of new bindings
// It uses knowledge bases to find these bindings
type ProblemSolver struct {
	knowledgeBases		 []knowledge.KnowledgeBase
	factBases            []knowledge.FactBase
	ruleBases            []knowledge.RuleBase
	functionBases        []knowledge.FunctionBase
	aggregateBases       []knowledge.AggregateBase
	nestedStructureBases []knowledge.NestedStructureBase
	matcher              *mentalese.RelationMatcher
	modifier             *FactBaseModifier
	dialogContext        *DialogContext
	log                  *common.SystemLog
	SolveDepth           int
}

func NewProblemSolver(matcher *mentalese.RelationMatcher, dialogContext *DialogContext, log *common.SystemLog) *ProblemSolver {
	return &ProblemSolver{
		knowledgeBases: []knowledge.KnowledgeBase{},
		factBases:      []knowledge.FactBase{},
		ruleBases:      []knowledge.RuleBase{},
		functionBases:  []knowledge.FunctionBase{},
		aggregateBases: []knowledge.AggregateBase{},
		matcher:        matcher,
		modifier:       NewFactBaseModifier(log),
		dialogContext:  dialogContext,
		log:            log,
	}
}

func (solver *ProblemSolver) AddFactBase(factBase knowledge.FactBase) {
	solver.factBases = append(solver.factBases, factBase)
	solver.knowledgeBases = append(solver.knowledgeBases, factBase)
}

func (solver *ProblemSolver) AddFunctionBase(functionBase knowledge.FunctionBase) {
	solver.functionBases = append(solver.functionBases, functionBase)
	solver.knowledgeBases = append(solver.knowledgeBases, functionBase)
}

func (solver *ProblemSolver) AddRuleBase(ruleBase knowledge.RuleBase) {
	solver.ruleBases = append(solver.ruleBases, ruleBase)
	solver.knowledgeBases = append(solver.knowledgeBases, ruleBase)
}

func (solver *ProblemSolver) AddMultipleBindingsBase(source knowledge.AggregateBase) {
	solver.aggregateBases = append(solver.aggregateBases, source)
	solver.knowledgeBases = append(solver.knowledgeBases, source)
}

func (solver *ProblemSolver) AddNestedStructureBase(base knowledge.NestedStructureBase) {
	solver.nestedStructureBases = append(solver.nestedStructureBases, base)
	solver.knowledgeBases = append(solver.knowledgeBases, base)
}

// set e.g. [ father(X, Y) father(Y, Z) ]
// bindings e.g. [{X: john, Z: jack} {}]
// return e.g. [
//  { X: john, Z: jack, Y: billy }
//  { X: john, Z: jack, Y: bob }
// ]
func (solver ProblemSolver) SolveRelationSet(set mentalese.RelationSet, bindings mentalese.Bindings) mentalese.Bindings {

	solver.log.StartProduction("Solve Set", set.String() + " " + bindings.String())

	for _, relation := range set {
		if !solver.isPredicateSupported(relation.Predicate) {
			solver.log.AddError("Predicate not supported by any knowledge base: " + relation.Predicate)
			return mentalese.Bindings{}
		}
	}

	newBindings := bindings
	for _, relation := range set {
		newBindings = solver.solveSingleRelationMultipleBindings(relation, newBindings)

		if len(newBindings) == 0 {
			break
		}
	}

	// remove duplicates because they cause unnecessary work and they cause problems for the generator
	newBindings = newBindings.UniqueBindings()

	solver.log.EndProduction("Solve Set", newBindings.String())

	return newBindings
}

func (solver ProblemSolver) isPredicateSupported(predicate string) bool {
	for _, knowledgeBase := range solver.knowledgeBases {
		if knowledgeBase.HandlesPredicate(predicate) {
			return true
		}
	}
	return false
}

// goal e.g. father(Y, Z)
// bindings e.g. {
//  { {X='john', Y='jack'} }
//  { {X='bob', Y='jonathan'} }
// }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) solveSingleRelationMultipleBindings(relation mentalese.Relation, bindings []mentalese.Binding) []mentalese.Binding {

	solver.log.StartProduction("Solve Relation", relation.String() + " " + fmt.Sprint(bindings))

	newBindings := mentalese.Bindings{}
	multiFound := false
	aggregateBindings := mentalese.Bindings{}

	// Note: aggregate base relations are currently the only ones whose bindings are not limited to the variables of the arguments
	// As long as these relations are simple, this is not a problem.
	for _, aggregateBase := range solver.aggregateBases {
		aggregateBindings, multiFound = aggregateBase.Execute(relation, bindings)
		if multiFound {
			newBindings = aggregateBindings
			break
		}
	}

	if !multiFound {

		if len(bindings) == 0 {
			newBindings = solver.solveSingleRelationSingleBinding(relation, mentalese.Binding{})
		} else {
			for _, binding := range bindings {
				newBindings = append(newBindings, solver.solveSingleRelationSingleBinding(relation, binding)...)
			}
		}
	}

	solver.log.EndProduction("Solve Relation", fmt.Sprint(newBindings))

	return newBindings
}

// goalRelation e.g. father(Y, Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) solveSingleRelationSingleBinding(relation mentalese.Relation, binding mentalese.Binding) []mentalese.Binding {

	relationVariables := relation.GetVariableNames()
	simpleBinding := binding.FilterVariablesByName(relationVariables)

	solver.log.StartProduction("Solve Simple Binding", relation.String() + " " + fmt.Sprint(simpleBinding))

	newBindings := mentalese.Bindings{ }

	// go through all fact bases
	for _, factBase := range solver.factBases {
		newBindings = append(newBindings, solver.FindFacts(factBase, relation, simpleBinding)...)
	}

	// go through all rule bases
	for _, ruleBase := range solver.ruleBases {
		newBindings = append(newBindings, solver.solveSingleRelationSingleBindingSingleRuleBase(relation, simpleBinding, ruleBase)...)
	}

	// go through all function bases
	for _, functionBase := range solver.functionBases {
		resultBinding, functionFound := functionBase.Execute(relation, simpleBinding)
		if functionFound && resultBinding != nil {
			newBindings = append(newBindings, resultBinding)
		}
	}

	// go through all nested structure bases
	for _, nestedStructureBase := range solver.nestedStructureBases {
		newBindings = append(newBindings, nestedStructureBase.SolveNestedStructure(relation, simpleBinding)...)
	}

	// do assert / retract
	newBindings = append(newBindings, solver.modifyKnowledgeBase(relation, simpleBinding)...)

	solver.log.EndProduction("Solve Simple Binding", fmt.Sprint(newBindings))

	// compose the result set
	completedBindings := mentalese.Bindings{}
	for _, newBinding := range newBindings {
		// remove temporary variables
		essentialResultBinding := newBinding.FilterVariablesByName(relationVariables)
		// combine the source binding with the clean results
		completedBinding := binding.Merge(essentialResultBinding)
		completedBindings = append(completedBindings, completedBinding)
	}

	return completedBindings
}

// Creates bindings for the free variables in 'relations', by resolving them in factBase
func (solver ProblemSolver) FindFacts(factBase knowledge.FactBase, relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	dbBindings := mentalese.Bindings{}

	for _, ds2db := range factBase.GetReadMappings() {

		activeBinding, match := solver.matcher.MatchTwoRelations(relation, ds2db.Goal, mentalese.Binding{})
		if !match { continue }

		activeBinding2, match2 := solver.matcher.MatchTwoRelations(ds2db.Goal, relation, mentalese.Binding{})
		if !match2 { continue }

		dbRelations := ds2db.Pattern.ImportBinding(activeBinding2)

		localIdBinding := solver.replaceSharedIdsByLocalIds(binding, factBase)

		relevantBinding := localIdBinding.Select(dbRelations.GetVariableNames())
		newDbBindings := solver.solveMultipleRelationSingleFactBase(dbRelations, relevantBinding, factBase)

		for _, newDbBinding := range newDbBindings {

			dbBinding := activeBinding.Merge(newDbBinding)

			combinedBinding := localIdBinding.Merge(dbBinding.Select(relation.GetVariableNames()))
			sharedBinding := solver.replaceLocalIdBySharedId(combinedBinding, factBase)
			dbBindings = append(dbBindings, sharedBinding)
		}
	}

	return dbBindings
}

func (solver ProblemSolver) solveMultipleRelationSingleFactBase(relations []mentalese.Relation, binding mentalese.Binding, factBase knowledge.FactBase) mentalese.Bindings {

	sequenceBindings := mentalese.Bindings{ binding }

	for _, relation := range relations {
		sequenceBindings = solver.solveSingleRelationSingleFactBase(relation, sequenceBindings, factBase)
	}

	return sequenceBindings
}

func (solver ProblemSolver) solveSingleRelationSingleFactBase(relation mentalese.Relation, bindings mentalese.Bindings, factBase knowledge.FactBase) mentalese.Bindings {

	solver.log.StartProduction("Database" + " " + factBase.GetName(), relation.String() + " " + bindings.String())

	relationBindings := mentalese.Bindings{}

	multiFound := false
	aggregateBindings := mentalese.Bindings{}

	for _, aggregateBase := range solver.aggregateBases {
		aggregateBindings, multiFound = aggregateBase.Execute(relation, bindings)
		if multiFound {
			relationBindings = aggregateBindings
			break
		}
	}

	if !multiFound {

		for _, binding := range bindings {

			resultBindings := factBase.MatchRelationToDatabase(relation, binding)

			// found bindings must be extended with the bindings already present
			for _, resultBinding := range resultBindings {
				newRelationBinding := binding.Merge(resultBinding)
				relationBindings = append(relationBindings, newRelationBinding)
			}
		}
	}

	solver.log.EndProduction("Database" + " " + factBase.GetName(), relationBindings.String())

	return relationBindings
}

func (solver ProblemSolver) replaceSharedIdsByLocalIds(binding mentalese.Binding, factBase knowledge.FactBase) mentalese.Binding {

	newBinding := mentalese.Binding{}

	for key, value := range binding {
		newValue := value

		if value.IsId() {
			sharedId := value.TermValue
			entityType := value.TermEntityType
			if entityType != "" {
				localId := factBase.GetLocalId(sharedId, entityType)
				if localId == "" {
					solver.log.AddError(fmt.Sprintf("Local id %s not found for %s in fact base %s", sharedId, entityType, factBase.GetName()))
					return mentalese.Binding{}
				}
				newValue = mentalese.NewTermId(localId, entityType)
			}
		}

		newBinding[key] = newValue
	}

	return newBinding
}

func (solver ProblemSolver) replaceLocalIdBySharedId(binding mentalese.Binding, factBase knowledge.FactBase) mentalese.Binding {

	newBinding := mentalese.Binding{}

	for key, value := range binding {
		newValue := value

		if value.IsId() {
			localId := value.TermValue
			entityType := value.TermEntityType
			if entityType != "" {
				sharedId := factBase.GetSharedId(localId, entityType)
				if sharedId == "" {
					solver.log.AddError(fmt.Sprintf("Shared id %s not found for %s in fact base %s", localId, entityType, factBase.GetName()))
					return mentalese.Binding{}
				}
				newValue = mentalese.NewTermId(sharedId, entityType)
			}
		}

		newBinding[key] = newValue
	}

	return newBinding
}

func (solver ProblemSolver) modifyKnowledgeBase(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	newBindings := mentalese.Bindings{ }

	if len(relation.Arguments) == 0 { return newBindings }

	argument := relation.Arguments[0]

	if relation.Predicate == mentalese.PredicateAssert {
		if argument.IsRelationSet() {
			for _, factBase := range solver.factBases {
				localIdBinding := solver.replaceSharedIdsByLocalIds(binding, factBase)
				boundRelation := relation.BindSingle(localIdBinding)
				singleRelation := boundRelation.Arguments[0].TermValueRelationSet[0]
				if (singleRelation.IsBound()) {
					solver.modifier.Assert(singleRelation, factBase)
				} else {
					solver.log.AddError("Cannot assert unbound relation " + singleRelation.String())
					return mentalese.Bindings{}
				}
				binding = solver.replaceLocalIdBySharedId(binding, factBase)
				newBindings = append(newBindings, binding)
			}
		} else if argument.IsRule() {
			for _, ruleBase := range solver.ruleBases {
				rule := relation.Arguments[0].TermValueRule.BindSingle(binding)
				ruleBase.Assert(rule)
				newBindings = append(newBindings, binding)
				//  only add the rule to a single rulebase
				break
			}
		} else if argument.IsList() {
			panic("assert not implemented for list")
		}
	} else if relation.Predicate == mentalese.PredicateRetract {
		if argument.IsRelationSet() {
			for _, factBase := range solver.factBases {
				localIdBinding := solver.replaceSharedIdsByLocalIds(binding, factBase)
				boundRelation := relation.BindSingle(localIdBinding)
				solver.modifier.Retract(boundRelation.Arguments[0].TermValueRelationSet[0], factBase)
				binding = solver.replaceLocalIdBySharedId(binding, factBase)
				newBindings = append(newBindings, binding)
			}
		}
	}

	return newBindings
}

// goalRelation e.g. father('jack', Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) solveSingleRelationSingleBindingSingleRuleBase(goalRelation mentalese.Relation, binding mentalese.Binding, ruleBase knowledge.RuleBase) mentalese.Bindings {

	inputVariables := goalRelation.GetVariableNames()

	goalBindings := mentalese.Bindings{}

	// match rules from the rule base to the goalRelation
	sourceSubgoalSets, _ := ruleBase.Bind(goalRelation, binding)

	for _, sourceSubgoalSet := range sourceSubgoalSets {

		subgoalResultBindings := mentalese.Bindings{binding}

		for _, subGoal := range sourceSubgoalSet {

			subgoalResultBindings = solver.SolveRelationSet([]mentalese.Relation{subGoal}, subgoalResultBindings)
			if len(subgoalResultBindings) == 0 {
				break
			}
		}

		for _, subgoalResultBinding := range subgoalResultBindings {

			// filter out the input variables
			filteredBinding := subgoalResultBinding.FilterVariablesByName(inputVariables)

			// make sure all variables of the original binding are present
			goalBinding := binding.Merge(filteredBinding)

			goalBindings = append(goalBindings, goalBinding)
		}
	}

	return goalBindings
}
