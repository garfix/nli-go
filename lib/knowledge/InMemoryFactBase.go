package knowledge

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type InMemoryFactBase struct {
	KnowledgeBaseCore
	facts   	mentalese.RelationSet
	ds2db   	[]mentalese.Rule
	ds2dbWrite 	[]mentalese.Rule
	entities 	mentalese.Entities
	sharedIds 	SharedIds
	matcher 	*mentalese.RelationMatcher
	log     	*common.SystemLog
}

func NewInMemoryFactBase(name string, facts mentalese.RelationSet, matcher *mentalese.RelationMatcher, ds2db []mentalese.Rule, ds2dbWrite []mentalese.Rule, entities mentalese.Entities, log *common.SystemLog) *InMemoryFactBase {
	return &InMemoryFactBase{
		KnowledgeBaseCore: KnowledgeBaseCore{ Name: name },
		facts: facts,
		ds2db: ds2db,
		ds2dbWrite: ds2dbWrite,
		entities: entities,
		sharedIds: SharedIds{},
		matcher: matcher,
		log: log,
	}
}

func (factBase *InMemoryFactBase) HandlesPredicate(predicate string) bool {
	for _, rule := range factBase.ds2db {
		if rule.Goal.Predicate == predicate {
			return true
		}
	}
	if len(factBase.ds2dbWrite) > 0 && (predicate == mentalese.PredicateAssert || predicate == mentalese.PredicateRetract) {
		return true
	}
	return false
}

func (factBase *InMemoryFactBase) GetMappings() []mentalese.Rule {
	return factBase.ds2db
}

func (factBase *InMemoryFactBase) GetWriteMappings() []mentalese.Rule {
	return factBase.ds2dbWrite
}

func (factBase *InMemoryFactBase) GetEntities() mentalese.Entities {
	return factBase.entities
}

func (factBase *InMemoryFactBase) SetSharedIds(sharedIds SharedIds) {
	factBase.sharedIds = sharedIds
}

func (factBase *InMemoryFactBase) GetLocalId(inId string, entityType string) string {
	outId := ""

	_, found := factBase.sharedIds[entityType]
	if !found { return inId }

	for localId, sharedId := range factBase.sharedIds[entityType] {
		if inId == sharedId {
			outId = localId
			break
		}
	}

	return outId
}

func (factBase *InMemoryFactBase) GetSharedId(inId string, entityType string) string {
	outId := ""

	_, found := factBase.sharedIds[entityType]
	if !found { return inId }

	for localId, sharedId := range factBase.sharedIds[entityType] {
		if inId == localId {
			outId = sharedId
			break
		}
	}

	return outId
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

func (factBase *InMemoryFactBase) MatchRelationToDatabase(needleRelation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	bindings, _ := factBase.matcher.MatchRelationToSet(needleRelation, factBase.facts, binding)
	return bindings
}