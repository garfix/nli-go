package function

import (
	"nli-go/lib/mentalese"
)

func (base *SystemSolverFunctionBase) not(notRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

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

func (base *SystemSolverFunctionBase) and(andRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	first := andRelation.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := andRelation.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := mentalese.InitBindingSet(binding)

	newBindings = base.solver.SolveRelationSet(first, newBindings)

	if !newBindings.IsEmpty() {
		newBindings = base.solver.SolveRelationSet(second, newBindings)
	}

	return newBindings
}

func (base *SystemSolverFunctionBase) or(orRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	first := orRelation.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := orRelation.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := mentalese.InitBindingSet(binding)

	firstBindings := base.solver.SolveRelationSet(first, newBindings)
	secondBindings := base.solver.SolveRelationSet(second, newBindings)

	result := firstBindings.Copy()
	result.AddMultiple(secondBindings)

	return result
}

func (base *SystemSolverFunctionBase) xor(orRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	first := orRelation.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := orRelation.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := base.solver.SolveRelationSet(first, mentalese.InitBindingSet(binding))

	if newBindings.IsEmpty() {
		newBindings = base.solver.SolveRelationSet(second, mentalese.InitBindingSet(binding))
	}

	return newBindings
}

