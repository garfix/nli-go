package knowledge

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
)

// nested query structures (quant, or)
type NestedStructureBase struct {
	KnowledgeBase
	log     *common.SystemLog
}

func NewNestedStructureBase(log *common.SystemLog) NestedStructureBase {
	return NestedStructureBase{log: log}
}

func (base  NestedStructureBase) GetMatchingGroups(set mentalese.RelationSet, knowledgeBaseIndex int) []RelationGroup {

	matchingGroups := []RelationGroup{}
	predicates := []string{mentalese.Predicate_Quant}

	for _, setRelation := range set {
		for _, predicate:= range predicates {
			if predicate == setRelation.Predicate {
// TODO calculate real cost
				matchingGroups = append(matchingGroups, RelationGroup{mentalese.RelationSet{setRelation}, knowledgeBaseIndex, worst_cost})
			}
		}
	}

	return matchingGroups
}
