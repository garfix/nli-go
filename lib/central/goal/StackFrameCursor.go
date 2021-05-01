package goal

import (
	"nli-go/lib/mentalese"
)

type StackFrameCursor struct {
	Type                     string
	MutableVariables 		 map[string]bool
	State                    map[string]int
	AllStepBindings          []mentalese.BindingSet
	ChildFrameResultBindings mentalese.BindingSet
}

func NewStackFrameCursor() *StackFrameCursor {
	return &StackFrameCursor{
		Type: 					  mentalese.FrameTypePlain,
		MutableVariables: 		  map[string]bool{},
		State:                    map[string]int{},
		AllStepBindings:          []mentalese.BindingSet{},
		ChildFrameResultBindings: mentalese.NewBindingSet(),
	}
}

func (c *StackFrameCursor) HasMutableVariable(variable string) bool {
	_, found := c.MutableVariables[variable]
	return found
}

func (c *StackFrameCursor) AddMutableVariable(variable string) {
	c.MutableVariables[variable] = true
}

func (c *StackFrameCursor) UpdateMutableVariable(variable string, value mentalese.Term) {
	for _, bindingSet := range c.AllStepBindings {
		for _, binding := range bindingSet.GetAll() {
			if binding.ContainsVariable(variable) {
				binding.Set(variable, value)
			}
		}
	}
	for _, binding := range c.ChildFrameResultBindings.GetAll() {
		if binding.ContainsVariable(variable) {
			binding.Set(variable, value)
		}
	}
}

func (c *StackFrameCursor) SetType(t string) {
	c.Type = t
}

func (c *StackFrameCursor) GetType() string {
	return c.Type
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

func (c *StackFrameCursor) GetAllStepBindings() []mentalese.BindingSet {
	return c.AllStepBindings
}

func (c *StackFrameCursor) AddStepBindings(bindings mentalese.BindingSet) {
	c.AllStepBindings = append(c.AllStepBindings, bindings)
}

func (c *StackFrameCursor) GetChildFrameResultBindings() mentalese.BindingSet {
	return c.ChildFrameResultBindings
}