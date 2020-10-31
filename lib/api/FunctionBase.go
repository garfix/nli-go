package api

// Knowledge bases that processes functions with a single binding
// These functions can be used everywhere relations are used
type FunctionBase interface {
	KnowledgeBase
	GetFunctions() map[string]SimpleFunction
}
