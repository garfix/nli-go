package central

import (
	"nli-go/lib/mentalese"
)

type Process struct {
	GoalId     string
	Stack      []*StackFrame
	Slots      map[string]mentalese.Term
	WaitingFor mentalese.RelationSet
}

func NewProcess(goalId string, goalSet mentalese.RelationSet, bindings mentalese.BindingSet) *Process {
	return &Process{
		GoalId: goalId,
		Stack: []*StackFrame{
			NewStackFrame(goalSet, bindings),
		},
		Slots:      map[string]mentalese.Term{},
		WaitingFor: nil,
	}
}

func (p *Process) SetWaitingFor(set mentalese.RelationSet) {
	p.WaitingFor = set
}

func (p *Process) GetWaitingFor() mentalese.RelationSet {
	return p.WaitingFor
}

func (p *Process) AddMutableVariable(variable string) {
	for i := len(p.Stack) - 1; i >= 0; i-- {
		frame := p.Stack[i]
		if frame.Cursor.GetType() == mentalese.FrameTypeScope {
			frame.Cursor.AddMutableVariable(variable)
			break
		}
	}
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
		newLastFrame.Cursor.ChildFrameResultBindings = resultBindings
	}
}

// prepare the active binding to be fed to a function
func (p *Process) GetPreparedBinding(f *StackFrame) mentalese.Binding {

	binding := f.GetCurrentInBinding()

	// filter out only the variables needed by the relation
	binding = binding.FilterVariablesByName(f.GetCurrentRelation().GetVariableNames())

	return binding
}

func (p *Process) CreateMessenger(processRunner *ProcessRunner, process *Process) *Messenger {
	frame := p.GetLastFrame()

	return NewMessenger(processRunner, process, frame.Cursor, p.Slots)
}

func (p *Process) ProcessMessenger(messenger *Messenger, currentFrame *StackFrame) *StackFrame {

	outBindings := messenger.GetOutBindings()

	for slot, value := range messenger.newSlots {
		p.Slots[slot] = value
	}

	p.updateMutableVariables(outBindings)

	currentFrame, outBindings = p.executeProcessInstructions(messenger, currentFrame, outBindings)

	currentFrame.AddOutBindings(currentFrame.GetCurrentInBinding(), outBindings)

	if messenger.GetChildFrame() != nil {
		p.PushFrame(messenger.GetChildFrame())
	}

	return currentFrame
}

func (p *Process) executeProcessInstructions(messenger *Messenger, currentFrame *StackFrame, outBindings mentalese.BindingSet) (*StackFrame, mentalese.BindingSet) {

	for instruction, value := range messenger.GetProcessInstructions() {
		switch instruction {
		case mentalese.ProcessInstructionLet:
			p.AddMutableVariable(value)
		case mentalese.ProcessInstructionBreak:
			outBindings = currentFrame.InBindings
			currentFrame = p.executeBreak(currentFrame)
		case mentalese.ProcessInstructionCancel:
			outBindings = mentalese.NewBindingSet()
			currentFrame = p.executeBreak(currentFrame)
		case mentalese.ProcessInstructionReturn:
			outBindings = currentFrame.InBindings
			currentFrame = p.executeReturn(currentFrame)
			//case mentalese.ProcessInstructionSendMessage:
			//	lastFrame := p.GetLastFrame()
			//	binding := lastFrame.InBindings.Get(lastFrame.InBindingIndex)
			//	boundRelations := lastFrame.Relations.BindSingle(binding)
			//	p.messageManager.NotifyListeners(boundRelations)
		}
	}

	return currentFrame, outBindings
}

func (p *Process) executeBreak(currentFrame *StackFrame) *StackFrame {
	done := false
	for !done {
		frame := p.GetLastFrame()
		if frame == nil {
			// todo: log error: break without loop
			done = true
		}

		frameType := frame.Cursor.GetType()

		switch frameType {
		case mentalese.FrameTypeLoop:
			currentFrame = frame
			done = true
		case mentalese.FrameTypeScope:
			// todo: log error: break without loop
			done = true
		default:
			p.PopFrame()
		}
	}

	return currentFrame
}

func (p *Process) executeReturn(currentFrame *StackFrame) *StackFrame {
	for true {
		frame := p.GetLastFrame()

		if frame == nil {
			break
		}
		if frame.Cursor.GetType() == mentalese.FrameTypeScope {
			currentFrame = frame
			break
		} else {
			p.PopFrame()
		}
	}

	return currentFrame
}

func (p *Process) updateMutableVariables(outBindings mentalese.BindingSet) {
	for _, outBinding := range outBindings.GetAll() {
		for variable, value := range outBinding.GetAll() {
			p.updateMutableVariable(variable, value)
		}
	}
}

func (p *Process) updateMutableVariable(variable string, value mentalese.Term) {

	found := false
	for _, frame := range p.Stack {
		if !found {
			// cursor with mutable variable
			if frame.Cursor.HasMutableVariable(variable) {
				frame.Cursor.UpdateMutableVariable(variable, value)
				found = true
			}
		} else {
			// frames below cursor with variable
			frame.UpdateMutableVariable(variable, value)
		}
	}
}

func (p *Process) ProcessMessengerMultipleBindings(messenger *Messenger, frame *StackFrame) {

	// add bindings without variable validation
	frame.OutBindings.AddMultiple(messenger.GetOutBindings())

	// skip the bindings
	frame.InBindingIndex = frame.InBindings.GetLength() - 1

	if messenger.GetChildFrame() != nil {
		p.PushFrame(messenger.GetChildFrame())
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
