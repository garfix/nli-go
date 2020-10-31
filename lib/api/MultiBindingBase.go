package api

// Knowledge bases that take all current bindings as input at once
type MultiBindingBase interface {
	KnowledgeBase
	GetFunctions() map[string]MultiBindingFunction
}
