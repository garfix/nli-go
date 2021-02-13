package goal

import "nli-go/lib/mentalese"

// RelationIndex must always point to a real relation!

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

func (f *StackFrame) GetCurrentInBinding() mentalese.Binding {
	return f.InBindings.Get(f.InBindingIndex)
}

// prepare the active binding to be fed to a function
func (f *StackFrame) GetPreparedBinding() mentalese.Binding {

	binding := f.GetCurrentInBinding()

	// filter out only the variables needed by the relation
	binding = binding.FilterVariablesByName(f.GetCurrentRelation().GetVariableNames())

	return binding
}

// prepare the active binding to be fed to a function
func (f *StackFrame) GetPreparedBindings() mentalese.BindingSet {

	bindings := f.InBindings

	// filter out only the variables needed by the relation
	bindings = bindings.FilterVariablesByName(f.GetCurrentRelation().GetVariableNames())

	return bindings
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