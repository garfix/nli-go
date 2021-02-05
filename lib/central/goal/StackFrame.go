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

// prepare the active binding to be fed to a function
func (f *StackFrame) GetInBinding() mentalese.Binding {

	binding := f.GetCurrentBinding()

	// filter out only the variable needed by the relation
	binding = binding.FilterVariablesByName(f.GetCurrentRelation().GetVariableNames())

	return binding
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