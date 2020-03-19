package central

import (
	"nli-go/lib/mentalese"
)

func (solver ProblemSolver) Call(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	child := relation.Arguments[0].TermValueRelationSet

	newBindings := solver.SolveRelationSet(child, mentalese.Bindings{ binding })

	return newBindings
}
