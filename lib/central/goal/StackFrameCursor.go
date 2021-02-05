package goal

import "nli-go/lib/mentalese"

type StackFrameCursor struct {
	State                    mentalese.Binding
	StepBindings             mentalese.BindingSet
	ChildFrameResultBindings mentalese.BindingSet
}

func NewStackFrameCursor() *StackFrameCursor {
	return &StackFrameCursor{
		State:                    mentalese.NewBinding(),
		StepBindings:             mentalese.NewBindingSet(),
		ChildFrameResultBindings: mentalese.NewBindingSet(),
	}
}