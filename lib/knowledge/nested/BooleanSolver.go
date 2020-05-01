package nested

import (
	"nli-go/lib/mentalese"
)

func (base *SystemNestedStructureBase) SolveNot(notRelation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	scope := notRelation.Arguments[mentalese.NotScopeIndex].TermValueRelationSet

	newBindings := base.solver.SolveRelationSet(scope, mentalese.Bindings{ binding })
	resultBindings := mentalese.Bindings{}

	if len(newBindings) > 0 {
		resultBindings = mentalese.Bindings{}
	} else {
		resultBindings = mentalese.Bindings{ binding }
	}

	return resultBindings
}

func (base *SystemNestedStructureBase) SolveAnd(andRelation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	first := andRelation.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := andRelation.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

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
