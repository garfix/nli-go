package example3

type simpleProblemSolver struct {
	sources []*SimpleKnowledgeBase
}

func NewSimpleProblemSolver() *simpleProblemSolver {
	return &simpleProblemSolver{sources: []SimpleKnowledgeBase{}}
}

func (solver *simpleProblemSolver) AddKnowledgeBase(source *SimpleKnowledgeBase) {
	solver.sources = append(solver.sources, source)
}