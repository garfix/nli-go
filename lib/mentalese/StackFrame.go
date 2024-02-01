package mentalese

import "fmt"

// RelationIndex must always point to a real relation!

type StackFrame struct {
	Relations      RelationSet
	RelationIndex  int
	InBindings     BindingSet
	InBindingIndex int
	HandlerIndex   int
	HandlerCount   int
	OutBindings    BindingSet
	Cursor         *StackFrameCursor
}

func NewStackFrame(relations RelationSet, bindings BindingSet) *StackFrame {
	return &StackFrame{
		Relations:      relations,
		InBindings:     bindings,
		OutBindings:    NewBindingSet(),
		HandlerCount:   0,
		InBindingIndex: 0,
		HandlerIndex:   0,
		RelationIndex:  0,
		Cursor:         NewStackFrameCursor(),
	}
}

func (f *StackFrame) IsDone() bool {
	return f.RelationIndex >= len(f.Relations)
}

func (f *StackFrame) GetCurrentRelation() Relation {
	return f.Relations[f.RelationIndex]
}

func (f *StackFrame) GetCurrentInBinding() Binding {
	return f.InBindings.Get(f.InBindingIndex)
}

func (f *StackFrame) AddOutBinding(outBinding Binding) {
	f.OutBindings.Add(outBinding)
}

func (f *StackFrame) AsId() string {
	return fmt.Sprintf("%p-%d-%d-%d", f, f.RelationIndex, f.HandlerIndex, f.InBindingIndex)
}
