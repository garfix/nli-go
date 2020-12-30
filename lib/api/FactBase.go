package api

import "nli-go/lib/mentalese"

// Knowledge base that retrieves, assets and retracts single facts
type FactBase interface {
	KnowledgeBase
	MatchRelationToDatabase(needleRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet
	Assert(relation mentalese.Relation)
	Retract(relation mentalese.Relation)
	GetReadMappings() []mentalese.Rule
	GetWriteMappings() []mentalese.Rule
	GetLocalId(sharedId string, sort string) string
	GetSharedId(localId string, sort string) string
}

type SessionBasedFactBase interface {
	ResetSession()
}