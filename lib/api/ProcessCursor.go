package api

import "nli-go/lib/mentalese"

// This struct is a working environment of a single step in a stack frame
// It is needed by relations that have child stack frames:
// when a stack frame has finished, the parent relation is re-entered and continued

type ProcessCursor interface {
	GetPhase() string
	GetType() string
	SetType(string)
	GetState(string, int) int
	SetState(string, int)
	AddStepBindings(bindings mentalese.BindingSet)
	GetAllStepBindings() []mentalese.BindingSet
	GetChildFrameResultBindings() mentalese.BindingSet
}
