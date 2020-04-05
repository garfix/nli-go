package central

import (
	"fmt"
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"strings"
)

// The problem solver takes a relation set and a set of bindings
// and returns a set of new bindings
// It uses knowledge bases to find these bindings
type ProblemSolver struct {
	allKnowledgeBases   []knowledge.KnowledgeBase
	factBases           []knowledge.FactBase
	ruleBases           []knowledge.RuleBase
	functionBases		[]knowledge.FunctionBase
	aggregateBases      []knowledge.AggregateBase
	nestedStructureBase []knowledge.NestedStructureBase
	matcher             *mentalese.RelationMatcher
	predicates 			mentalese.Predicates
	modifier            *FactBaseModifier
	dialogContext 		*DialogContext
	log                 *common.SystemLog
	SolveDepth int
}

func NewProblemSolver(matcher *mentalese.RelationMatcher, predicates mentalese.Predicates, dialogContext *DialogContext, log *common.SystemLog) *ProblemSolver {
	return &ProblemSolver{
		factBases:      []knowledge.FactBase{},
		ruleBases:      []knowledge.RuleBase{},
		functionBases:  []knowledge.FunctionBase{},
		aggregateBases: []knowledge.AggregateBase{},
		matcher:        matcher,
		predicates:		predicates,
		modifier:       NewFactBaseModifier(log),
		dialogContext:  dialogContext,
		log:            log,
		SolveDepth:     0,
	}
}

func (solver *ProblemSolver) AddFactBase(factBase knowledge.FactBase) {
	solver.factBases = append(solver.factBases, factBase)
	solver.allKnowledgeBases = append(solver.allKnowledgeBases, factBase)
}

func (solver *ProblemSolver) AddFunctionBase(functionBase knowledge.FunctionBase) {
	solver.functionBases = append(solver.functionBases, functionBase)
	solver.allKnowledgeBases = append(solver.allKnowledgeBases, functionBase)
}

func (solver *ProblemSolver) AddRuleBase(ruleBase knowledge.RuleBase) {
	solver.ruleBases = append(solver.ruleBases, ruleBase)
	solver.allKnowledgeBases = append(solver.allKnowledgeBases, ruleBase)
}

func (solver *ProblemSolver) AddMultipleBindingsBase(source knowledge.AggregateBase) {
	solver.aggregateBases = append(solver.aggregateBases, source)
	solver.allKnowledgeBases = append(solver.allKnowledgeBases, source)
}

func (solver *ProblemSolver) AddNestedStructureBase(base knowledge.NestedStructureBase) {
	solver.nestedStructureBase = append(solver.nestedStructureBase, base)
	solver.allKnowledgeBases = append(solver.allKnowledgeBases, base)
}

// set e.g. [ father(X, Y) father(Y, Z) ]
// bindings e.g. [{X: john, Z: jack} {}]
// return e.g. [
//  { X: john, Z: jack, Y: billy }
//  { X: john, Z: jack, Y: bob }
// ]
func (solver ProblemSolver) SolveRelationSet(set mentalese.RelationSet, bindings mentalese.Bindings) mentalese.Bindings {

	solver.SolveDepth++

	head := strings.Repeat("  ", solver.SolveDepth)

	solver.log.AddProduction(head + "Solve Set", set.String() + " " + bindings.String())

	newBindings := bindings
	for _, relation := range set {
		newBindings = solver.SolveSingleRelationMultipleBindings(relation, newBindings)

		if len(newBindings) == 0 {
			break
		}
	}

	// remove duplicates because they cause unnecessary work and they cause problems for the generator
	newBindings = mentalese.UniqueBindings(newBindings)

	solver.log.AddProduction(head + "Solve Set", newBindings.String())

	solver.SolveDepth--

	return newBindings
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
func (solver ProblemSolver) SolveSingleRelationMultipleBindings(relation mentalese.Relation, bindings []mentalese.Binding) []mentalese.Binding {

	solver.log.StartDebug("SolveSingleRelationMultipleBindings", relation, bindings)

	newBindings := []mentalese.Binding{}
	multiFound := false

	for _, aggregateBase := range solver.aggregateBases {
		newBindings, multiFound = aggregateBase.Bind(relation, bindings)
		if multiFound {
			break
		}
	}

	if !multiFound {

		if len(bindings) == 0 {
			newBindings = solver.SolveSingleRelationSingleBinding(relation, mentalese.Binding{})
		} else {
			for _, binding := range bindings {
				newBindings = append(newBindings, solver.SolveSingleRelationSingleBinding(relation, binding)...)
			}
		}
	}

	solver.log.EndDebug("SolveSingleRelationMultipleBindings", newBindings)

	return newBindings
}

// goalRelation e.g. father(Y, Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) SolveSingleRelationSingleBinding(relation mentalese.Relation, binding mentalese.Binding) []mentalese.Binding {

	solver.log.StartDebug("SolveSingleRelationSingleBinding", relation, binding)

	newBindings := []mentalese.Binding{}

	// go through all fact bases
	for _, factBase := range solver.factBases {
		newBindings = append(newBindings, solver.SolveSingleRelationSingleBindingSingleFactBase(relation, binding, factBase)...)
	}

	// go through all rule bases
	for _, ruleBase := range solver.ruleBases {
		newBindings = append(newBindings, solver.SolveSingleRelationSingleBindingSingleRuleBase(relation, binding, ruleBase)...)
	}

	// go through all function bases
	for _, functionBase := range solver.functionBases {
		resultBinding, functionFound := functionBase.Execute(relation, binding)
		if functionFound {
			newBindings = append(newBindings, resultBinding)
		}
	}

	// go through all nested structure bases
	newBindings = append(newBindings, solver.solveChildStructures(relation, binding)...)

	solver.log.EndDebug("SolveSingleRelationSingleBinding", newBindings)

	return newBindings
}

func (solver ProblemSolver) solveChildStructures(goal mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	solver.log.StartDebug("NestedStructureBase BindChildStructures", goal, binding)

	var newBindings mentalese.Bindings

	if goal.Predicate == mentalese.PredicateQuant {

		newBindings = solver.SolveQuant(goal, binding)

	} else if goal.Predicate == mentalese.PredicateSequence {

		newBindings = solver.SolveSeq(goal, binding)

	} else if goal.Predicate == mentalese.PredicateNot {

		newBindings = solver.SolveNot(goal, binding)

	} else if goal.Predicate == mentalese.PredicateCall {

		newBindings = solver.Call(goal, binding)

	}

	solver.log.EndDebug("NestedStructureBase BindChildStructures", newBindings)

	return newBindings
}

// Creates bindings for the free variables in 'relations', by resolving them in factBase
func (solver ProblemSolver) FindFacts(factBase knowledge.FactBase, relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	solver.log.StartDebug("FindFacts", relation, binding)

	dbBindings := mentalese.Bindings{}

	for _, ds2db := range factBase.GetMappings() {

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

	solver.log.EndDebug("FindFacts", dbBindings)

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

	solver.SolveDepth++

	head := strings.Repeat("  ", solver.SolveDepth)

	solver.log.AddProduction(head + "Database" + " " + factBase.GetName(), relation.String() + " " + bindings.String())

	relationBindings := mentalese.Bindings{}

	aggregateFunctionFound := false
	for _, aggregateBase := range solver.aggregateBases {
		newRelationBindings, ok := aggregateBase.Bind(relation, bindings)
		if ok {
			relationBindings = newRelationBindings
			aggregateFunctionFound = true
			break
		}
	}

	if !aggregateFunctionFound {

		for _, binding := range bindings {

			resultBindings := factBase.MatchRelationToDatabase(relation, binding)

			// found bindings must be extended with the bindings already present
			for _, resultBinding := range resultBindings {
				newRelationBinding := binding.Merge(resultBinding)
				relationBindings = append(relationBindings, newRelationBinding)
			}
		}
	}

	solver.log.AddProduction(head + "Database" + " " + factBase.GetName(), relationBindings.String())

	solver.SolveDepth--

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
				newValue = mentalese.NewId(localId, entityType)
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
				newValue = mentalese.NewId(sharedId, entityType)
			}
		}

		newBinding[key] = newValue
	}

	return newBinding
}

func (solver ProblemSolver) SolveSingleRelationSingleBindingSingleFactBase(relation mentalese.Relation, binding mentalese.Binding, factBase knowledge.FactBase) mentalese.Bindings {

	newBindings := mentalese.Bindings{}

	if relation.Predicate == mentalese.PredicateAssert {

		localIdBinding := solver.replaceSharedIdsByLocalIds(binding, factBase)
		boundRelation := relation.BindSingleRelationSingleBinding(localIdBinding)
		solver.modifier.Assert(boundRelation.Arguments[0].TermValueRelationSet[0], factBase)
		binding = solver.replaceLocalIdBySharedId(binding, factBase)
		newBindings = append(newBindings, binding)

	} else if relation.Predicate == mentalese.PredicateRetract {

		localIdBinding := solver.replaceSharedIdsByLocalIds(binding, factBase)
		boundRelation := relation.BindSingleRelationSingleBinding(localIdBinding)
		solver.modifier.Retract(boundRelation.Arguments[0].TermValueRelationSet[0], factBase)
		binding = solver.replaceLocalIdBySharedId(binding, factBase)
		newBindings = append(newBindings, binding)

	} else {

		newBindings = solver.FindFacts(factBase, relation, binding)
	}

	return newBindings
}

// goalRelation e.g. father('jack', Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) SolveSingleRelationSingleBindingSingleRuleBase(goalRelation mentalese.Relation, binding mentalese.Binding, ruleBase knowledge.RuleBase) mentalese.Bindings {

	solver.log.StartDebug("SolveSingleRelationSingleBindingSingleRuleBase", goalRelation, binding)

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

	solver.log.EndDebug("SolveSingleRelationSingleBindingSingleRuleBase", goalBindings)

	return goalBindings
}
