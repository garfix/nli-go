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
		mentalese.PredicateLet,
	}

	for _, p := range predicates {
		if p == predicate {
			return true
		}
	}
	return false
}

func (base *SystemNestedStructureBase) sort(input mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	if !knowledge.Validate(input, "va", base.log) {
		return nil
	}

	return mentalese.Bindings{ binding }
}

func (base *SystemNestedStructureBase) intent(input mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	bound := input.BindSingle(binding)

	if !knowledge.Validate(bound, "a*", base.log) {
		return nil
	}

	return mentalese.Bindings{ binding }
}

func (base *SystemNestedStructureBase) SolveLet(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	if !knowledge.Validate(relation, "v*", base.log) { return nil }

	variable := relation.Arguments[0].TermValue
	value := relation.Arguments[1]
	variables := base.solver.GetCurrentScope().GetVariables()
	(*variables).Set(variable, value)

	return mentalese.Bindings{ binding }
}

func (base *SystemNestedStructureBase) SolveNestedStructure(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {
	var newBindings mentalese.Bindings

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

	} else if relation.Predicate == mentalese.PredicateLet {

		newBindings = base.SolveLet(relation, binding)
	}

	return newBindings
}
