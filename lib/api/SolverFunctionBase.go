package api

import "nli-go/lib/mentalese"

// A function base whose predicates cannot be used everywhere, only in the solving process
type SolverFunctionBase interface {
	KnowledgeBase
	Execute(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet
}
