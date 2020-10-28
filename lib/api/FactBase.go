package api

import "nli-go/lib/mentalese"

type FactBase interface {
	KnowledgeBase
	MatchRelationToDatabase(needleRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet
	Assert(relation mentalese.Relation)
	Retract(relation mentalese.Relation)
	GetReadMappings() []mentalese.Rule
	GetWriteMappings() []mentalese.Rule
	GetLocalId(sharedId string, entityType string) string
	GetSharedId(localId string, entityType string) string
}
