package goal

import (
	"nli-go/lib/api"
	"nli-go/lib/mentalese"
)

type Messenger struct {
	cursor *StackFrameCursor
	outBindings mentalese.BindingSet
	childFrame *StackFrame
}

func NewMessenger(relation mentalese.Relation, binding mentalese.Binding, cursor *StackFrameCursor) *Messenger {
	return &Messenger{
		cursor: cursor,
		outBindings: mentalese.NewBindingSet(),
		childFrame: nil,
	}
}

func (i *Messenger) GetCursor() api.ProcessCursor {
	return i.cursor
}

func (i *Messenger) AddOutBinding(binding mentalese.Binding) {
	i.outBindings.Add(binding)
}

func (i *Messenger) AddOutBindings(bindings mentalese.BindingSet) {
	i.outBindings.AddMultiple(bindings)
}

func (i *Messenger) CreateChildStackFrame(relations mentalese.RelationSet, bindings mentalese.BindingSet) {

	if relations.IsEmpty() {
		panic("Cannot create stack frame with no relations")
	}

	i.childFrame = NewStackFrame(relations, bindings)
}

func (i *Messenger) GetChildFrame() *StackFrame {
	return i.childFrame
}

func (i *Messenger) GetOutBindings() mentalese.BindingSet {
	return i.outBindings
}
