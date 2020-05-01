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
	solver *central.ProblemSolver
	dialogContext *central.DialogContext
	predicates mentalese.Predicates
	log     *common.SystemLog
}

func NewSystemNestedStructureBase(solver *central.ProblemSolver, dialogContext *central.DialogContext, predicates mentalese.Predicates, log *common.SystemLog) *SystemNestedStructureBase {
	return &SystemNestedStructureBase{
		KnowledgeBaseCore: knowledge.KnowledgeBaseCore{ Name: "nested-structure" },
		solver: solver,
		dialogContext: dialogContext,
		predicates: predicates,
		log:               log,
	}
}

func (base *SystemNestedStructureBase) HandlesPredicate(predicate string) bool {
	predicates := []string{
		mentalese.PredicateDo,
		mentalese.PredicateFind,
		mentalese.PredicateCall,
		mentalese.PredicateAnd,
		mentalese.PredicateOr,
		mentalese.PredicateNot,
		mentalese.PredicateBackReference,
		mentalese.PredicateDefiniteBackReference}

	for _, p := range predicates {
		if p == predicate {
			return true
		}
	}
	return false
}

func (base *SystemNestedStructureBase) SolveNestedStructure(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {
	var newBindings mentalese.Bindings

	if relation.Predicate == mentalese.PredicateBackReference {

		newBindings = base.SolveBackReference(relation, binding)

	} else if relation.Predicate == mentalese.PredicateDefiniteBackReference {

		newBindings = base.SolveDefiniteReference(relation, binding)

	} else if relation.Predicate == mentalese.PredicateFind {

		newBindings = base.SolveFind(relation, binding)

	} else if relation.Predicate == mentalese.PredicateDo {

		newBindings = base.SolveDo(relation, binding)

	} else if relation.Predicate == mentalese.PredicateAnd {

		newBindings = base.SolveAnd(relation, binding)

	} else if relation.Predicate == mentalese.PredicateOr {

		newBindings = base.SolveOr(relation, binding)

	} else if relation.Predicate == mentalese.PredicateNot {

		newBindings = base.SolveNot(relation, binding)

	} else if relation.Predicate == mentalese.PredicateCall {

		newBindings = base.Call(relation, binding)

	}

	return newBindings
}
