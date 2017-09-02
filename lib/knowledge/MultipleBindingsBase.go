package knowledge

import "nli-go/lib/mentalese"

type MultipleBindingsBase interface {
	KnowledgeBase
	// Returns false if none of the predicates matches
	Bind(goal mentalese.Relation, bindings []mentalese.Binding) ([]mentalese.Binding, bool)
}
