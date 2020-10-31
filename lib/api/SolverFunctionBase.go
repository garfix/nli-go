package api

// A function base whose predicates cannot be used everywhere, only in the solving process
type SolverFunctionBase interface {
	KnowledgeBase
	GetFunctions() map[string]SolverFunction
}
