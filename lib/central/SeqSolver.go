package central

import (
	"nli-go/lib/mentalese"
)

func (solver ProblemSolver) SolveSeq(seq mentalese.Relation, keyCabinet *mentalese.KeyCabinet, binding mentalese.Binding) mentalese.Bindings {

	first := seq.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := seq.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := mentalese.Bindings{binding}

	newBindings = solver.SolveRelationSet(first, keyCabinet, newBindings)

	if len(newBindings) > 0 {
		newBindings = solver.SolveRelationSet(second, keyCabinet, newBindings)
	}

	return newBindings
}
