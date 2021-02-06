package goal

import "nli-go/lib/mentalese"

// This struct is a working environment of a single step in a stack frame
// It is needed by relations that have child stack frames:
// when a stack frame has finished, the parent relation is re-entered and continued

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