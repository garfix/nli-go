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

func (base *SystemNestedStructureBase) SolveOr(orRelation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	first := orRelation.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := orRelation.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := mentalese.Bindings{binding}

	firstBindings := base.solver.SolveRelationSet(first, newBindings)

	if len(firstBindings) != 0 {
		return firstBindings
	}

	secondBindings := base.solver.SolveRelationSet(second, newBindings)

	return secondBindings
}
