package central

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/knowledge"
	"nli-go/lib/common"
)

type simpleProblemSolver struct {
	sources []knowledge.SimpleKnowledgeBase
	matcher *mentalese.SimpleRelationMatcher
}

func NewSimpleProblemSolver() *simpleProblemSolver {
	return &simpleProblemSolver{sources: []knowledge.SimpleKnowledgeBase{}, matcher:mentalese.NewSimpleRelationMatcher()}
}

func (solver *simpleProblemSolver) AddKnowledgeBase(source knowledge.SimpleKnowledgeBase) {
	solver.sources = append(solver.sources, source)
}

// goals e.g. { father(X, Y), father(Y, Z)}
// return e.g. {
//  { father('john', 'jack'), father('jack', 'joe') }
//  { father('bob', 'jonathan'), father('jonathan', 'bill') }
// }
func (solver simpleProblemSolver) Solve(goals []mentalese.SimpleRelation) []mentalese.SimpleRelationSet {

	common.Log("Solve start\n")
	bindings := solver.SolveMultipleRelations(goals)
	solutions := solver.matcher.BindMultipleRelationsMultipleBindings(goals, bindings)

	common.Logf("Solve end %v\n", solutions)
	return solutions
}

// goals e.g. { father(X, Y), father(Y, Z)}
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver simpleProblemSolver) SolveMultipleRelations(goals []mentalese.SimpleRelation) []mentalese.SimpleBinding {

	common.Logf("SolveMultipleRelations start %v\n", goals)
	
	bindings := []mentalese.SimpleBinding{}

	for _, goal := range goals {
		bindings = solver.SolveSingleRelationMultipleBindings(goal, bindings)
	}

	common.Logf("SolveMultipleRelations end %v\n", bindings)

	return bindings
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
func (solver simpleProblemSolver) SolveSingleRelationMultipleBindings(goalRelation mentalese.SimpleRelation, bindings []mentalese.SimpleBinding) []mentalese.SimpleBinding {

	common.Logf("SolveSingleRelationMultipleBindings start %v ; %v\n", goalRelation, bindings)

	if len(bindings) == 0 {
		return solver.SolveSingleRelationSingleBinding(goalRelation, mentalese.SimpleBinding{})
	}

	newBindings := []mentalese.SimpleBinding{}

	for _, binding := range bindings {
		newBindings = append(newBindings, solver.SolveSingleRelationSingleBinding(goalRelation, binding)...)
	}

	common.Logf("SolveSingleRelationMultipleBindings end %v\n", newBindings)

	return newBindings
}

// goalRelation e.g. father(Y, Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver simpleProblemSolver) SolveSingleRelationSingleBinding(goalRelation mentalese.SimpleRelation, binding mentalese.SimpleBinding) []mentalese.SimpleBinding {

	common.Logf("SolveSingleRelationSingleBinding start %v %v\n", goalRelation, binding)

	newBindings := []mentalese.SimpleBinding{}

	boundRelation := solver.matcher.BindSingleRelationSingleBinding(goalRelation, binding)

	// go through all knowledge sources
	for _, source := range solver.sources {
		newBindings = append(newBindings, solver.SolveSingleRelationSingleBindingSingleSource(boundRelation, binding, source)...)
	}

	common.Logf("SolveSingleRelationSingleBinding end %v\n", newBindings)

	return newBindings
}

// boundRelation e.g. father('jack', Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver simpleProblemSolver) SolveSingleRelationSingleBindingSingleSource(boundRelation mentalese.SimpleRelation, binding mentalese.SimpleBinding, source knowledge.SimpleKnowledgeBase) []mentalese.SimpleBinding {

	common.Logf("SolveSingleRelationSingleBindingSingleSource start %v %v\n", boundRelation, binding)

	newBindings := []mentalese.SimpleBinding{}

	// boundRelation e.g. father(X, 'john')
	// subgoalSets e.g. {
	//    { male(X), parent(X, 'john') },
	//    { child('john', X), male(X) }
	// }
	// bindings e.g. {
	//    { X='Jack' },
	// }
	// Note: bindings are linked to subgoalSets, one on one; but usually just one of the arrays is used
	sourceSubgoalSets, sourceBindings := source.Bind(boundRelation)

	for i, sourceSubgoalSet := range sourceSubgoalSets {

		sourceBinding := sourceBindings[i]

		combinedBinding := binding.Merge(sourceBinding)

		subgoalSetBindings := solver.SolveMultipleRelationsSingleBinding(sourceSubgoalSet, combinedBinding)
		newBindings = append(newBindings, subgoalSetBindings...)
	}

	common.Logf("SolveSingleRelationSingleBindingSingleSource end %v\n", newBindings)

	return newBindings
}

// goal e.g. { father(X, Y), father(Y, Z)}
// bindings {X='john', Y='jack'}
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver simpleProblemSolver) SolveMultipleRelationsSingleBinding(goals []mentalese.SimpleRelation, binding mentalese.SimpleBinding) []mentalese.SimpleBinding {

	common.Logf("SolveMultipleRelationsSingleBinding start %v %v\n", goals, binding)

	bindings := []mentalese.SimpleBinding{binding}

	for _, goal := range goals {
		bindings = solver.SolveSingleRelationMultipleBindings(goal, bindings)
	}

	common.Logf("SolveMultipleRelationsSingleBinding end %v\n", bindings)

	return bindings
}