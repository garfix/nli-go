package mentalese

const StateOk = "ok"
const StateInterrupted = "interrupted"

type StackFrameCursor struct {
	Type                     string
	MutableVariableValues    Binding
	ChildFrameResultBindings BindingSet
	State                    string
}

func NewStackFrameCursor() *StackFrameCursor {
	return &StackFrameCursor{
		Type:                     FrameTypePlain,
		MutableVariableValues:    NewBinding(),
		ChildFrameResultBindings: NewBindingSet(),
		State:                    StateOk,
	}
}

func (c *StackFrameCursor) SetState(phase string) {
	c.State = phase
}

func (c *StackFrameCursor) GetState() string {
	return c.State
}

func (c *StackFrameCursor) SetType(t string) {
	c.Type = t
}

func (c *StackFrameCursor) GetType() string {
	return c.Type
}
