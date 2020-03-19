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
//
// Structures used:
// RelationSet: describes the query
// SolutionRoute: describes a single path through a series of knowledge base calls
// RelationGroup: a single step in a solution route; it is the input for a specified knowledge base, having a calculation cost
type ProblemSolver struct {
	allKnowledgeBases   []knowledge.KnowledgeBase
	factBases           []knowledge.FactBase
	ruleBases           []knowledge.RuleBase
	aggregateBases      []knowledge.AggregateBase
	nestedStructureBase []knowledge.NestedStructureBase
	matcher             *mentalese.RelationMatcher
	predicates 			mentalese.Predicates
	optimizer           Optimizer
	modifier            *FactBaseModifier
	dialogContext 		*DialogContext
	log                 *common.SystemLog
	SolveDepth int
}

func NewProblemSolver(matcher *mentalese.RelationMatcher, predicates mentalese.Predicates, dialogContext *DialogContext, log *common.SystemLog) *ProblemSolver {
	return &ProblemSolver{
		factBases:      []knowledge.FactBase{},
		ruleBases:      []knowledge.RuleBase{},
		aggregateBases: []knowledge.AggregateBase{},
		matcher:        matcher,
		predicates:		predicates,
		optimizer:      NewOptimizer(matcher),
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

	newBindings := mentalese.Bindings{}

	// remove duplicates because they cause unnecessary work and the optimizer can't deal with them
	set = set.RemoveDuplicates()

	// sort the relations to reduce the number of tuples retrieved from the fact bases
	solutionRoutes, remainingRelations, ok := solver.optimizer.CreateSolutionRoutes(set, solver.allKnowledgeBases)

	solver.log.AddProduction(head + "Solution Routes", solutionRoutes.String())

	if !ok {

		solver.log.AddError("Cannot find these relations in any knowledge base: " + remainingRelations.String())

	} else {

		for _, solutionRoute := range solutionRoutes {
			newBindings = append(newBindings, solver.solveSingleSolutionRouteMultipleBindings(solutionRoute, bindings)...)
		}
	}

	// remove duplicates because they cause unnecessary work and they cause problems for the generator
	newBindings = mentalese.UniqueBindings(newBindings)

	solver.log.AddProduction(head + "Solution Routes Result", newBindings.String())

	solver.SolveDepth--

	solver.log.EndDebug("SolveRelationSet", newBindings)

	return newBindings
}

func (solver ProblemSolver) solveSingleSolutionRouteMultipleBindings(solutionRoute knowledge.SolutionRoute, bindings mentalese.Bindings) mentalese.Bindings {

	newBindings := bindings

	solver.SolveDepth++

	head := strings.Repeat("  ", solver.SolveDepth)

	for _, relationGroup := range solutionRoute {

		solver.log.AddProduction(head + "Solve RelationGroup", relationGroup.String() + " " + bindings.String())

		newBindings = solver.solveSingleRelationGroupMultipleBindings(relationGroup, newBindings)

		solver.log.AddProduction(head + "Solve RelationGroup Result", newBindings.String())

		if len(newBindings) == 0 {
			break
		}
	}

	solver.SolveDepth--

	return newBindings
}

func (solver ProblemSolver) findKnowledgeBaseByName(name string) knowledge.KnowledgeBase {
	for _, knowledgeBase := range solver.allKnowledgeBases {
		if knowledgeBase.GetName() == name {
			return knowledgeBase
		}
	}

	return nil
}

func (solver ProblemSolver) solveSingleRelationGroupMultipleBindings(relationGroup knowledge.RelationGroup, bindings mentalese.Bindings) mentalese.Bindings {

	solver.log.StartDebug("solveSingleRelationGroupMultipleBindings", relationGroup, bindings)

	newBindings := mentalese.Bindings{}

	knowledgeBase := solver.findKnowledgeBaseByName(relationGroup.KnowledgeBaseName)
	aggregateBase, isAggregateBase := knowledgeBase.(knowledge.AggregateBase)

	if isAggregateBase {

		mbBindings, ok := aggregateBase.Bind(relationGroup.Relations[0], bindings)

		if ok {
			newBindings = append(newBindings, mbBindings...)
		}

	} else {

		for _, binding := range bindings {
			groupBindings := solver.solveSingleRelationGroupSingleBinding(relationGroup, binding)
			newBindings = append(newBindings, groupBindings...)
		}
	}

	solver.log.EndDebug("solveSingleRelationGroupMultipleBindings", newBindings)

	return newBindings
}

func (solver ProblemSolver) solveSingleRelationGroupSingleBinding(relationGroup knowledge.RelationGroup, binding mentalese.Binding) mentalese.Bindings {

	solver.log.StartDebug("solveSingleRelationGroupSingleBinding", relationGroup, binding)

	knowledgeBase := solver.findKnowledgeBaseByName(relationGroup.KnowledgeBaseName)
	factBase, isFactBase := knowledgeBase.(knowledge.FactBase)
	ruleBase, isRuleBase := knowledgeBase.(knowledge.RuleBase)
	functionBase, isFunctionBase := knowledgeBase.(knowledge.FunctionBase)
	_, isNestedStructureBase := knowledgeBase.(knowledge.NestedStructureBase)

	var newBindings mentalese.Bindings

	if isFactBase {

		newBindings = solver.solveSingleRelationGroupSingleBindingFactBase(relationGroup, binding, factBase)

	} else if isFunctionBase {

		var relation = relationGroup.Relations[0]
		resultBinding, functionFound := functionBase.Execute(relation, binding)
		if functionFound {
			newBindings = append(newBindings, resultBinding)
		}

	} else if isRuleBase {

		newBindings = append(newBindings, solver.SolveSingleRelationSingleBindingSingleRuleBase(relationGroup.Relations[0], binding, ruleBase)...)

	} else if isNestedStructureBase {

		newBindings = solver.solveChildStructures(relationGroup.Relations[0], binding)

	}

	solver.log.EndDebug("solveSingleRelationGroupSingleBinding", newBindings)

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

func (solver ProblemSolver) solveSingleRelationGroupSingleBindingFactBase(relationGroup knowledge.RelationGroup, binding mentalese.Binding, factBase knowledge.FactBase) mentalese.Bindings {

	newBindings := mentalese.Bindings{}

	if len(relationGroup.Relations) == 1 && relationGroup.Relations[0].Predicate == mentalese.PredicateAssert {

		localIdBinding := solver.replaceSharedIdsByLocalIds(relationGroup.Relations, binding, factBase)
		boundRelations := relationGroup.Relations.BindSingle(localIdBinding)
		solver.modifier.Assert(boundRelations[0].Arguments[0].TermValueRelationSet, factBase)
		binding = solver.replaceLocalIdBySharedId(relationGroup.Relations, binding, factBase)
		newBindings = append(newBindings, binding)

	} else if len(relationGroup.Relations) == 1 && relationGroup.Relations[0].Predicate == mentalese.PredicateRetract {

		localIdBinding := solver.replaceSharedIdsByLocalIds(relationGroup.Relations, binding, factBase)
		boundRelations := relationGroup.Relations.BindSingle(localIdBinding)
		solver.modifier.Retract(boundRelations[0].Arguments[0].TermValueRelationSet, factBase)
		binding = solver.replaceLocalIdBySharedId(relationGroup.Relations, binding, factBase)
		newBindings = append(newBindings, binding)

	} else {

		newBindings = solver.FindFacts(factBase, relationGroup.Relations, binding)
	}

	return newBindings
}

// Creates bindings for the free variables in 'relations', by resolving them in factBase
func (solver ProblemSolver) FindFacts(factBase knowledge.FactBase, relations mentalese.RelationSet, binding mentalese.Binding) mentalese.Bindings {

	solver.log.StartDebug("FindFacts", relations, binding)

	// shared ids to local ids
	localIdBinding := solver.replaceSharedIdsByLocalIds(relations, binding, factBase)

	dbBindings := mentalese.Bindings{}

	for _, ds2db := range factBase.GetMappings() {

		activeBindings, match := solver.matcher.MatchSequenceToSet(relations, ds2db.Pattern, mentalese.Binding{})
		if !match { continue }
		activeBinding := activeBindings[0]

		activeBindings2, match2 := solver.matcher.MatchSequenceToSet(ds2db.Pattern, relations, mentalese.Binding{})
		if !match2 { continue }
		activeBinding2 := activeBindings2[0]

		dbRelations := ds2db.Replacement.ImportBinding(activeBinding2)

		relevantBinding := localIdBinding.Select(dbRelations.GetVariableNames())
		newDbBindings := solver.solveMultipleRelationSingleFactBase(dbRelations, relevantBinding, factBase)

		for _, newDbBinding := range newDbBindings {
			dbBindings = append(dbBindings, activeBinding.Merge(newDbBinding))
		}
	}

	// adapt bindings for outer world
	newBindings := mentalese.Bindings{}
	for _, dbBinding := range dbBindings {
		// extend the original binding with the db binding
		combinedBinding := localIdBinding.Merge(dbBinding.Select(relations.GetVariableNames()))
		// local ids to shared ids
		sharedBinding := solver.replaceLocalIdBySharedId(relations, combinedBinding, factBase)
		newBindings = append(newBindings, sharedBinding)
	}

	solver.log.EndDebug("FindFacts", newBindings)

	return newBindings
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

			boundRelation := relation.BindSingleRelationSingleBinding(binding)
			resultBindings := factBase.MatchRelationToDatabase(boundRelation)

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

func (solver ProblemSolver) replaceSharedIdsByLocalIds(relationSet mentalese.RelationSet, binding mentalese.Binding, factBase knowledge.FactBase) mentalese.Binding {

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

func (solver ProblemSolver) replaceLocalIdBySharedId(relationSet mentalese.RelationSet, binding mentalese.Binding, factBase knowledge.FactBase) mentalese.Binding {

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
