package knowledge

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type InMemoryFactBase struct {
	KnowledgeBaseCore
	facts   mentalese.RelationSet
	ds2db   []mentalese.RelationTransformation
	ds2dbWrite []mentalese.RelationTransformation
	stats	mentalese.DbStats
	entities 		  mentalese.Entities
	matcher *mentalese.RelationMatcher
	log     *common.SystemLog
}

func NewInMemoryFactBase(name string, facts mentalese.RelationSet, matcher *mentalese.RelationMatcher, ds2db []mentalese.RelationTransformation, ds2dbWrite []mentalese.RelationTransformation, stats mentalese.DbStats, entities mentalese.Entities, log *common.SystemLog) *InMemoryFactBase {
	return &InMemoryFactBase{
		KnowledgeBaseCore: KnowledgeBaseCore{ Name: name },
		facts: facts,
		ds2db: ds2db,
		ds2dbWrite: ds2dbWrite,
		stats: stats,
		entities: entities,
		matcher: matcher,
		log: log,
	}
}

func (factBase *InMemoryFactBase) GetMappings() []mentalese.RelationTransformation {
	return factBase.ds2db
}

func (factBase *InMemoryFactBase) GetWriteMappings() []mentalese.RelationTransformation {
	return factBase.ds2dbWrite
}

func (factBase *InMemoryFactBase) GetMatchingGroups(set mentalese.RelationSet, keyCabinet *mentalese.KeyCabinet) []RelationGroup {
	return getFactBaseMatchingGroups(factBase.matcher, set, factBase, keyCabinet)
}

func (factBase *InMemoryFactBase) GetStatistics() mentalese.DbStats {
	return factBase.stats
}

func (factBase *InMemoryFactBase) GetEntities() mentalese.Entities {
	return factBase.entities
}

func (factBase *InMemoryFactBase) SetRelations(relations mentalese.RelationSet) {
	factBase.facts = relations
}

func (factBase *InMemoryFactBase) GetRelations() mentalese.RelationSet {
	return factBase.facts
}

func (factBase *InMemoryFactBase) AddRelation(relation mentalese.Relation) {
	factBase.facts = append(factBase.facts, relation)
}

func (factBase *InMemoryFactBase) Assert(relation mentalese.Relation) {

	for _, fact := range factBase.facts {
		_, found := factBase.matcher.MatchTwoRelations(relation, fact, mentalese.Binding{})
		if found {
			return
		}
	}

	factBase.facts = append(factBase.facts, relation)
}

func (factBase *InMemoryFactBase) Retract(relation mentalese.Relation) {
	factBase.RemoveRelation(relation)
}

// Removes all facts that match relation
func (factBase *InMemoryFactBase) RemoveRelation(relation mentalese.Relation) {
	newFacts := []mentalese.Relation{}

	for _, fact := range factBase.facts {
		_, found := factBase.matcher.MatchTwoRelations(relation, fact, mentalese.Binding{})
		if !found {
			newFacts = append(newFacts, fact)
		}
	}

	factBase.facts = newFacts
}

func (factBase *InMemoryFactBase) MatchRelationToDatabase(needleRelation mentalese.Relation) []mentalese.Binding {

	bindings, _ := factBase.matcher.MatchRelationToSet(needleRelation, factBase.facts, mentalese.Binding{})
	return bindings
}