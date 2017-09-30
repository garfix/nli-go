package knowledge

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type InMemoryFactBase struct {
	facts   mentalese.RelationSet
	ds2db   []mentalese.RelationTransformation
	stats	mentalese.DbStats
	matcher *mentalese.RelationMatcher
	log     *common.SystemLog
}

func NewInMemoryFactBase(facts mentalese.RelationSet, matcher *mentalese.RelationMatcher, ds2db []mentalese.RelationTransformation, stats mentalese.DbStats, log *common.SystemLog) FactBase {
	return InMemoryFactBase{facts: facts, ds2db: ds2db, stats: stats, matcher: matcher, log: log}
}

func (factBase InMemoryFactBase) GetMappings() []mentalese.RelationTransformation {
	return factBase.ds2db
}

func (factBase InMemoryFactBase) GetMatchingGroups(set mentalese.RelationSet, knowledgeBaseIndex int) []RelationGroup {
	return getFactBaseMatchingGroups(factBase.matcher, set, factBase, knowledgeBaseIndex)
}

func (factBase InMemoryFactBase) GetStatistics() mentalese.DbStats {
	return factBase.stats
}

func (factBase InMemoryFactBase) MatchRelationToDatabase(needleRelation mentalese.Relation) []mentalese.Binding {

	bindings, _ := factBase.matcher.MatchRelationToSet(needleRelation, factBase.facts, mentalese.Binding{})
	return bindings
}