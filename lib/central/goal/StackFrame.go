package goal

import "nli-go/lib/mentalese"

type StackFrame struct {
	Relations      mentalese.RelationSet
	RelationIndex  int
	InBindings     mentalese.BindingSet
	InBindingIndex int
	OutBindings    mentalese.BindingSet
	Cursor         *StackFrameCursor
}

func NewStackFrame(relations mentalese.RelationSet, bindings mentalese.BindingSet) *StackFrame {
	return &StackFrame{
		Relations:      relations,
		InBindings:     bindings,
		OutBindings:    mentalese.NewBindingSet(),
		InBindingIndex: 0,
		RelationIndex:  0,
		Cursor:         NewStackFrameCursor(),
	}
}

func (f *StackFrame) IsDone() bool {
	return f.RelationIndex >= len(f.Relations)
}

func (f *StackFrame) GetCurrentRelation() mentalese.Relation {
	return f.Relations[f.RelationIndex]
}

func (f *StackFrame) GetCurrentBinding() mentalese.Binding {
	return f.InBindings.Get(f.InBindingIndex)
}