package central

import (
	"nli-go/lib/api"
	"nli-go/lib/mentalese"
)

type Messenger struct {
	processRunner       *ProcessRunner
	process             *Process
	cursor              *StackFrameCursor
	outBindings         mentalese.BindingSet
	childFrame          *StackFrame
	processInstructions map[string]string
	oldSlots            map[string]mentalese.Term
	newSlots            map[string]mentalese.Term
}

func NewMessenger(processRunner *ProcessRunner, process *Process, cursor *StackFrameCursor, slots map[string]mentalese.Term) *Messenger {
	return &Messenger{
		processRunner:       processRunner,
		process:             process,
		cursor:              cursor,
		outBindings:         mentalese.NewBindingSet(),
		childFrame:          nil,
		processInstructions: map[string]string{},
		oldSlots:            slots,
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

func (i *Messenger) CreateChildStackFrame(relations mentalese.RelationSet, bindings mentalese.BindingSet) {
	i.childFrame = NewStackFrame(relations, bindings)
}

func (i *Messenger) SendMessage(message mentalese.RelationSet) {
	i.processRunner.list.messageManager.NotifyListeners(message)
}

func (i *Messenger) ExecuteChildStackFrame(relations mentalese.RelationSet, bindings mentalese.BindingSet) (mentalese.BindingSet, bool) {

	if len(relations) == 0 {
		return bindings, false
	}

	newBindings := i.processRunner.PushAndRun(i.process, relations, bindings)

	return newBindings, false
}

func (i *Messenger) StartProcess(relations mentalese.RelationSet, binding mentalese.Binding) {
	i.processRunner.StartProcess(relations, binding)
}

func (i *Messenger) GetProcessSlot(slot string) (mentalese.Term, bool) {
	value, found := i.oldSlots[slot]
	return value, found
}

func (i *Messenger) SetProcessSlot(slot string, value mentalese.Term) {
	i.oldSlots[slot] = value
	i.newSlots[slot] = value
}

func (i *Messenger) GetChildFrame() *StackFrame {
	return i.childFrame
}

func (i *Messenger) GetOutBindings() mentalese.BindingSet {
	return i.outBindings
}