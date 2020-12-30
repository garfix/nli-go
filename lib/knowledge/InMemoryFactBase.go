package knowledge

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type InMemoryFactBase struct {
	KnowledgeBaseCore
	originalFacts     mentalese.RelationSet
	facts     mentalese.RelationSet
	readMap   []mentalese.Rule
	writeMap  []mentalese.Rule
	entities  mentalese.Entities
	sharedIds SharedIds
	matcher   *central.RelationMatcher
	storage *common.FileStorage
	log       *common.SystemLog
}

func NewInMemoryFactBase(name string, facts mentalese.RelationSet, matcher *central.RelationMatcher, readMap []mentalese.Rule, writeMap []mentalese.Rule, storage *common.FileStorage, log *common.SystemLog) *InMemoryFactBase {
	factBase := InMemoryFactBase{
		KnowledgeBaseCore: KnowledgeBaseCore{ Name: name },
		originalFacts: 	   facts,
		facts:             facts.Copy(),
		readMap:           readMap,
		writeMap:          writeMap,
		sharedIds:         SharedIds{},
		matcher:           matcher,
		storage:           storage,
		log:               log,
	}

	if storage != nil {
		storage.Read(&factBase.facts)
	}

	return &factBase
}

func (factBase *InMemoryFactBase) GetReadMappings() []mentalese.Rule {
	return factBase.readMap
}

func (factBase *InMemoryFactBase) GetWriteMappings() []mentalese.Rule {
	return factBase.writeMap
}

func (factBase *InMemoryFactBase) SetSharedIds(sharedIds SharedIds) {
	factBase.sharedIds = sharedIds
}

func (factBase *InMemoryFactBase) GetLocalId(inId string, sort string) string {
	outId := ""

	_, found := factBase.sharedIds[sort]
	if !found { return inId }

	for localId, sharedId := range factBase.sharedIds[sort] {
		if inId == sharedId {
			outId = localId
			break
		}
	}

	return outId
}

func (factBase *InMemoryFactBase) GetSharedId(inId string, sort string) string {
	outId := ""

	_, found := factBase.sharedIds[sort]
	if !found { return inId }

	for localId, sharedId := range factBase.sharedIds[sort] {
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
		_, found := factBase.matcher.MatchTwoRelations(relation, fact, mentalese.NewBinding())
		if found {
			return
		}
	}

	factBase.facts = append(factBase.facts, relation)

	factBase.persist()
}

func (factBase *InMemoryFactBase) Retract(relation mentalese.Relation) {
	factBase.RemoveRelation(relation)

	factBase.persist()
}

func (factBase *InMemoryFactBase) ResetSession() {
	factBase.facts = factBase.originalFacts.Copy()

	factBase.persist()
}

func (factBase *InMemoryFactBase) persist() {
	if factBase.storage != nil {
		factBase.storage.Write(factBase.facts)
	}
}

// Removes all facts that match relation
func (factBase *InMemoryFactBase) RemoveRelation(relation mentalese.Relation) {
	newFacts := []mentalese.Relation{}

	for _, fact := range factBase.facts {
		_, found := factBase.matcher.MatchTwoRelations(relation, fact, mentalese.NewBinding())
		if !found {
			newFacts = append(newFacts, fact)
		}
	}

	factBase.facts = newFacts
}

func (factBase *InMemoryFactBase) MatchRelationToDatabase(needleRelation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bindings, _ := factBase.matcher.MatchRelationToSet(needleRelation, factBase.facts, binding)
	return bindings
}
