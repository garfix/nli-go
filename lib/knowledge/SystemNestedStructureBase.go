package knowledge

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

// nested query structures (quant, or)
type SystemNestedStructureBase struct {
	KnowledgeBaseCore
	log     *common.SystemLog
}

func NewSystemNestedStructureBase(log *common.SystemLog) *SystemNestedStructureBase {
	return &SystemNestedStructureBase{
		KnowledgeBaseCore: KnowledgeBaseCore{ Name: "nested-structure" },
		log: log,
	}
}

func (base *SystemNestedStructureBase) GetMatchingGroups(set mentalese.RelationSet, knowledgeBaseName string) []RelationGroup {

	matchingGroups := []RelationGroup{}
	predicates := []string{mentalese.Predicate_Quant}

	for _, setRelation := range set {
		for _, predicate:= range predicates {
			if predicate == setRelation.Predicate {
// TODO calculate real cost
				matchingGroups = append(matchingGroups, RelationGroup{mentalese.RelationSet{setRelation}, knowledgeBaseName, worst_cost})
			}
		}
	}

	return matchingGroups
}
