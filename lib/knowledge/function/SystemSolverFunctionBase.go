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
	solverAsync     *central.ProblemSolver
	dialogContext   *central.DialogContext
	meta            *mentalese.Meta
	log             *common.SystemLog
	clientConnector api.ClientConnector
	choices         map[string]string
}

func NewSystemSolverFunctionBase(dialogContext *central.DialogContext, meta *mentalese.Meta, log *common.SystemLog, clientConnector api.ClientConnector) *SystemSolverFunctionBase {
	return &SystemSolverFunctionBase{
		KnowledgeBaseCore: knowledge.KnowledgeBaseCore{Name: "nested-structure"},
		dialogContext:     dialogContext,
		meta:              meta,
		log:               log,
		clientConnector:   clientConnector,
		choices:           map[string]string{},
	}
}

func (base *SystemSolverFunctionBase) GetFunctions() map[string]api.SolverFunction {
	return map[string]api.SolverFunction{

		// quant
		mentalese.PredicateCheck:            base.quantCheck,
		mentalese.PredicateDo:               base.quantForeach,
		mentalese.PredicateQuantOrderedList: base.quantOrderedList,
		mentalese.PredicateQuantMatch:       base.quantMatch,

		// control
		mentalese.PredicateIfThen:         base.ifThen,
		mentalese.PredicateIfThenElse:     base.ifThenElse,
		mentalese.PredicateIfThenBool:     base.ifThenBool,
		mentalese.PredicateIfThenElseBool: base.ifThenElseBool,
		mentalese.PredicateFail:           base.fail,
		mentalese.PredicateReturn:         base.returnStatement,
		mentalese.PredicateAssign:         base.assign,
		mentalese.PredicateRangeForeach:   base.rangeForEach,
		mentalese.PredicateForRelations:   base.forRelations,
		mentalese.PredicateForIndexValue:  base.forIndexValue,
		mentalese.PredicateListIndex2:     base.listIndex2,
		mentalese.PredicateForRange:       base.forRange,
		mentalese.PredicateBreak:          base.doBreak,
		mentalese.PredicateCancel:         base.cancel,
		mentalese.PredicateWaitFor:        base.waitFor,
		mentalese.PredicateCall:           base.call,
		mentalese.PredicateApply:          base.apply,
		mentalese.PredicateIgnore:         base.ignore,
		mentalese.PredicateAnd:            base.and,
		mentalese.PredicateXor:            base.xor,
		mentalese.PredicateOr:             base.or,
		mentalese.PredicateNot:            base.not,
		mentalese.PredicateExec:           base.exec,
		mentalese.PredicateExecResponse:   base.execResponse,

		// dialog context
		mentalese.PredicateContextSet:    base.contextSet,
		mentalese.PredicateContextGet:    base.contextGet,
		mentalese.PredicateContextExtend: base.contextExtend,
		mentalese.PredicateContextClear:  base.contextClear,
		mentalese.PredicateContextCall:   base.contextCall,
		mentalese.PredicateCreateGoal:    base.createGoal,

		// list
		mentalese.PredicateListOrder:       base.listOrder,
		mentalese.PredicateListAppend:      base.listAppend,
		mentalese.PredicateListSet:         base.listSet,
		mentalese.PredicateListForeach:     base.listForeach,
		mentalese.PredicateListDeduplicate: base.listDeduplicate,
		mentalese.PredicateListSort:        base.listSort,
		mentalese.PredicateListRemove:      base.listRemove,
		mentalese.PredicateListIndex:       base.listIndex,
		mentalese.PredicateListExpand:      base.listExpand,
	}
}
