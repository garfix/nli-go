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

// Note! An internal fact base would use the same predicates as the domain language;
// This is an simulation of an external database
func (factBase InMemoryFactBase) Bind(goal []mentalese.Relation) ([]mentalese.Binding, bool) {

	factBase.log.StartDebug("Factbase Bind", goal)

	internalBindings, match := factBase.matcher.MatchSequenceToSet(goal, factBase.facts, mentalese.Binding{})

	factBase.log.EndDebug("Factbase Bind", internalBindings, match)

	return internalBindings, match
}
