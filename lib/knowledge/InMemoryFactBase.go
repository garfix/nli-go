package knowledge

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type InMemoryFactBase struct {
	facts   mentalese.RelationSet
	ds2db   []mentalese.DbMapping
	stats	mentalese.DbStats
	matcher *mentalese.RelationMatcher
	log     *common.SystemLog
}

func NewInMemoryFactBase(facts mentalese.RelationSet, ds2db []mentalese.DbMapping, stats mentalese.DbStats, log *common.SystemLog) mentalese.FactBase {
	return InMemoryFactBase{facts: facts, ds2db: ds2db, stats: stats, matcher: mentalese.NewRelationMatcher(log), log: log}
}

func (factBase InMemoryFactBase) GetMappings() []mentalese.DbMapping {
	return factBase.ds2db
}

func (factBase InMemoryFactBase) GetStatistics() mentalese.DbStats {
	return factBase.stats
}

// Note! An internal fact base would use the same predicates as the domain language;
// This is an simulation of an external database
func (factBase InMemoryFactBase) Bind(goal []mentalese.Relation) ([]mentalese.Binding, bool) {

	factBase.log.StartDebug("Factbase Bind", goal)

	internalBindings, _, match := factBase.matcher.MatchSequenceToSet(goal, factBase.facts, mentalese.Binding{})

	factBase.log.EndDebug("Factbase Bind", internalBindings, match)

	return internalBindings, match
}
