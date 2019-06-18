package central

import (
	"nli-go/lib/mentalese"
)

func (solver ProblemSolver) SolveSeq(seq mentalese.Relation, nameStore *mentalese.ResolvedNameStore, binding mentalese.Binding) []mentalese.Binding {

	first := seq.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := seq.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := []mentalese.Binding{binding}

	newBindings = solver.SolveRelationSet(first, nameStore, newBindings)

	if len(newBindings) > 0 {
		newBindings = solver.SolveRelationSet(second, nameStore, newBindings)
	}

	return newBindings
}
