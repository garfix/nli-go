package function

import (
	"nli-go/lib/api"
	"nli-go/lib/mentalese"
)

func (base *SystemSolverFunctionBase) not(messenger api.ProcessMessenger, notRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	scope := notRelation.Arguments[mentalese.NotScopeIndex].TermValueRelationSet

	//cursor := messenger.GetCursor()
	//state := cursor.GetState("state", 0)
	//
	//if state == 0 {
	//	cursor.SetState("state", 1)
	//	messenger.CreateChildStackFrame(scope, mentalese.InitBindingSet(binding))
	//	return mentalese.NewBindingSet()
	//} else {
	//	newBindings := cursor.GetChildFrameResultBindings()
	//	resultBindings := mentalese.NewBindingSet()
	//	if !newBindings.IsEmpty() {
	//		resultBindings = mentalese.NewBindingSet()
	//	} else {
	//		resultBindings.Add(binding)
	//	}
	//	return resultBindings
	//}

	newBindings, _ := messenger.ExecuteChildStackFrameAsync(scope, mentalese.InitBindingSet(binding))
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

	//cursor := messenger.GetCursor()
	//state := cursor.GetState("state", 0)

	//if state == 0 {
	//	cursor.SetState("state", 1)
	//	messenger.CreateChildStackFrame(first, mentalese.InitBindingSet(binding))
	//	return mentalese.NewBindingSet()
	//} else if state == 1 {
	//	cursor.SetState("state", 2)
	//	childBindings := cursor.GetChildFrameResultBindings()
	//	if childBindings.IsEmpty() {
	//		return childBindings
	//	}
	//	messenger.CreateChildStackFrame(second, childBindings)
	//	return mentalese.NewBindingSet()
	//} else {
	//	childBindings := cursor.GetChildFrameResultBindings()
	//	newBindings = childBindings
	//}

	childBindings, _ := messenger.ExecuteChildStackFrameAsync(first, newBindings)
	if childBindings.IsEmpty() {
		newBindings = childBindings
	} else {
		newBindings, _ = messenger.ExecuteChildStackFrameAsync(second, childBindings)
	}

	return newBindings
}

func (base *SystemSolverFunctionBase) or(messenger api.ProcessMessenger, orRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	first := orRelation.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := orRelation.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := mentalese.NewBindingSet()

	//cursor := messenger.GetCursor()
	//state := cursor.GetState("state", 0)
	//
	//if state == 0 {
	//	cursor.SetState("state", 1)
	//	messenger.CreateChildStackFrame(first, mentalese.InitBindingSet(binding))
	//	return mentalese.NewBindingSet()
	//} else if state == 1 {
	//	cursor.SetState("state", 2)
	//	childBindings := cursor.GetChildFrameResultBindings()
	//	cursor.AddStepBindings(childBindings)
	//	messenger.CreateChildStackFrame(second, mentalese.InitBindingSet(binding))
	//	return mentalese.NewBindingSet()
	//} else {
	//	childBindings := cursor.GetChildFrameResultBindings()
	//	cursor.AddStepBindings(childBindings)
	//	for _, childBindings := range cursor.GetAllStepBindings() {
	//		result.AddMultiple(childBindings)
	//	}
	//}

	childBindings, _ := messenger.ExecuteChildStackFrameAsync(first, mentalese.InitBindingSet(binding))
	newBindings.AddMultiple(childBindings)
	childBindings, _ = messenger.ExecuteChildStackFrameAsync(second, mentalese.InitBindingSet(binding))
	newBindings.AddMultiple(childBindings)

	return newBindings
}

func (base *SystemSolverFunctionBase) xor(messenger api.ProcessMessenger, orRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	first := orRelation.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet
	second := orRelation.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet

	newBindings := mentalese.InitBindingSet(binding)

	//cursor := messenger.GetCursor()
	//state := cursor.GetState("state", 0)
	//
	//if state == 0 {
	//	cursor.SetState("state", 1)
	//	messenger.CreateChildStackFrame(first, mentalese.InitBindingSet(binding))
	//	return mentalese.NewBindingSet()
	//} else if state == 1 {
	//	cursor.SetState("state", 2)
	//	childBindings := cursor.GetChildFrameResultBindings()
	//	if !childBindings.IsEmpty() {
	//		return childBindings
	//	}
	//	messenger.CreateChildStackFrame(second, mentalese.InitBindingSet(binding))
	//	return mentalese.NewBindingSet()
	//} else {
	//	childBindings := cursor.GetChildFrameResultBindings()
	//	newBindings = childBindings
	//}

	childBindings, _ := messenger.ExecuteChildStackFrameAsync(first, mentalese.InitBindingSet(binding))
	if !childBindings.IsEmpty() {
		newBindings = childBindings
	} else {
		childBindings, _ = messenger.ExecuteChildStackFrameAsync(second, mentalese.InitBindingSet(binding))
		newBindings = childBindings
	}

	return newBindings
}
