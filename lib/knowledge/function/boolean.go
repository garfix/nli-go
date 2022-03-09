package function

import (
	"nli-go/lib/api"
	"nli-go/lib/mentalese"
)

func (base *SystemSolverFunctionBase) not(messenger api.ProcessMessenger, notRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	scope := notRelation.Arguments[mentalese.NotScopeIndex].TermValueRelationSet

	newBindings := messenger.ExecuteChildStackFrame(scope, mentalese.InitBindingSet(binding))
	if !newBindings.IsEmpty() {
		return mentalese.NewBindingSet()
	} else {
		return mentalese.InitBindingSet(binding)
	}
}

func (base *SystemSolverFunctionBase) and(messenger api.ProcessMessenger, andRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	first := andRelation.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := andRelation.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := mentalese.InitBindingSet(binding)

	childBindings := messenger.ExecuteChildStackFrame(first, newBindings)
	if childBindings.IsEmpty() {
		newBindings = childBindings
	} else {
		newBindings = messenger.ExecuteChildStackFrame(second, childBindings)
	}

	return newBindings
}

func (base *SystemSolverFunctionBase) or(messenger api.ProcessMessenger, orRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	first := orRelation.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := orRelation.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := mentalese.NewBindingSet()

	childBindings := messenger.ExecuteChildStackFrame(first, mentalese.InitBindingSet(binding))
	newBindings.AddMultiple(childBindings)
	childBindings = messenger.ExecuteChildStackFrame(second, mentalese.InitBindingSet(binding))
	newBindings.AddMultiple(childBindings)

	return newBindings
}

func (base *SystemSolverFunctionBase) xor(messenger api.ProcessMessenger, orRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	first := orRelation.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := orRelation.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := mentalese.InitBindingSet(binding)

	childBindings := messenger.ExecuteChildStackFrame(first, mentalese.InitBindingSet(binding))
	if !childBindings.IsEmpty() {
		newBindings = childBindings
	} else {
		childBindings = messenger.ExecuteChildStackFrame(second, mentalese.InitBindingSet(binding))
		newBindings = childBindings
	}

	return newBindings
}
