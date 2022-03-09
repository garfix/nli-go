package central

import (
	"nli-go/lib/mentalese"
)

const PhaseOk = "ok"
const PhaseInterrupted = "interrupted"
const PhaseBreaked = "breaked"
const PhaseCanceled = "canceled"

type StackFrameCursor struct {
	Type                     string
	MutableVariableValues    mentalese.Binding
	State                    map[string]int
	AllStepBindings          []mentalese.BindingSet
	ChildFrameResultBindings mentalese.BindingSet
	Phase                    string
}

func NewStackFrameCursor() *StackFrameCursor {
	return &StackFrameCursor{
		Type:                     mentalese.FrameTypePlain,
		MutableVariableValues:    mentalese.NewBinding(),
		State:                    map[string]int{},
		AllStepBindings:          []mentalese.BindingSet{},
		ChildFrameResultBindings: mentalese.NewBindingSet(),
		Phase:                    PhaseOk,
	}
}

//func (c *StackFrameCursor) UpdateMutableVariable(variable string, value mentalese.Term) {
//	for _, bindingSet := range c.AllStepBindings {
//		for _, binding := range bindingSet.GetAll() {
//			if binding.ContainsVariable(variable) {
//				binding.Set(variable, value)
//			}
//		}
//	}
//	for _, binding := range c.ChildFrameResultBindings.GetAll() {
//		if binding.ContainsVariable(variable) {
//			binding.Set(variable, value)
//		}
//	}
//}

func (c *StackFrameCursor) SetPhase(phase string) {
	c.Phase = phase
}

func (c *StackFrameCursor) GetPhase() string {
	return c.Phase
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

func (c *StackFrameCursor) SetState(name string, value int) {
	c.State[name] = value
}

func (c *StackFrameCursor) GetChildFrameResultBindings() mentalese.BindingSet {
	return c.ChildFrameResultBindings
}
