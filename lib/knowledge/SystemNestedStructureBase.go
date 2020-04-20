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

func (base *SystemNestedStructureBase) HandlesPredicate(predicate string) bool {
	predicates := []string{mentalese.PredicateDo, mentalese.PredicateFind, mentalese.PredicateCall, mentalese.PredicateSequence, mentalese.PredicateNot}

	for _, p := range predicates {
		if p == predicate {
			return true
		}
	}
	return false
}