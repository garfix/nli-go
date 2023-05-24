package central

import (
	"nli-go/lib/mentalese"
)

type Process struct {
	ProcessType string
	GoalId      string
	Stack       []*StackFrame
	Slots       map[string]mentalese.Term
	channel     chan mentalese.Request
}

func NewProcess(processType string, goalId string, goalSet mentalese.RelationSet, bindings mentalese.BindingSet) *Process {
	return &Process{
		ProcessType: processType,
		GoalId:      goalId,
		Stack: []*StackFrame{
			NewStackFrame(goalSet, bindings),
		},
		Slots:   map[string]mentalese.Term{},
		channel: make(chan mentalese.Request),
	}
}

func (p *Process) GetType() string {
	return p.ProcessType
}

func (p *Process) GetChannel() chan mentalese.Request {
	return p.channel
}

func (p *Process) PushFrame(frame *StackFrame) {
	p.Stack = append(p.Stack, frame)
}

func (p *Process) Clear() {
	p.Stack = []*StackFrame{}
}

func (p *Process) EmptyRelationCheck() {
	frame := p.GetLastFrame()
	if frame.Relations.IsEmpty() {
		p.advanceFrame(frame)
	}
}

func (p *Process) Advance() {

	frame := p.GetLastFrame()
	frame.Cursor = NewStackFrameCursor()

	p.AdvanceHandler()
}

func (p *Process) AdvanceHandler() {

	frame := p.GetLastFrame()
	frame.HandlerIndex++

	if frame.HandlerIndex >= frame.HandlerCount {
		p.AdvanceBinding()
	}
}

func (p *Process) AdvanceBinding() {

	frame := p.GetLastFrame()
	frame.HandlerIndex = 0
	frame.InBindingIndex++

	if frame.InBindingIndex >= frame.InBindings.GetLength() {
		p.advanceRelation(frame)
	}
}

func (p *Process) advanceRelation(frame *StackFrame) {

	frame.InBindings = frame.OutBindings
	frame.InBindingIndex = 0

	frame.OutBindings = mentalese.NewBindingSet()

	frame.RelationIndex++

	if frame.InBindings.IsEmpty() {
		// process failed due to no result bindings
		p.advanceFrame(frame)
	} else if frame.IsDone() {
		p.advanceFrame(frame)
	}
}

func (p *Process) advanceFrame(frame *StackFrame) {

	p.PopFrame()

	// transfer child bindings to parent
	resultBindings := frame.InBindings
	newLastFrame := p.GetLastFrame()
	if newLastFrame != nil {
		newLastFrame.Cursor.ChildFrameResultBindings.AddMultiple(resultBindings)
	}
}

// prepare the active binding to be fed to a function
func (p *Process) GetPreparedBinding(f *StackFrame) mentalese.Binding {

	binding := f.GetCurrentInBinding()
	relation := f.GetCurrentRelation()

	binding = p.addMutableVariables(relation, binding)

	// filter out only the variables needed by the relation
	binding = binding.FilterVariablesByName(relation.GetVariableNames())

	return binding
}

func (p *Process) addMutableVariables(relation mentalese.Relation, binding mentalese.Binding) mentalese.Binding {
	scope := p.GetContextScope()
	if scope != nil {
		mutables := scope.Cursor.MutableVariableValues
		filtered := mutables.FilterVariablesByName(relation.GetVariableNames())
		return binding.Merge(filtered)
	} else {
		return binding
	}
}

func (p *Process) AddMutableVariablesMultiple(relation mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {
	scope := p.GetContextScope()
	if scope != nil {
		mutables := scope.Cursor.MutableVariableValues
		filtered := mutables.FilterVariablesByName(relation.GetVariableNames())
		newBindings := mentalese.NewBindingSet()
		for _, binding := range bindings.GetAll() {
			newBindings.Add(binding.Merge(filtered))
		}
		return newBindings
	} else {
		return bindings
	}
}

func (p *Process) CreateMessenger(processRunner *ProcessRunner, process *Process) *Messenger {
	frame := p.GetLastFrame()

	return NewMessenger(processRunner, process, frame.Cursor, p.Slots)
}

func (p *Process) ProcessMessenger(messenger *Messenger, frame *StackFrame) *StackFrame {

	for slot, value := range messenger.newSlots {
		p.Slots[slot] = value
	}

	outBindings := messenger.GetOutBindings()
	relationVariables := frame.GetCurrentRelation().GetVariableNames()

	if messenger.cursor.GetType() == mentalese.FrameTypePlain {
		mutableBindings := outBindings.FilterMutableVariables()
		for k, v := range mutableBindings.FilterVariablesByName(relationVariables).GetAll() {
			messenger.SetMutableVariable(k, v)
		}
	}

	outBindingsWithoutMutables := outBindings.RemoveMutableVariables()

	processedOutBindings := mentalese.NewBindingSet()
	for _, outBinding := range outBindingsWithoutMutables.GetAll() {
		// filter out temporary variables
		cleanBinding := outBinding.FilterVariablesByName(relationVariables)
		// make sure the original values are present
		cleanBinding = cleanBinding.Merge(frame.GetCurrentInBinding())
		processedOutBindings.Add(cleanBinding)
	}

	outBindings = p.executeProcessInstructions(messenger, processedOutBindings)

	frame.OutBindings.AddMultiple(outBindings)

	return frame
}

func (p *Process) ProcessMessengerMultipleBindings(messenger *Messenger, frame *StackFrame) {

	outBindings := messenger.GetOutBindings()
	outBindingsWithoutMutables := outBindings.RemoveMutableVariables()

	// add bindings without variable validation
	frame.OutBindings.AddMultiple(outBindingsWithoutMutables)

	// skip all bindings
	frame.InBindingIndex = frame.InBindings.GetLength() - 1
}

func (p *Process) executeProcessInstructions(messenger *Messenger, outBindings mentalese.BindingSet) mentalese.BindingSet {

	for instruction := range messenger.GetProcessInstructions() {
		switch instruction {
		case mentalese.ProcessInstructionBreak:
			p.executeBreak(outBindings)
			outBindings = mentalese.NewBindingSet()
		case mentalese.ProcessInstructionCancel:
			p.executeBreak(outBindings)
			outBindings = mentalese.NewBindingSet()
		case mentalese.ProcessInstructionReturn:
			p.executeReturn(outBindings)
			outBindings = mentalese.NewBindingSet()
		}
	}

	return outBindings
}

func (p *Process) GetCurrentScope() *StackFrame {

	var scope *StackFrame = nil
	i := len(p.Stack) - 1

	for i >= 0 {
		frame := p.Stack[i]

		if frame.Cursor.GetType() == mentalese.FrameTypeScope {
			return frame
		}

		i--
	}

	return scope
}

func (p *Process) GetContextScope() *StackFrame {

	var scope *StackFrame = nil
	i := len(p.Stack) - 2

	for i >= 0 {
		frame := p.Stack[i]

		if frame.Cursor.GetType() == mentalese.FrameTypeScope {
			return frame
		}

		i--
	}

	return scope
}

func (p *Process) executeBreak(bindings mentalese.BindingSet) {
	done := false
	i := len(p.Stack) - 1

	for !done && i >= 0 {

		frame := p.Stack[i]
		frameType := frame.Cursor.GetType()
		frame.Cursor.SetState(StateInterrupted)

		if frameType == mentalese.FrameTypeLoop {
			frame.Cursor.ChildFrameResultBindings.AddMultiple(bindings)
			done = true
		}

		i--
	}
}

func (p *Process) executeReturn(bindings mentalese.BindingSet) {
	done := false
	i := len(p.Stack) - 1

	for !done && i >= 0 {
		frame := p.Stack[i]

		if frame.Cursor.GetType() == mentalese.FrameTypeScope {
			frame.Cursor.ChildFrameResultBindings.AddMultiple(bindings)
			done = true
		} else {
			frame.Cursor.SetState(StateInterrupted)
		}

		i--
	}
}

func (p *Process) GetCursor() *StackFrameCursor {
	frame := p.GetLastFrame()
	if frame == nil {
		return nil
	}
	return frame.Cursor
}

func (p *Process) GetLastFrame() *StackFrame {
	if len(p.Stack) == 0 {
		return nil
	} else {
		frame := p.Stack[len(p.Stack)-1]
		return frame
	}
}

func (p *Process) GetBeforeLastFrame() *StackFrame {
	if len(p.Stack) < 2 {
		return nil
	} else {
		return p.Stack[len(p.Stack)-2]
	}
}

func (p *Process) PopFrame() {
	p.Stack = p.Stack[0 : len(p.Stack)-1]
}

func (p *Process) IsDone() bool {
	return len(p.Stack) == 0
}
