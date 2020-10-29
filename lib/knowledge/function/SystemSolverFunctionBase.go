package function

import (
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

func (base *SystemSolverFunctionBase) HandlesPredicate(predicate string) bool {
	predicates := []string {
		mentalese.PredicateIntent,
		mentalese.PredicateQuantForeach,
		mentalese.PredicateQuantCheck,
		mentalese.PredicateCall,
		mentalese.PredicateAnd,
		mentalese.PredicateOr,
		mentalese.PredicateXor,
		mentalese.PredicateNot,
		mentalese.PredicateBackReference,
		mentalese.PredicateIfThenElse,
		mentalese.PredicateDefiniteBackReference,
		mentalese.PredicateQuantOrderedList,
		mentalese.PredicateListOrder,
		mentalese.PredicateListForeach,
		mentalese.PredicateListDeduplicate,
		mentalese.PredicateListSort,
		mentalese.PredicateListIndex,
		mentalese.PredicateListGet,
		mentalese.PredicateListLength,
		mentalese.PredicateListExpand,
		mentalese.PredicateLet,
		mentalese.PredicateRangeForeach,
	}

	for _, p := range predicates {
		if p == predicate {
			return true
		}
	}
	return false
}

func (base *SystemSolverFunctionBase) Execute(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
	newBindings := mentalese.NewBindingSet()

	switch relation.Predicate {

	// grammar
	case mentalese.PredicateIntent:
		newBindings = base.intent(relation, binding)
	case mentalese.PredicateBackReference:
		newBindings = base.backReference(relation, binding)
	case mentalese.PredicateDefiniteBackReference:
		newBindings = base.definiteReference(relation, binding)

	// quant
	case mentalese.PredicateQuantCheck:
		newBindings = base.quantCheck(relation, binding)
	case mentalese.PredicateQuantForeach:
		newBindings = base.quantForeach(relation, binding)
	case mentalese.PredicateQuantOrderedList:
		newBindings = base.quantOrderedList(relation, binding)

	// control
	case mentalese.PredicateIfThenElse:
		newBindings = base.ifThenElse(relation, binding)
	case mentalese.PredicateLet:
		newBindings = base.let(relation, binding)
	case mentalese.PredicateRangeForeach:
		newBindings = base.rangeForEach(relation, binding)
	case mentalese.PredicateCall:
		newBindings = base.call(relation, binding)
	case mentalese.PredicateAnd:
		newBindings = base.and(relation, binding)
	case mentalese.PredicateXor:
		newBindings = base.xor(relation, binding)
	case mentalese.PredicateOr:
		newBindings = base.or(relation, binding)
	case mentalese.PredicateNot:
		newBindings = base.not(relation, binding)

	//list
	case mentalese.PredicateListOrder:
		newBindings = base.listOrder(relation, binding)
	case mentalese.PredicateListForeach:
		newBindings = base.listForeach(relation, binding)
	case mentalese.PredicateListDeduplicate:
		newBindings = base.listDeduplicate(relation, binding)
	case mentalese.PredicateListSort:
		newBindings = base.listSort(relation, binding)
	case mentalese.PredicateListIndex:
		newBindings = base.listIndex(relation, binding)
	case mentalese.PredicateListGet:
		newBindings = base.listGet(relation, binding)
	case mentalese.PredicateListLength:
		newBindings = base.listLength(relation, binding)
	case mentalese.PredicateListExpand:
		newBindings = base.listExpand(relation, binding)
	}

	return newBindings
}

