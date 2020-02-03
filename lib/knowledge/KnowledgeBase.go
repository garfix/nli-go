package knowledge

import (
	"nli-go/lib/mentalese"
)

type KnowledgeBase interface {

	GetName() string
	GetMatchingGroups(set mentalese.RelationSet) []RelationGroup
}
