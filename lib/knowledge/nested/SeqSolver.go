package nested

import (
	"nli-go/lib/mentalese"
)

func (base *SystemNestedStructureBase) SolveSeq(seq mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	first := seq.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := seq.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := mentalese.Bindings{binding}

	newBindings = base.solver.SolveRelationSet(first, newBindings)

	if len(newBindings) > 0 {
		newBindings = base.solver.SolveRelationSet(second, newBindings)
	}

	return newBindings
}
