package knowledge

import "nli-go/lib/mentalese"

type RuleBase interface {
	KnowledgeBase
	Bind(goal mentalese.Relation) ([]mentalese.RelationSet, mentalese.Bindings)
}