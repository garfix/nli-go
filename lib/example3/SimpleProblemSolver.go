package example3

type simpleProblemSolver struct {
	sources []SimpleKnowledgeBase
	matcher *simpleRelationMatcher
}

func NewSimpleProblemSolver() *simpleProblemSolver {
	return &simpleProblemSolver{sources: []*SimpleKnowledgeBase{}, matcher:NewSimpleRelationMatcher()}
}

func (solver *simpleProblemSolver) AddKnowledgeBase(source SimpleKnowledgeBase) {
	solver.sources = append(solver.sources, source)
}

// goals e.g. { father(X, Y), father(Y, Z)}
// return e.g. {
//  { father('john', 'jack'), father('jack', 'joe') }
//  { father('bob', 'jonathan'), father('jonathan', 'bill') }
// }
func (solver simpleProblemSolver) Solve(goals SimpleRelationSet) [][]SimpleRelation {

	bindings := solver.SolveMultipleRelations(goals.GetRelations())
	solutions := solver.matcher.bindMultipleRelationsMultipleBindings(goals, bindings)

	return solutions
}

// goals e.g. { father(X, Y), father(Y, Z)}
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver simpleProblemSolver) SolveMultipleRelations(goals []SimpleRelation) []SimpleBinding {

	bindings := []SimpleBinding{}

	for _, goal := range goals {
		bindings = solver.SolveSingleRelationMultipleBindings(goal, bindings)
	}

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
func (solver simpleProblemSolver) SolveSingleRelationMultipleBindings(goalRelation SimpleRelation, bindings []SimpleBinding) []SimpleBinding {

	newBindings := []SimpleBinding{}

	for _, binding := range bindings {
		newBindings = append(newBindings, solver.SolveSingleRelationSingleBinding(goalRelation, binding))
	}

	return newBindings
}

// goal e.g. father(Y, Z)
// bindings e.g. {
//  { {X='john', Y='jack'} }
//  { {X='bob', Y='jonathan'} }
// }
// return e.g. {X='bob', Y='jonathan', Z='bill'}
func (solver simpleProblemSolver) SolveSingleRelationSingleBinding(goalRelation SimpleRelation, binding SimpleBinding) SimpleBinding {

	bindings := []SimpleBinding{}

	// go through all knowledge sources
	for _, source := range solver.sources {
		bindings = append(bindings, source.Bind(goalRelation))
	}

	return bindings
}