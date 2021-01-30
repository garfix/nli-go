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

func (p *Process) PushFrame(goalSet mentalese.RelationSet, bindings mentalese.BindingSet) {
	p.Stack = append(p.Stack,
		NewStackFrame(goalSet, bindings))
}

func (p *Process) Clear() {
	p.Stack = []*StackFrame{}
}

// advance the cursor in the frame
// pop the frame when done, and transfer child bindings to parent
func (p *Process) Advance() {
	if !p.IsDone() {

		// advance binding
		frame := p.GetLastFrame()
		frame.BindingIndex++
		if frame.BindingIndex >= frame.Bindings.GetLength() {
			frame.BindingIndex = 0

			// advance position
			frame.Position++
			frame.Cursor = NewStackFrameCursor()

			if frame.IsDone() {
				p.PopFrame()

				// transfer child bindings to parent
				resultBindings := frame.Bindings
				newLastFrame := p.GetLastFrame()
				if newLastFrame != nil {
					newLastFrame.Cursor.ChildFrameResultBindings = resultBindings
				}
			}
		}
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