package api

import "nli-go/lib/mentalese"

// nested query structures (quant, or)
type SolverFunctionBase interface {
	KnowledgeBase
	Execute(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet
}
