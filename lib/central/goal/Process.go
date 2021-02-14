package goal

import "nli-go/lib/mentalese"

type Process struct {
	GoalId int
	Stack  []*StackFrame
}

func NewProcess(goalId int, goalSet mentalese.RelationSet) *Process {
	return &Process{
		GoalId: goalId,
		Stack: []*StackFrame{
			NewStackFrame(goalSet, mentalese.InitBindingSet(mentalese.NewBinding())),
		},
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

// advance the cursor in the frame
// pop the frame when done, and transfer child bindings to parent
func (p *Process) Advance() {

	// advance binding
	frame := p.GetLastFrame()
	frame.InBindingIndex++

	// create a new working environment
	frame.Cursor = NewStackFrameCursor()

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
		p.Clear()
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

func (p *Process) CreateMessenger() *Messenger {
	frame := p.GetLastFrame()

	return NewMessenger(
		frame.GetCurrentRelation(),
		frame.GetPreparedBinding(),
		frame.Cursor)
}

func (p *Process) ProcessMessenger(messenger *Messenger, frame *StackFrame) {
	frame.AddOutBindings(frame.GetCurrentInBinding(), messenger.GetOutBindings())

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