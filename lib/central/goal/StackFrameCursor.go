package goal

import "nli-go/lib/mentalese"

type StackFrameCursor struct {
	State                    mentalese.Binding
	OutBindings              mentalese.BindingSet
	ChildFrameResultBindings mentalese.BindingSet
}

func NewStackFrameCursor() *StackFrameCursor {
	return &StackFrameCursor{
		State:                    mentalese.NewBinding(),
		OutBindings:              mentalese.NewBindingSet(),
		ChildFrameResultBindings: mentalese.NewBindingSet(),
	}
}