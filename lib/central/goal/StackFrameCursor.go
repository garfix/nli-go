package goal

import "nli-go/lib/mentalese"

type StackFrameCursor struct {
	State                    map[string]int
	StepBindings             mentalese.BindingSet
	ChildFrameResultBindings mentalese.BindingSet
}

func NewStackFrameCursor() *StackFrameCursor {
	return &StackFrameCursor{
		State:                    map[string]int{},
		StepBindings:             mentalese.NewBindingSet(),
		ChildFrameResultBindings: mentalese.NewBindingSet(),
	}
}

func (c *StackFrameCursor) GetState(name string, fallback int) int {
	value, found := c.State[name]
	if found {
		return value
	} else {
		return fallback
	}
}

func (c *StackFrameCursor) SetState(name string, value int)  {
	c.State[name] = value
}

func (c *StackFrameCursor) GetStepBindings() mentalese.BindingSet {
	return c.StepBindings
}

func (c *StackFrameCursor) AddStepBinding(binding mentalese.Binding) {
	c.StepBindings.Add(binding)
}

func (c *StackFrameCursor) AddStepBindings(bindings mentalese.BindingSet) {
	c.StepBindings.AddMultiple(bindings)
}

func (c *StackFrameCursor) GetChildFrameResultBindings() mentalese.BindingSet {
	return c.ChildFrameResultBindings
}