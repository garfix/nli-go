package knowledge

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"sync"
)

type Find struct {
	index    int
	binding  mentalese.Binding
	relation mentalese.Relation
}

type InMemoryFactBase struct {
	KnowledgeBaseCore
	originalFacts map[string]mentalese.RelationSet
	facts         map[string]mentalese.RelationSet
	readMap       []mentalese.Rule
	writeMap      []mentalese.Rule
	entities      mentalese.SortProperties
	sharedIds     SharedIds
	matcher       *central.RelationMatcher
	log           *common.SystemLog
	mutex         sync.Mutex
}

func NewInMemoryFactBase(name string, facts mentalese.RelationSet, matcher *central.RelationMatcher, readMap []mentalese.Rule, writeMap []mentalese.Rule, log *common.SystemLog) *InMemoryFactBase {

	indexedFacts := map[string]mentalese.RelationSet{}
	for _, fact := range facts {
		_, found := indexedFacts[fact.Predicate]
		if !found {
			indexedFacts[fact.Predicate] = mentalese.RelationSet{}
		}
		indexedFacts[fact.Predicate] = append(indexedFacts[fact.Predicate], fact)
	}

	factBase := InMemoryFactBase{
		KnowledgeBaseCore: KnowledgeBaseCore{Name: name},
		originalFacts:     indexedFacts,
		facts:             indexedFacts,
		readMap:           readMap,
		writeMap:          writeMap,
		sharedIds:         SharedIds{},
		matcher:           matcher,
		log:               log,
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
	if !found {
		return inId
	}

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
	if !found {
		return inId
	}

	for localId, sharedId := range factBase.sharedIds[sort] {
		if inId == localId {
			outId = sharedId
			break
		}
	}

	return outId
}

func (factBase *InMemoryFactBase) EnsurePresent(predicate string) {
	_, found := factBase.facts[predicate]
	if !found {
		factBase.facts[predicate] = mentalese.RelationSet{}
	}
}

func (factBase *InMemoryFactBase) FindFacts(needleRelation mentalese.Relation, binding mentalese.Binding) []Find {
	finds := []Find{}
	_, found := factBase.facts[needleRelation.Predicate]
	if !found {
		return finds
	}
	for i, fact := range factBase.facts[needleRelation.Predicate] {
		b, found := factBase.matcher.MatchTwoRelations(needleRelation, fact, binding)
		if found {
			finds = append(finds, Find{i, b, fact})
		}
	}
	return finds
}

func (factBase *InMemoryFactBase) MatchRelationToDatabase(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	factBase.mutex.Lock()

	finds := factBase.FindFacts(relation, binding)
	bindings := mentalese.NewBindingSet()
	for _, find := range finds {
		bindings.Add(find.binding)
	}

	factBase.mutex.Unlock()

	return bindings
}

func (factBase *InMemoryFactBase) Assert(relation mentalese.Relation) {

	factBase.mutex.Lock()

	predicate := relation.Predicate
	finds := factBase.FindFacts(relation, mentalese.NewBinding())
	if len(finds) == 0 {
		factBase.EnsurePresent(relation.Predicate)
		factBase.facts[predicate] = append(factBase.facts[predicate], relation)
	}

	factBase.mutex.Unlock()
}

// Removes all facts that match relation
func (factBase *InMemoryFactBase) Retract(relation mentalese.Relation) {

	factBase.mutex.Lock()

	predicate := relation.Predicate
	finds := factBase.FindFacts(relation, mentalese.NewBinding())
	if len(finds) > 0 {
		newFacts := []mentalese.Relation{}

		for i, fact := range factBase.facts[predicate] {
			found := false
			for _, find := range finds {
				if find.index == i {
					found = true
					break
				}
			}
			if !found {
				newFacts = append(newFacts, fact)
			}
		}
		factBase.facts[predicate] = newFacts
	}

	factBase.mutex.Unlock()
}

func (factBase *InMemoryFactBase) ResetSession() {

	factBase.mutex.Lock()

	c := map[string]mentalese.RelationSet{}
	for predicate, facts := range factBase.facts {
		c[predicate] = facts.Copy()
	}

	factBase.facts = c

	factBase.mutex.Unlock()
}
