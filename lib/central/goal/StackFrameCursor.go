package goal

import "nli-go/lib/mentalese"

type StackFrameCursor struct {
	State                    mentalese.Binding
	Bindings                 mentalese.BindingSet
	ChildFrameResultBindings mentalese.BindingSet
}

func NewStackFrameCursor() *StackFrameCursor {
	return &StackFrameCursor{
		State:    mentalese.NewBinding(),
		Bindings: mentalese.NewBindingSet(),
		ChildFrameResultBindings: mentalese.NewBindingSet(),
	}
}