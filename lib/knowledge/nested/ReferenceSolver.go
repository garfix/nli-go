package nested

import "nli-go/lib/mentalese"

func (base *SystemNestedStructureBase) SolveBackReference(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	variable := relation.Arguments[0].TermValue
	set := relation.Arguments[1].TermValueRelationSet

	newBindings := mentalese.Bindings{}

	for _, group := range *base.dialogContext.AnaphoraQueue {

		ref := group[0]

		refBinding := binding.Merge(mentalese.Binding{ variable: mentalese.NewTermId(ref.Id, ref.EntityType)})

		// empty set ("it")
		if len(set) == 0 {
			newBindings = mentalese.Bindings{ refBinding }
			break
		}

		if !base.quickAcceptabilityCheck(variable, ref.EntityType, set) {
			continue
		}

		testRangeBindings := base.solver.SolveRelationSet(set, mentalese.Bindings{refBinding})
		if len(testRangeBindings) == 1 {
			newBindings = testRangeBindings
			break
		}
	}

	return newBindings
}

func (base *SystemNestedStructureBase) SolveDefiniteReference(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	variable := relation.Arguments[0].TermValue
	set := relation.Arguments[1].TermValueRelationSet

	newBindings := base.SolveBackReference(relation, binding)

	if len(newBindings) == 0 {
		newBindings = base.solver.SolveRelationSet(set, mentalese.Bindings{binding})

		if len(newBindings) > 1 {
			rangeIndex, found := base.rangeIndexClarification(newBindings, variable)
			if found {
				newBindings = newBindings[rangeIndex:rangeIndex + 1]
			} else {
				return mentalese.Bindings{}
			}
		}
	}

	return newBindings
}