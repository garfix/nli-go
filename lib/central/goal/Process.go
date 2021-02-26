package goal

import "nli-go/lib/mentalese"

type Process struct {
	GoalId           int
	Stack            []*StackFrame
	MutableVariables map[string]bool
}

func NewProcess(goalId int, goalSet mentalese.RelationSet) *Process {
	return &Process{
		GoalId: goalId,
		Stack: []*StackFrame{
			NewStackFrame(goalSet, mentalese.InitBindingSet(mentalese.NewBinding())),
		},
		MutableVariables: map[string]bool{},
	}
}

func (p *Process) AddMutableVariable(variable string) {
	p.MutableVariables[variable] = true
}

func (p* Process) IsMutableVariable(variable string) bool {
	_, found := p.MutableVariables[variable]
	return found
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

func (p *Process) CreateMessenger() *Messenger {
	frame := p.GetLastFrame()

	return NewMessenger(frame.Cursor)
}

func (p *Process) ProcessMessenger(messenger *Messenger, frame *StackFrame) {

	outBindings := messenger.GetOutBindings()

	p.executeProcessInstructions(messenger, frame)

	p.updateMutableVariables(outBindings)

	frame.AddOutBindings(frame.GetCurrentInBinding(), outBindings)

	if messenger.GetChildFrame() != nil {
		p.PushFrame(messenger.GetChildFrame())
	}
}

func (p *Process) executeProcessInstructions(messenger *Messenger, frame *StackFrame) {

	for instruction, value := range messenger.GetProcessInstructions() {
		switch instruction {
		case mentalese.ProcessInstructionLet:
			p.AddMutableVariable(value)
		case mentalese.ProcessInstructionType:
			frame.SetType(value)
		case mentalese.ProcessInstructionBreak:
			p.executeBreak()
		}
	}
}

func (p *Process) executeBreak() {
	done := false
	for !done {
		frame := p.GetLastFrame()
		if frame == nil {
			// todo: log error: break without loop
			done = true
		}
		frameType := frame.GetType()

		p.PopFrame()

		if frameType == mentalese.FrameTypeLoop {
			done = true
		}
		if frameType == mentalese.FrameTypeScope {
			// todo: log error: break without loop
			done = true
		}
	}
}

func (p *Process) updateMutableVariables(outBindings mentalese.BindingSet) {
	for _, outBinding := range outBindings.GetAll() {
		for variable, value := range outBinding.GetAll() {
			if p.IsMutableVariable(variable) {
				p.updateMutableVariable(variable, value)
			}
		}
	}
}

func (p *Process) updateMutableVariable(variable string, value mentalese.Term) {

	for _, frame := range p.Stack {
		frame.UpdateMutableVariable(variable, value)
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
		return p.Stack[len(p.Stack) - 1]
	}
}

func (p *Process) PopFrame() {
	p.Stack = p.Stack[0 : len(p.Stack) - 1]
}

func (p *Process) IsDone() bool {
	return len(p.Stack) == 0
}