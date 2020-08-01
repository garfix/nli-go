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
	secondBindings := base.solver.SolveRelationSet(second, newBindings)

	result := append(firstBindings, secondBindings...)

	return result.UniqueBindings()
}

func (base *SystemNestedStructureBase) SolveXor(orRelation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	first := orRelation.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := orRelation.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := base.solver.SolveRelationSet(first, mentalese.Bindings{ binding })

	if len(newBindings) == 0 {
		newBindings = base.solver.SolveRelationSet(second, mentalese.Bindings{ binding })
	}

	return newBindings
}


func (base *SystemNestedStructureBase) SolveIfThenElse(ifThenElse mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	condition := ifThenElse.Arguments[0].TermValueRelationSet
	action := ifThenElse.Arguments[1].TermValueRelationSet
	alternative := ifThenElse.Arguments[2].TermValueRelationSet

	newBindings := base.solver.SolveRelationSet(condition, mentalese.Bindings{ binding })

	if len(newBindings) > 0 {
		newBindings = base.solver.SolveRelationSet(action, newBindings )
	} else {
		newBindings = base.solver.SolveRelationSet(alternative, mentalese.Bindings{ binding } )
	}

	return newBindings
}
