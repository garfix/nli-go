package central

import (
	"nli-go/lib/api"
	"nli-go/lib/mentalese"
)

type Messenger struct {
	processRunner       *ProcessRunner
	process             *Process
	cursor              *mentalese.StackFrameCursor
	outBindings         mentalese.BindingSet
	processInstructions map[string]string
	oldSlots            map[string]mentalese.Term
	newSlots            map[string]mentalese.Term
}

func NewMessenger(processRunner *ProcessRunner, process *Process, cursor *mentalese.StackFrameCursor, slots map[string]mentalese.Term) *Messenger {
	return &Messenger{
		processRunner:       processRunner,
		process:             process,
		cursor:              cursor,
		outBindings:         mentalese.NewBindingSet(),
		processInstructions: map[string]string{},
		oldSlots:            slots,
		newSlots:            map[string]mentalese.Term{},
	}
}

func NewSimpleMessenger() *Messenger {
	return &Messenger{
		processRunner:       nil,
		process:             nil,
		cursor:              nil,
		outBindings:         mentalese.NewBindingSet(),
		processInstructions: map[string]string{},
		oldSlots:            nil,
		newSlots:            map[string]mentalese.Term{},
	}
}

func (i *Messenger) GetProcess() api.Process {
	return i.process
}

func (i *Messenger) GetCursor() api.ProcessCursor {
	return i.cursor
}

func (i *Messenger) AddProcessInstruction(name string, value string) {
	i.processInstructions[name] = value
}

func (i *Messenger) GetProcessInstructions() map[string]string {
	return i.processInstructions
}

func (i *Messenger) AddOutBinding(binding mentalese.Binding) {
	i.outBindings.Add(binding)
}

func (i *Messenger) AddOutBindings(bindings mentalese.BindingSet) {
	i.outBindings.AddMultiple(bindings)
}

// executes `relations` expecting immutable variables
// modifies only pre-existing mutable variables
func (i *Messenger) ExecuteChildStackFrame(relations mentalese.RelationSet, bindings mentalese.BindingSet) mentalese.BindingSet {

	// mark the calling function as non-plain;
	// checking for plain-ness is done when bindings of mutable variables need to be stored in the scope
	// function `ProcessMessenger`
	if i.cursor.GetType() == mentalese.FrameTypePlain {
		i.cursor.SetType(mentalese.FrameTypeComplex)
	}

	// when a loop has been breaked, the remaining calls are nullified
	if i.cursor.GetState() == mentalese.StateInterrupted {
		return mentalese.NewBindingSet()
	}

	if len(relations) == 0 {
		return bindings
	}

	if bindings.GetLength() == 0 {
		return bindings
	}

	// save the old state of mutable variables
	saved := i.GetMutableVariables()

	newBindings := i.processRunner.PushAndRun(i.process, relations, bindings)

	// update saved values with result
	for key, value := range i.GetMutableVariables().GetAll() {
		_, found := saved.Get(key)
		if found {
			saved.Set(key, value)
		}
	}
	i.SetMutableVariables(saved)

	return newBindings
}

// as `ExecuteChildStackFrame`, but expecting and returning mutable variables
// relations: relations with mutable variables
// return: mutable variables
func (i *Messenger) ExecuteChildStackFrameMutable(relationsMut mentalese.RelationSet, binding mentalese.Binding) mentalese.BindingSet {

	relationsIm := relationsMut.ConvertVariablesToImmutables()

	// create an immutable variable binding from mutable variables
	bindingIm := i.collectImmutableVariablesForRelations(relationsMut).ConvertVariablesToImmutables()
	newBindings := i.ExecuteChildStackFrame(relationsIm, mentalese.InitBindingSet(bindingIm))

	return newBindings.ConvertVariablesToMutables()
}

func (i *Messenger) collectImmutableVariablesForRelations(relations mentalese.RelationSet) mentalese.Binding {
	binding := mentalese.NewBinding()
	variables := relations.GetVariableNames()
	for _, variable := range variables {
		value, found := i.GetMutableVariable(variable)
		if found {
			binding.Set(variable, value)
		}
	}
	return binding
}

func (i *Messenger) StartProcess(resource string, relations mentalese.RelationSet, binding mentalese.Binding) bool {
	return i.processRunner.StartProcess(resource, relations, binding)
}

func (i *Messenger) GetProcessSlot(slot string) (mentalese.Term, bool) {
	value, found := i.oldSlots[slot]
	return value, found
}

func (i *Messenger) SetProcessSlot(slot string, value mentalese.Term) {
	i.oldSlots[slot] = value
	i.newSlots[slot] = value
}

func (i *Messenger) GetOutBindings() mentalese.BindingSet {
	return i.outBindings
}

func (p *Messenger) SetMutableVariable(variable string, value mentalese.Term) {
	scope := p.process.GetCurrentScope()
	if scope != nil {
		scope.Cursor.MutableVariableValues.Set(variable, value)
	}
}

func (p *Messenger) GetMutableVariable(variable string) (mentalese.Term, bool) {
	scope := p.process.GetCurrentScope()
	if scope != nil {
		value, found := scope.Cursor.MutableVariableValues.Get(variable)
		return value, found
	}
	return mentalese.Term{}, false
}

func (p *Messenger) SetMutableVariables(binding mentalese.Binding) {
	scope := p.process.GetCurrentScope()
	if scope != nil {
		scope.Cursor.MutableVariableValues = binding
	}
}

func (p *Messenger) GetMutableVariables() mentalese.Binding {
	scope := p.process.GetCurrentScope()
	if scope != nil {
		return scope.Cursor.MutableVariableValues.Copy()
	}
	return mentalese.NewBinding()
}
