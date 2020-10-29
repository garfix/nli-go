package api

import "nli-go/lib/mentalese"

// Knowledge bases that take all current bindings as input at once
type MultiBindingBase interface {
	KnowledgeBase
	// Returns false if none of the predicates matches
	Execute(goal mentalese.Relation, bindings mentalese.BindingSet) (mentalese.BindingSet, bool)
}
