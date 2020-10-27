package nested

import (
	"nli-go/lib/mentalese"
)

func (base *SystemNestedStructureBase) SolveNot(notRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	scope := notRelation.Arguments[mentalese.NotScopeIndex].TermValueRelationSet

	newBindings := base.solver.SolveRelationSet(scope, mentalese.InitBindingSet(binding))
	resultBindings := mentalese.NewBindingSet()

	if !newBindings.IsEmpty() {
		resultBindings = mentalese.NewBindingSet()
	} else {
		resultBindings.Add(binding)
	}

	return resultBindings
}

func (base *SystemNestedStructureBase) SolveAnd(andRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	first := andRelation.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := andRelation.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := mentalese.InitBindingSet(binding)

	newBindings = base.solver.SolveRelationSet(first, newBindings)

	if !newBindings.IsEmpty() {
		newBindings = base.solver.SolveRelationSet(second, newBindings)
	}

	return newBindings
}

func (base *SystemNestedStructureBase) SolveOr(orRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	first := orRelation.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := orRelation.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := mentalese.InitBindingSet(binding)

	firstBindings := base.solver.SolveRelationSet(first, newBindings)
	secondBindings := base.solver.SolveRelationSet(second, newBindings)

	result := firstBindings.Copy()
	result.AddMultiple(secondBindings)

	return result
}

func (base *SystemNestedStructureBase) SolveXor(orRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	first := orRelation.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := orRelation.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := base.solver.SolveRelationSet(first, mentalese.InitBindingSet(binding))

	if newBindings.IsEmpty() {
		newBindings = base.solver.SolveRelationSet(second, mentalese.InitBindingSet(binding))
	}

	return newBindings
}


func (base *SystemNestedStructureBase) SolveIfThenElse(ifThenElse mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	condition := ifThenElse.Arguments[0].TermValueRelationSet
	action := ifThenElse.Arguments[1].TermValueRelationSet
	alternative := ifThenElse.Arguments[2].TermValueRelationSet

	newBindings := base.solver.SolveRelationSet(condition, mentalese.InitBindingSet(binding))

	if !newBindings.IsEmpty() {
		newBindings = base.solver.SolveRelationSet(action, newBindings )
	} else {
		newBindings = base.solver.SolveRelationSet(alternative, mentalese.InitBindingSet(binding))
	}

	return newBindings
}
