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
	scopeStack           *mentalese.ScopeStack
	matcher              *mentalese.RelationMatcher
	modifier             *FactBaseModifier
	dialogContext        *DialogContext
	log                  *common.SystemLog
}

func NewProblemSolver(matcher *mentalese.RelationMatcher, dialogContext *DialogContext, log *common.SystemLog) *ProblemSolver {
	return &ProblemSolver{
		knowledgeBases: []knowledge.KnowledgeBase{},
		factBases:      []knowledge.FactBase{},
		ruleBases:      []knowledge.RuleBase{},
		functionBases:  []knowledge.FunctionBase{},
		aggregateBases: []knowledge.AggregateBase{},
		scopeStack: 	mentalese.NewScopeStack(),
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

func (solver *ProblemSolver) GetCurrentScope() *mentalese.Scope {
	return solver.scopeStack.GetCurrentScope()
}

// set e.g. [ father(X, Y) father(Y, Z) ]
// bindings e.g. [{X: john, Z: jack} {}]
// return e.g. [
//  { X: john, Z: jack, Y: billy }
//  { X: john, Z: jack, Y: bob }
// ]
func (solver ProblemSolver) SolveRelationSet(set mentalese.RelationSet, bindings mentalese.BindingSet) mentalese.BindingSet {

	solver.log.StartProduction("Solve Set", set.String() + " " + bindings.String())

	for _, relation := range set {
		if !solver.isPredicateSupported(relation.Predicate) {
			solver.log.AddError("Predicate not supported by any knowledge base: " + relation.Predicate)
			return mentalese.NewBindingSet()
		}
	}

	newBindings := bindings
	for _, relation := range set {
		newBindings = solver.solveSingleRelationMultipleBindings(relation, newBindings)

		if newBindings.IsEmpty() {
			break
		}
	}

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
func (solver ProblemSolver) solveSingleRelationMultipleBindings(relation mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

	solver.log.StartProduction("Solve Relation", relation.String() + " " + fmt.Sprint(bindings))

	newBindings := mentalese.NewBindingSet()
	multiFound := false
	aggregateBindings := mentalese.NewBindingSet()

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

		if bindings.IsEmpty() {
			newBindings = solver.solveSingleRelationSingleBinding(relation, mentalese.NewScopedBinding(solver.scopeStack.GetCurrentScope()))
		} else {
			for _, binding := range bindings.GetAll() {
				newBindings.AddMultiple(solver.solveSingleRelationSingleBinding(relation, binding))
			}
		}
	}

	solver.log.EndProduction("Solve Relation", relation.String() + ": " + fmt.Sprint(newBindings))

	return newBindings
}

// goalRelation e.g. father(Y, Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) solveSingleRelationSingleBinding(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	relationVariables := relation.GetVariableNames()
	simpleBinding := binding.FilterVariablesByName(relationVariables)

	solver.log.StartProduction("Solve Simple Binding", relation.String() + " " + fmt.Sprint(simpleBinding))

	newBindings := mentalese.NewBindingSet()

	// go through all fact bases
	for _, factBase := range solver.factBases {
		newBindings.AddMultiple(solver.FindFacts(factBase, relation, simpleBinding))
	}

	// go through all rule bases
	for _, ruleBase := range solver.ruleBases {
		newBindings.AddMultiple(solver.solveSingleRelationSingleBindingSingleRuleBase(relation, simpleBinding, ruleBase))
	}

	// go through all function bases
	for _, functionBase := range solver.functionBases {
		resultBinding, functionFound, success := functionBase.Execute(relation, simpleBinding)
		if functionFound && success {
			newBindings.Add(resultBinding)
		}
	}

	// go through all nested structure bases
	for _, nestedStructureBase := range solver.nestedStructureBases {
		newBindings.AddMultiple(nestedStructureBase.SolveNestedStructure(relation, simpleBinding))
	}

	// do assert / retract
	newBindings.AddMultiple(solver.modifyKnowledgeBase(relation, simpleBinding))

	solver.log.EndProduction("Solve Simple Binding", relation.String() + ": " + fmt.Sprint(newBindings))

	// compose the result set
	completedBindings := mentalese.NewBindingSet()
	for _, newBinding := range newBindings.GetAll() {
		// remove temporary variables
		essentialResultBinding := newBinding.FilterVariablesByName(relationVariables)
		// combine the source binding with the clean results
		completedBinding := binding.Merge(essentialResultBinding)
		completedBindings.Add(completedBinding)
	}

	return completedBindings
}

// Creates bindings for the free variables in 'relations', by resolving them in factBase
func (solver ProblemSolver) FindFacts(factBase knowledge.FactBase, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	dbBindings := mentalese.NewBindingSet()

	for _, ds2db := range factBase.GetReadMappings() {

		activeBinding, match := solver.matcher.MatchTwoRelations(relation, ds2db.Goal, mentalese.NewBinding())
		if !match { continue }

		activeBinding2, match2 := solver.matcher.MatchTwoRelations(ds2db.Goal, relation, mentalese.NewBinding())
		if !match2 { continue }

		dbRelations := ds2db.Pattern.ImportBinding(activeBinding2)

		localIdBinding := solver.replaceSharedIdsByLocalIds(binding, factBase)

		relevantBinding := localIdBinding.Select(dbRelations.GetVariableNames())
		newDbBindings := solver.solveMultipleRelationSingleFactBase(dbRelations, relevantBinding, factBase)

		for _, newDbBinding := range newDbBindings.GetAll() {

			dbBinding := activeBinding.Merge(newDbBinding)

			combinedBinding := localIdBinding.Merge(dbBinding.Select(relation.GetVariableNames()))
			sharedBinding := solver.replaceLocalIdBySharedId(combinedBinding, factBase)
			dbBindings.Add(sharedBinding)
		}
	}

	return dbBindings
}

func (solver ProblemSolver) solveMultipleRelationSingleFactBase(relations []mentalese.Relation, binding mentalese.Binding, factBase knowledge.FactBase) mentalese.BindingSet {

	sequenceBindings := mentalese.InitBindingSet(binding)

	for _, relation := range relations {
		sequenceBindings = solver.solveSingleRelationSingleFactBase(relation, sequenceBindings, factBase)
	}

	return sequenceBindings
}

func (solver ProblemSolver) solveSingleRelationSingleFactBase(relation mentalese.Relation, bindings mentalese.BindingSet, factBase knowledge.FactBase) mentalese.BindingSet {

	solver.log.StartProduction("Database" + " " + factBase.GetName(), relation.String() + " " + bindings.String())

	relationBindings := mentalese.NewBindingSet()

	multiFound := false
	aggregateBindings := mentalese.NewBindingSet()

	for _, aggregateBase := range solver.aggregateBases {
		aggregateBindings, multiFound = aggregateBase.Execute(relation, bindings)
		if multiFound {
			relationBindings = aggregateBindings
			break
		}
	}

	if !multiFound {

		for _, binding := range bindings.GetAll() {

			resultBindings := factBase.MatchRelationToDatabase(relation, binding)

			// found bindings must be extended with the bindings already present
			for _, resultBinding := range resultBindings.GetAll() {
				newRelationBinding := binding.Merge(resultBinding)
				relationBindings.Add(newRelationBinding)
			}
		}
	}

	solver.log.EndProduction("Database" + " " + factBase.GetName(), relationBindings.String())

	return relationBindings
}

func (solver ProblemSolver) replaceSharedIdsByLocalIds(binding mentalese.Binding, factBase knowledge.FactBase) mentalese.Binding {

	newBinding := mentalese.NewScopedBinding(binding.GetScope())

	for key, value := range binding.GetAll() {
		newValue := value

		if value.IsId() {
			sharedId := value.TermValue
			entityType := value.TermEntityType
			if entityType != "" {
				localId := factBase.GetLocalId(sharedId, entityType)
				if localId == "" {
					solver.log.AddError(fmt.Sprintf("Local id %s not found for %s in fact base %s", sharedId, entityType, factBase.GetName()))
					return newBinding
				}
				newValue = mentalese.NewTermId(localId, entityType)
			}
		}

		newBinding.Set(key, newValue)
	}

	return newBinding
}

func (solver ProblemSolver) replaceLocalIdBySharedId(binding mentalese.Binding, factBase knowledge.FactBase) mentalese.Binding {

	newBinding := mentalese.NewScopedBinding(binding.GetScope())

	for key, value := range binding.GetAll() {
		newValue := value

		if value.IsId() {
			localId := value.TermValue
			entityType := value.TermEntityType
			if entityType != "" {
				sharedId := factBase.GetSharedId(localId, entityType)
				if sharedId == "" {
					solver.log.AddError(fmt.Sprintf("Shared id %s not found for %s in fact base %s", localId, entityType, factBase.GetName()))
					return newBinding
				}
				newValue = mentalese.NewTermId(sharedId, entityType)
			}
		}

		newBinding.Set(key, newValue)
	}

	return newBinding
}

func (solver ProblemSolver) modifyKnowledgeBase(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	newBindings := mentalese.NewBindingSet()

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
					return mentalese.NewBindingSet()
				}
				binding = solver.replaceLocalIdBySharedId(binding, factBase)
				newBindings.Add(binding)
			}
		} else if argument.IsRule() {
			for _, ruleBase := range solver.ruleBases {
				rule := relation.Arguments[0].TermValueRule.BindSingle(binding)
				ruleBase.Assert(rule)
				newBindings.Add(binding)
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
				newBindings.Add(binding)
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
func (solver ProblemSolver) solveSingleRelationSingleBindingSingleRuleBase(goalRelation mentalese.Relation, binding mentalese.Binding, ruleBase knowledge.RuleBase) mentalese.BindingSet {

	inputVariables := goalRelation.GetVariableNames()

	goalBindings := mentalese.NewBindingSet()

	// match rules from the rule base to the goalRelation
	sourceSubgoalSets, _ := ruleBase.Bind(goalRelation, binding)

	for _, sourceSubgoalSet := range sourceSubgoalSets {

		scope := mentalese.NewScope()
		solver.scopeStack.Push(scope)

		scopedBinding := mentalese.NewScopedBinding(scope).Merge(binding)
		subgoalResultBindings := mentalese.InitBindingSet(scopedBinding)

		for _, subGoal := range sourceSubgoalSet {

			subgoalResultBindings = solver.SolveRelationSet([]mentalese.Relation{subGoal}, subgoalResultBindings)
			if subgoalResultBindings.IsEmpty() {
				break
			}
		}

		for _, subgoalResultBinding := range subgoalResultBindings.GetAll() {

			// filter out the input variables
			filteredBinding := subgoalResultBinding.FilterVariablesByName(inputVariables)

			// make sure all variables of the original binding are present
			goalBinding := scopedBinding.Merge(filteredBinding)

			goalBindings.Add(goalBinding)
		}

		solver.scopeStack.Pop()
	}

	return goalBindings
}
