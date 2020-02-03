package central

import (
	"nli-go/lib/mentalese"
)

func (solver ProblemSolver) SolveSeq(seq mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	first := seq.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := seq.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := mentalese.Bindings{binding}

	newBindings = solver.SolveRelationSet(first, newBindings)

	if len(newBindings) > 0 {
		newBindings = solver.SolveRelationSet(second, newBindings)
	}

	return newBindings
}
