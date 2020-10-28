package api

import "nli-go/lib/mentalese"

type AggregateFunctionBase interface {
	KnowledgeBase
	// Returns false if none of the predicates matches
	Execute(goal mentalese.Relation, bindings mentalese.BindingSet) (mentalese.BindingSet, bool)
}
