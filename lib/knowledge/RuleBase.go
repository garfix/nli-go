package knowledge

import "nli-go/lib/mentalese"

type RuleBase interface {
	KnowledgeBase
	Bind(goal mentalese.Relation, binding mentalese.Binding) ([]mentalese.RelationSet, mentalese.Bindings)
}