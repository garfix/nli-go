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
	solver        *central.ProblemSolver
	dialogContext *central.DialogContext
	meta          *mentalese.Meta
	log           *common.SystemLog
}

func NewSystemSolverFunctionBase(solver *central.ProblemSolver, dialogContext *central.DialogContext, meta *mentalese.Meta, log *common.SystemLog) *SystemSolverFunctionBase {
	return &SystemSolverFunctionBase{
		KnowledgeBaseCore: knowledge.KnowledgeBaseCore{ Name: "nested-structure" },
		solver:            solver,
		dialogContext:     dialogContext,
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
		mentalese.PredicateLet: base.let, // todo: remove
		mentalese.PredicateRangeForeach: base.rangeForEach,
		mentalese.PredicateBreak: base.doBreak,
		mentalese.PredicateCall: base.call,
		mentalese.PredicateIgnore: base.ignore,
		mentalese.PredicateAnd:	base.and,
		mentalese.PredicateXor: base.xor,
		mentalese.PredicateOr: base.or,
		mentalese.PredicateNot: base.not,
		mentalese.PredicateExec: base.exec,
		mentalese.PredicateExecResponse: base.execResponse,
		
		// list
		mentalese.PredicateListOrder: base.listOrder,
		mentalese.PredicateListAppend: base.listAppend,
		mentalese.PredicateListForeach: base.listForeach,
		mentalese.PredicateListDeduplicate: base.listDeduplicate,
		mentalese.PredicateListSort: base.listSort,
		mentalese.PredicateListIndex: base.listIndex,
		mentalese.PredicateListGet: base.listGet,
		mentalese.PredicateListLength: base.listLength,
		mentalese.PredicateListExpand: base.listExpand,
	}
}

