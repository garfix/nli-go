package nested

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
)

// nested query structures (quant, or)
type SystemNestedStructureBase struct {
	knowledge.KnowledgeBaseCore
	solver        *central.ProblemSolver
	dialogContext *central.DialogContext
	meta          *mentalese.Meta
	log           *common.SystemLog
}

func NewSystemNestedStructureBase(solver *central.ProblemSolver, dialogContext *central.DialogContext, meta *mentalese.Meta, log *common.SystemLog) *SystemNestedStructureBase {
	return &SystemNestedStructureBase{
		KnowledgeBaseCore: knowledge.KnowledgeBaseCore{ Name: "nested-structure" },
		solver:            solver,
		dialogContext:     dialogContext,
		meta:              meta,
		log:               log,
	}
}

func (base *SystemNestedStructureBase) HandlesPredicate(predicate string) bool {
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

func (base *SystemNestedStructureBase) sort(input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	if !knowledge.Validate(input, "va", base.log) {
		return mentalese.NewBindingSet()
	}

	return mentalese.InitBindingSet(binding)
}

func (base *SystemNestedStructureBase) intent(input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !knowledge.Validate(bound, "a*", base.log) {
		return mentalese.NewBindingSet()
	}

	return mentalese.InitBindingSet(binding)
}

func (base *SystemNestedStructureBase) SolveLet(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "**", base.log) { return mentalese.NewBindingSet() }

	variable := relation.Arguments[0].TermValue
	value := bound.Arguments[1]
	variables := base.solver.GetCurrentScope().GetVariables()
	(*variables).Set(variable, value)

	return mentalese.InitBindingSet(binding)
}

func (base *SystemNestedStructureBase) SolveNestedStructure(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
	newBindings := mentalese.NewBindingSet()

	if relation.Predicate == mentalese.PredicateIntent {

		newBindings = base.intent(relation, binding)

	} else if relation.Predicate == mentalese.PredicateBackReference {

		newBindings = base.SolveBackReference(relation, binding)

	} else if relation.Predicate == mentalese.PredicateDefiniteBackReference {

		newBindings = base.SolveDefiniteReference(relation, binding)

	} else if relation.Predicate == mentalese.PredicateQuantCheck {

		newBindings = base.SolveQuantCheck(relation, binding)

	} else if relation.Predicate == mentalese.PredicateQuantForeach {

		newBindings = base.SolveQuantForeach(relation, binding)

	} else if relation.Predicate == mentalese.PredicateAnd {

		newBindings = base.SolveAnd(relation, binding)

	} else if relation.Predicate == mentalese.PredicateXor {

		newBindings = base.SolveXor(relation, binding)

	} else if relation.Predicate == mentalese.PredicateOr {

		newBindings = base.SolveOr(relation, binding)

	} else if relation.Predicate == mentalese.PredicateNot {

		newBindings = base.SolveNot(relation, binding)

	} else if relation.Predicate == mentalese.PredicateIfThenElse {

		newBindings = base.SolveIfThenElse(relation, binding)

	} else if relation.Predicate == mentalese.PredicateCall {

		newBindings = base.Call(relation, binding)

	} else if relation.Predicate == mentalese.PredicateQuantOrderedList {

		newBindings = base.SolveQuantOrderedList(relation, binding)

	} else if relation.Predicate == mentalese.PredicateListOrder {

		newBindings = base.SolveListOrder(relation, binding)

	} else if relation.Predicate == mentalese.PredicateListForeach {

		newBindings = base.SolveListForeach(relation, binding)

	} else if relation.Predicate == mentalese.PredicateListDeduplicate {

		newBindings = base.listDeduplicate(relation, binding)

	} else if relation.Predicate == mentalese.PredicateListSort {

		newBindings = base.listSort(relation, binding)

	} else if relation.Predicate == mentalese.PredicateListIndex {

		newBindings = base.listIndex(relation, binding)

	} else if relation.Predicate == mentalese.PredicateListGet {

		newBindings = base.listGet(relation, binding)

	} else if relation.Predicate == mentalese.PredicateListLength {

		newBindings = base.listLength(relation, binding)

	} else if relation.Predicate == mentalese.PredicateListExpand {

		newBindings = base.listExpand(relation, binding)

	} else if relation.Predicate == mentalese.PredicateLet {

		newBindings = base.SolveLet(relation, binding)

	} else if relation.Predicate == mentalese.PredicateRangeForeach {

		newBindings = base.RangeForeach(relation, binding)

	}

	return newBindings
}
