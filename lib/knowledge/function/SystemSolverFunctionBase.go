package function

import (
	"nli-go/lib/api"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
)

// nested query structures (quant, or)
type SystemSolverFunctionBase struct {
	knowledge.KnowledgeBaseCore
	solverAsync       *central.ProblemSolverAsync
	anaphoraQueue     *central.AnaphoraQueue
	deicticCenter     *central.DeicticCenter
	discourseEntities *mentalese.Binding
	meta              *mentalese.Meta
	log               *common.SystemLog
}

func NewSystemSolverFunctionBase(anaphoraQueue *central.AnaphoraQueue, deicticCenter *central.DeicticCenter, discourseEntities *mentalese.Binding,
	meta *mentalese.Meta, log *common.SystemLog) *SystemSolverFunctionBase {
	return &SystemSolverFunctionBase{
		KnowledgeBaseCore: knowledge.KnowledgeBaseCore{ Name: "nested-structure" },
		anaphoraQueue:     anaphoraQueue,
		deicticCenter:     deicticCenter,
		discourseEntities: discourseEntities,
		meta:              meta,
		log:               log,
	}
}

func (base *SystemSolverFunctionBase) GetFunctions() map[string]api.SolverFunction {
	return map[string]api.SolverFunction{
		// grammar
		mentalese.PredicateIntent: base.intent,
		mentalese.PredicateBackReference: base.backReference,
		mentalese.PredicateDefiniteBackReference: base.definiteReference,
		mentalese.PredicateSortalBackReference: base.sortalBackReference,

		// quant
		mentalese.PredicateQuantCheck: base.quantCheck,
		mentalese.PredicateQuantForeach: base.quantForeach,
		mentalese.PredicateQuantOrderedList: base.quantOrderedList,
		
		// control
		mentalese.PredicateIfThen: base.ifThen,
		mentalese.PredicateIfThenElse: base.ifThenElse,
		mentalese.PredicateFail: base.fail,
		mentalese.PredicateLet: base.let,
		mentalese.PredicateRangeForeach: base.rangeForEach,
		mentalese.PredicateBreak: base.doBreak,
		mentalese.PredicateCancel: base.cancel,
		mentalese.PredicateWaitFor: base.waitFor,
		mentalese.PredicateCall: base.call,
		mentalese.PredicateIgnore: base.ignore,
		mentalese.PredicateAnd:	base.and,
		mentalese.PredicateXor: base.xor,
		mentalese.PredicateOr: base.or,
		mentalese.PredicateNot: base.not,
		mentalese.PredicateExec: base.exec,
		mentalese.PredicateExecResponse: base.execResponse,

		// process slots
		mentalese.PredicateSlot:          base.slot,

		// dialog context
		mentalese.PredicateContextSet:    base.contextSet,
		mentalese.PredicateContextExtend: base.contextExtend,
		mentalese.PredicateContextClear:  base.contextClear,
		mentalese.PredicateContextCall:   base.contextCall,
		mentalese.PredicateDialogReadBindings: base.dialogReadBindings,
		mentalese.PredicateDialogWriteBindings: base.dialogWriteBindings,
		mentalese.PredicateDialogAnaphoraQueueLast: base.dialogAnaphoraQueueLast,

		// list
		mentalese.PredicateListOrder: base.listOrder,
		mentalese.PredicateListAppend: base.listAppend,
		mentalese.PredicateListForeach: base.listForeach,
		mentalese.PredicateListDeduplicate: base.listDeduplicate,
		mentalese.PredicateListSort: base.listSort,
		mentalese.PredicateListIndex: base.listIndex,
		mentalese.PredicateListExpand: base.listExpand,
	}
}

