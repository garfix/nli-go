package knowledge

import "nli-go/lib/mentalese"

type KnowledgeBase interface {

	GetMatchingGroups(set mentalese.RelationSet, knowledgeBaseIndex int) RelationGroups
}
