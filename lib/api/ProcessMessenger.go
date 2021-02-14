package api

import "nli-go/lib/mentalese"

// This type is an intermediary between a process and the rest of the framework
// Its purpose is to expose as little as possible of the internal state of the process

type ProcessMessenger interface {
	GetCursor() ProcessCursor
	CreateChildStackFrame(relations mentalese.RelationSet, bindings mentalese.BindingSet)
}
