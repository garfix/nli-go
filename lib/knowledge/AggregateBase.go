package knowledge

import "nli-go/lib/mentalese"

type AggregateBase interface {
	KnowledgeBase
	// Returns false if none of the predicates matches
	Bind(goal mentalese.Relation, bindings mentalese.Bindings) (mentalese.Bindings, bool)
}
