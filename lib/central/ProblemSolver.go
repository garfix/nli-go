package central

import (
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
)

// The problem solver takes a relation set and a set of bindings
// and returns a set of new bindings
// It uses knowledge bases to find these bindings
type ProblemSolver struct {
	allKnowledgeBases     []knowledge.KnowledgeBase
	factBases             []knowledge.FactBase
	ruleBases             []knowledge.RuleBase
	multipleBindingsBases []knowledge.MultipleBindingsBase
	nestedStructureBase   []knowledge.NestedStructureBase
	matcher               *mentalese.RelationMatcher
	optimizer			  Optimizer
	log                   *common.SystemLog
}

func NewProblemSolver(matcher *mentalese.RelationMatcher, log *common.SystemLog) *ProblemSolver {
	return &ProblemSolver{
		factBases:             []knowledge.FactBase{},
		ruleBases:             []knowledge.RuleBase{},
		multipleBindingsBases: []knowledge.MultipleBindingsBase{},
		matcher:               matcher,
		optimizer: 			   NewOptimizer(matcher),
		log:                   log,
	}
}

func (solver *ProblemSolver) AddFactBase(factBase knowledge.FactBase) {
	solver.factBases = append(solver.factBases, factBase)
	solver.allKnowledgeBases = append(solver.allKnowledgeBases, factBase)
}

func (solver *ProblemSolver) AddRuleBase(ruleBase knowledge.RuleBase) {
	solver.ruleBases = append(solver.ruleBases, ruleBase)
	solver.allKnowledgeBases = append(solver.allKnowledgeBases, ruleBase)
}

func (solver *ProblemSolver) AddMultipleBindingsBase(source knowledge.MultipleBindingsBase) {
	solver.multipleBindingsBases = append(solver.multipleBindingsBases, source)
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
func (solver ProblemSolver) SolveRelationSet(set mentalese.RelationSet, bindings []mentalese.Binding) []mentalese.Binding {

	solver.log.StartDebug("SolveRelationSet", set, bindings)

	newBindings := []mentalese.Binding{}

	// remove duplicates because they cause unnecessary work and the optimizer can't deal with them
	set = set.RemoveDuplicates()

	// sort the relations to reduce the number of tuples retrieved from the fact bases
	solutionRoutes, remainingRelations, ok := solver.optimizer.CreateSolutionRoutes(set, solver.allKnowledgeBases)

	solver.log.AddProduction("Solution Routes", solutionRoutes.String())

	if !ok {

		solver.log.AddError("Cannot find these relations in any knowledge base: " + remainingRelations.String())

	} else {

		for _, solutionRoute := range solutionRoutes {
			newBindings = append(newBindings, solver.solveSingleSolutionRouteMultipleBindings(solutionRoute, bindings)...)
		}

	}

	solver.log.EndDebug("SolveRelationSet", newBindings)

	return newBindings
}

func (solver ProblemSolver) solveSingleSolutionRouteMultipleBindings(solutionRoute knowledge.SolutionRoute, bindings []mentalese.Binding) []mentalese.Binding {

	newBindings := bindings

	for _, relationGroup := range solutionRoute {
		newBindings = solver.solveSingleRelationGroupMultipleBindings(relationGroup, newBindings)

		if len(newBindings) == 0 {
			break
		}
	}

	return newBindings
}

func (solver ProblemSolver) solveSingleRelationGroupMultipleBindings(relationGroup knowledge.RelationGroup, bindings []mentalese.Binding) []mentalese.Binding {

	solver.log.StartDebug("solveSingleRelationGroupMultipleBindings", relationGroup, bindings)

	newBindings := []mentalese.Binding{}

	knowledgeBase := solver.allKnowledgeBases[relationGroup.KnowledgeBaseIndex]
	multipleBindingsBase, isMultipleBindingsBase := knowledgeBase.(knowledge.MultipleBindingsBase)

	if isMultipleBindingsBase {

		mbBindings, ok := multipleBindingsBase.Bind(relationGroup.Relations[0], bindings)

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

func (solver ProblemSolver) solveSingleRelationGroupSingleBinding(relationGroup knowledge.RelationGroup, binding mentalese.Binding) []mentalese.Binding {

	solver.log.StartDebug("solveSingleRelationGroupSingleBinding", relationGroup, binding)

	knowledgeBase := solver.allKnowledgeBases[relationGroup.KnowledgeBaseIndex]
	factBase, isFactBase := knowledgeBase.(knowledge.FactBase)
	ruleBase, isRuleBase := knowledgeBase.(knowledge.RuleBase)
	_, isNestedStructureBase := knowledgeBase.(knowledge.NestedStructureBase)

	boundRelations := solver.matcher.BindRelationSetSingleBinding(relationGroup.Relations, binding)

	newBindings := []mentalese.Binding{}

	if isNestedStructureBase {

		newBindings = solver.SolveChildStructures(relationGroup.Relations[0], binding)

	} else if isFactBase {

		sourceBindings := solver.FindFacts(factBase, boundRelations)

		for _, sourceBinding := range sourceBindings {

			combinedBinding := binding.Merge(sourceBinding)
			newBindings = append(newBindings, combinedBinding)
		}

	} else if isRuleBase {

		newBindings = append(newBindings, solver.SolveSingleRelationSingleBindingSingleRuleBase(boundRelations[0], binding, ruleBase)...)

	}

	solver.log.EndDebug("solveSingleRelationGroupSingleBinding", newBindings)

	return newBindings
}


func (solver ProblemSolver) SolveChildStructures(goal mentalese.Relation, binding mentalese.Binding) []mentalese.Binding {

	solver.log.StartDebug("NestedStructureBase BindChildStructures", goal, binding)

	newBindings := []mentalese.Binding{}

	if goal.Predicate == mentalese.Predicate_Quant {

		newBindings = solver.SolveQuant(goal, binding)
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

				// match1 this bound version to the database
				internalBindings, match3 := factBase.Bind(boundConditions)

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

// goalRelation e.g. father('jack', Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) SolveSingleRelationSingleBindingSingleRuleBase(goalRelation mentalese.Relation, binding mentalese.Binding, ruleBase knowledge.RuleBase) []mentalese.Binding {

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

		subgoalResultBindings := solver.SolveRelationSet(sourceSubgoalSet, []mentalese.Binding{subgoalBinding})

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
