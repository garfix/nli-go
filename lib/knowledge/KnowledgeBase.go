package knowledge

import "nli-go/lib/mentalese"

type KnowledgeBase interface {
	Knows(relation mentalese.Relation) bool
}
