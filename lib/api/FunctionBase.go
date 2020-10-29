package api

import "nli-go/lib/mentalese"

// Knowledge bases that processes functions with a single binding
// These functions can be used everywhere relations are used
type FunctionBase interface {
	KnowledgeBase
	Execute(input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool, bool)
}
