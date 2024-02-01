package api

import "nli-go/lib/mentalese"

// This type is an intermediary between a process and the rest of the framework
// Its purpose is to expose as little as possible of the internal state of the process

type ProcessMessenger interface {
	GetProcess() Process
	GetCursor() ProcessCursor
	ExecuteChildStackFrame(relations mentalese.RelationSet, bindings mentalese.BindingSet) mentalese.BindingSet
	ExecuteChildStackFrameMutable(relations mentalese.RelationSet, bindings mentalese.Binding) mentalese.BindingSet
	StartProcess(resource string, relations mentalese.RelationSet, binding mentalese.Binding) bool
	AddProcessInstruction(name string, value string)
	GetProcessSlot(slot string) (mentalese.Term, bool)
	SetProcessSlot(slot string, value mentalese.Term)
	SetMutableVariable(variable string, value mentalese.Term)
	GetMutableVariable(variable string) (mentalese.Term, bool)
	SetMutableVariables(binding mentalese.Binding)
	GetMutableVariables() mentalese.Binding
}
