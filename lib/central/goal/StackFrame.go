package goal

import "nli-go/lib/mentalese"

type StackFrame struct {
	Relations mentalese.RelationSet
	Bindings  mentalese.BindingSet
	BindingIndex int
	Position int
	Cursor *StackFrameCursor
}

func NewStackFrame(relations mentalese.RelationSet, bindings mentalese.BindingSet) *StackFrame {
	return &StackFrame{
		Relations: relations,
		Bindings:  bindings,
		BindingIndex: 0,
		Position: 0,
		Cursor: NewStackFrameCursor(),
	}
}

func (f *StackFrame) IsDone() bool {
	return f.Position >= len(f.Relations)
}

func (f *StackFrame) GetCurrentRelation() mentalese.Relation {
	return f.Relations[f.Position]
}

func (f *StackFrame) GetCurrentBinding() mentalese.Binding {
	return f.Bindings.Get(f.BindingIndex)
}