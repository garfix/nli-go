package central

import (
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
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
	optimizer           Optimizer
	log                 *common.SystemLog
// todo refactor into something more decent
	quantLevel			int
}

func NewProblemSolver(matcher *mentalese.RelationMatcher, log *common.SystemLog) *ProblemSolver {
	return &ProblemSolver{
		factBases:      []knowledge.FactBase{},
		ruleBases:      []knowledge.RuleBase{},
		aggregateBases: []knowledge.AggregateBase{},
		matcher:        matcher,
		optimizer:      NewOptimizer(matcher),
		log:            log,
		quantLevel:     0,
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
func (solver ProblemSolver) SolveRelationSet(set mentalese.RelationSet, nameStore *ResolvedNameStore, bindings []mentalese.Binding) []mentalese.Binding {

	solver.log.StartDebug("SolveRelationSet", set, bindings)

	solver.log.AddProduction("Solve", set.String())

	var newBindings []mentalese.Binding

	// remove duplicates because they cause unnecessary work and the optimizer can't deal with them
	set = set.RemoveDuplicates()

	// sort the relations to reduce the number of tuples retrieved from the fact bases
	solutionRoutes, remainingRelations, ok := solver.optimizer.CreateSolutionRoutes(set, solver.allKnowledgeBases, nameStore)

	solver.log.AddProduction("Solution Routes", solutionRoutes.String())

	if !ok {

		solver.log.AddError("Cannot find these relations in any knowledge base: " + remainingRelations.String())

	} else {

		for _, solutionRoute := range solutionRoutes {
			newBindings = append(newBindings, solver.solveSingleSolutionRouteMultipleBindings(solutionRoute, nameStore, bindings)...)
		}

	}

	// remove duplicates because they cause unnecessary work and they cause problems for the generator
	newBindings = mentalese.UniqueBindings(newBindings)

	solver.log.EndDebug("SolveRelationSet", newBindings)

	return newBindings
}

func (solver ProblemSolver) solveSingleSolutionRouteMultipleBindings(solutionRoute knowledge.SolutionRoute, nameStore *ResolvedNameStore, bindings []mentalese.Binding) []mentalese.Binding {

	newBindings := bindings

	for _, relationGroup := range solutionRoute {
		newBindings = solver.solveSingleRelationGroupMultipleBindings(relationGroup, nameStore, newBindings)

		if len(newBindings) == 0 {
			break
		}
	}

	return newBindings
}

func (solver ProblemSolver) solveSingleRelationGroupMultipleBindings(relationGroup knowledge.RelationGroup, nameStore *ResolvedNameStore, bindings []mentalese.Binding) []mentalese.Binding {

	solver.log.StartDebug("solveSingleRelationGroupMultipleBindings", relationGroup, bindings)

	newBindings := []mentalese.Binding{}

	knowledgeBase := solver.allKnowledgeBases[relationGroup.KnowledgeBaseIndex]
	aggregateBase, isAggregateBase := knowledgeBase.(knowledge.AggregateBase)

	if isAggregateBase {

		mbBindings, ok := aggregateBase.Bind(relationGroup.Relations[0], bindings)

		if ok {
			newBindings = append(newBindings, mbBindings...)
		}

	} else {

		for _, binding := range bindings {
			groupBindings := solver.solveSingleRelationGroupSingleBinding(relationGroup, nameStore, binding)
			newBindings = append(newBindings, groupBindings...)
		}

	}

	solver.log.EndDebug("solveSingleRelationGroupMultipleBindings", newBindings)

	return newBindings
}

func (solver ProblemSolver) solveSingleRelationGroupSingleBinding(relationGroup knowledge.RelationGroup, nameStore *ResolvedNameStore, binding mentalese.Binding) []mentalese.Binding {

	solver.log.StartDebug("solveSingleRelationGroupSingleBinding", relationGroup, binding)

	knowledgeBase := solver.allKnowledgeBases[relationGroup.KnowledgeBaseIndex]
	factBase, isFactBase := knowledgeBase.(knowledge.FactBase)
	ruleBase, isRuleBase := knowledgeBase.(knowledge.RuleBase)
	functionBase, isFunctionBase := knowledgeBase.(knowledge.FunctionBase)
	_, isNestedStructureBase := knowledgeBase.(knowledge.NestedStructureBase)

	boundRelations := solver.matcher.BindRelationSetSingleBinding(relationGroup.Relations, binding)

	var newBindings []mentalese.Binding

	if isFactBase {

//		resolvedBoundRelations := solver.bindKnowledgeBaseVariables(boundRelations, nameStore, factBase.GetName())

		sourceBindings := solver.FindFacts(factBase, boundRelations)

		for _, sourceBinding := range sourceBindings {

			combinedBinding := binding.Merge(sourceBinding)
			newBindings = append(newBindings, combinedBinding)
		}

	} else if isFunctionBase {

		var relation = relationGroup.Relations[0]
		resultBinding, functionFound := functionBase.Execute(relation, binding)
		if functionFound {
			newBindings = append(newBindings, resultBinding)
		}

	} else if isRuleBase {

		newBindings = append(newBindings, solver.SolveSingleRelationSingleBindingSingleRuleBase(boundRelations[0], nameStore, binding, ruleBase)...)

	} else if isNestedStructureBase {

		newBindings = solver.SolveChildStructures(relationGroup.Relations[0], nameStore, binding)

	}

	solver.log.EndDebug("solveSingleRelationGroupSingleBinding", newBindings)

	return newBindings
}


//func (solver ProblemSolver) bindKnowledgeBaseVariables(set mentalese.RelationSet, nameStore *ResolvedNameStore, knowledgeBaseName string) mentalese.RelationSet {
//
//	values := nameStore.GetValues(knowledgeBaseName)
//
//	binding := mentalese.Binding{}
//
//	for key, value := range values {
//		binding[key] = mentalese.NewId(value)
//	}
//
//	boundRelations := solver.matcher.BindRelationSetSingleBinding(set, binding)
//
//	return boundRelations
//}


func (solver ProblemSolver) SolveChildStructures(goal mentalese.Relation, nameStore *ResolvedNameStore, binding mentalese.Binding) []mentalese.Binding {

	solver.log.StartDebug("NestedStructureBase BindChildStructures", goal, binding)

	var newBindings []mentalese.Binding

	if goal.Predicate == mentalese.Predicate_Quant {

		newBindings = solver.SolveQuant(goal, nameStore, binding)
	}

	solver.log.EndDebug("NestedStructureBase BindChildStructures", newBindings)

	return newBindings
}


func (solver ProblemSolver) FindFacts(factBase knowledge.FactBase, goal mentalese.RelationSet) []mentalese.Binding {

	solver.log.StartDebug("FindFacts", goal)

	subgoalBindings := []mentalese.Binding{}

	for _, ds2db := range factBase.GetMappings() {

		// gender(14, G), gender(A, male) => internalBinding: A = 14
		internalBindingsX, match1 := solver.matcher.MatchSequenceToSet(ds2db.Pattern, goal, mentalese.Binding{})
		if match1 {

			internalBinding := internalBindingsX[0]

			// gender(14, G), gender(A, male) => externalBinding: G = male
			externalBindings, match2 := solver.matcher.MatchSequenceToSet(goal, ds2db.Pattern, mentalese.Binding{})
			if match2 {

				externalBinding := externalBindings[0]

				// create a version of the conditions with bound variables
				boundConditions := solver.matcher.BindRelationSetSingleBinding(ds2db.Replacement, internalBinding)

				// match this bound version to the database
				internalBindings, match3 := solver.SolveMultipleRelationSingleFactBase(ds2db.Replacement, boundConditions, factBase)

				if match3 {
					for _, binding := range internalBindings {
						subgoalBindings = append(subgoalBindings, externalBinding.Intersection(binding))
					}
				}
			}
		}
	}

	solver.log.EndDebug("FindFacts", subgoalBindings)

	return subgoalBindings
}

func (solver ProblemSolver) SolveMultipleRelationSingleFactBase(unboundSequence []mentalese.Relation, boundSequence []mentalese.Relation, factBase knowledge.FactBase) ([]mentalese.Binding, bool) {

	solver.log.StartDebug("SolveMultipleRelationSingleFactBase", boundSequence)

	// bindings using database level variables
	sequenceBindings := []mentalese.Binding{}
	match := true

	for i, relation := range boundSequence {

		relationBindings := []mentalese.Binding{}

		aggregateFunctionFound := false
		for _, aggregateBase := range solver.aggregateBases {
			newRelationBindings, ok := aggregateBase.Bind(unboundSequence[i], sequenceBindings)
			if ok {
				relationBindings = newRelationBindings
				aggregateFunctionFound = true
				break
			}
		}

		if !aggregateFunctionFound {

			if len(sequenceBindings) == 0 {

				resultBindings := factBase.MatchRelationToDatabase(relation)
				relationBindings = resultBindings

			} else {

				//functionBindings, functionFound := solver.matcher.MatchRelationToFunction(relation, sequenceBindings)
				//if functionFound {
				//
				//	relationBindings = functionBindings
				//
				//} else {

				//// go through the bindings resulting from previous relation
				for _, binding := range sequenceBindings {

					boundRelation := solver.matcher.BindSingleRelationSingleBinding(relation, binding)
					resultBindings := factBase.MatchRelationToDatabase(boundRelation)

					// found bindings must be extended with the bindings already present
					for _, resultBinding := range resultBindings {
						newRelationBinding := binding.Merge(resultBinding)
						relationBindings = append(relationBindings, newRelationBinding)
					}
				}
				//			}

			}
		}

		sequenceBindings = relationBindings

		if len(sequenceBindings) == 0 {
			match = false
			break
		}
	}

	solver.log.EndDebug("SolveMultipleRelationSingleFactBase", sequenceBindings, match)

	return sequenceBindings, match
}

// goalRelation e.g. father('jack', Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) SolveSingleRelationSingleBindingSingleRuleBase(goalRelation mentalese.Relation, nameStore *ResolvedNameStore, binding mentalese.Binding, ruleBase knowledge.RuleBase) []mentalese.Binding {

	solver.log.StartDebug("SolveSingleRelationSingleBindingSingleRuleBase", goalRelation, binding)

	for _, val := range binding {
		if val.TermType == mentalese.Term_variable {
			panic("Variable bound to variable")
		}
	}

	goalBindings := []mentalese.Binding{}

	// match rules from the rule base to the goalRelation
	boundRelation := solver.matcher.BindSingleRelationSingleBinding(goalRelation, binding)
	sourceSubgoalSets, sourceBindings := ruleBase.Bind(boundRelation)

	for i, sourceSubgoalSet := range sourceSubgoalSets {

		// sourceBinding: from subgoal variable to goal argument
		sourceBinding := sourceBindings[i]

		// subgoalBinding: from subgoal variable to goal constant
		subgoalBinding := sourceBinding.RemoveVariables()

		subgoalResultBindings := solver.SolveRelationSet(sourceSubgoalSet, nameStore, []mentalese.Binding{subgoalBinding})

		// subgoalResultBinding: from subgoal variables to constants (contains temporary variables)
		for _, subgoalResultBinding := range subgoalResultBindings {

			// sourceBinding e.g. { A:X, B:'yellow' }
			// after swap { X:A }
			// bind { X:A } with { A:'red', B:'yellow', C:'blue' } results in { X:'red' }
			convertedBinding := sourceBinding.Swap().Bind(subgoalResultBinding)

			// start extending the new binding with goalRelation variables as keys
			goalBinding := binding.Merge(convertedBinding)
			goalBindings = append(goalBindings, goalBinding)
		}
	}

	solver.log.EndDebug("SolveSingleRelationSingleBindingSingleRuleBase", goalBindings)

	return goalBindings
}
