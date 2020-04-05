package knowledge

import (
	"nli-go/lib/common"
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
