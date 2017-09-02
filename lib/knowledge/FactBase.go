package knowledge

import "nli-go/lib/mentalese"

type FactBase interface {
	KnowledgeBase
	Bind(goal []mentalese.Relation) ([]mentalese.Binding, bool)
	GetMappings() []mentalese.DbMapping
	GetStatistics() mentalese.DbStats
}
