package goal

import "nli-go/lib/mentalese"

// RelationIndex must always point to a real relation!

type StackFrame struct {
	Relations      mentalese.RelationSet
	RelationIndex  int
	InBindings     mentalese.BindingSet
	InBindingIndex int
	HandlerIndex   int
	HandlerCount   int
	OutBindings    mentalese.BindingSet
	Cursor         *StackFrameCursor
}

func NewStackFrame(relations mentalese.RelationSet, bindings mentalese.BindingSet) *StackFrame {
	return &StackFrame{
		Relations:      relations,
		InBindings:     bindings,
		OutBindings:    mentalese.NewBindingSet(),
		HandlerCount: 0,
		InBindingIndex: 0,
		HandlerIndex: 0,
		RelationIndex:  0,
		Cursor:         NewStackFrameCursor(),
	}
}

func (f *StackFrame) UpdateMutableVariable(variable string, value mentalese.Term) {
	for _, binding := range f.InBindings.GetAll() {
		if binding.ContainsVariable(variable) {
			binding.Set(variable, value)
		}
	}
	for _, binding := range f.OutBindings.GetAll() {
		if binding.ContainsVariable(variable) {
			binding.Set(variable, value)
		}
	}
	f.Cursor.UpdateMutableVariable(variable, value)
}

func (f *StackFrame) IsDone() bool {
	return f.RelationIndex >= len(f.Relations)
}

func (f *StackFrame) GetCurrentRelation() mentalese.Relation {
	return f.Relations[f.RelationIndex]
}

func (f *StackFrame) GetCurrentInBinding() mentalese.Binding {
	return f.InBindings.Get(f.InBindingIndex)
}

func (f *StackFrame) AddOutBinding(inBinding mentalese.Binding, outBinding mentalese.Binding) {

	relationVariables := f.GetCurrentRelation().GetVariableNames()

	// filter out temporary variables
	cleanBinding := outBinding.FilterVariablesByName(relationVariables)
	// make sure the original values are present
	cleanBinding = cleanBinding.Merge(inBinding)

	f.OutBindings.Add(cleanBinding)
}

func (f *StackFrame) AddOutBindings(inBinding mentalese.Binding, outBindings mentalese.BindingSet) {
	for _, outBinding := range outBindings.GetAll() {
		f.AddOutBinding(inBinding, outBinding)
	}
}