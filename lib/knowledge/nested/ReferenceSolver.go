package nested

import "nli-go/lib/mentalese"

func (base *SystemNestedStructureBase) SolveBackReference(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	variable := relation.Arguments[0].TermValue
	set := relation.Arguments[1].TermValueRelationSet

	newBindings := mentalese.NewBindingSet()

	for _, group := range *base.dialogContext.AnaphoraQueue {

		ref := group[0]

		b := mentalese.NewBinding()
		b.Set(variable, mentalese.NewTermId(ref.Id, ref.EntityType))

		refBinding := binding.Merge(b)

		// empty set ("it")
		if len(set) == 0 {
			newBindings = mentalese.InitBindingSet(refBinding)
			break
		}

		if !base.quickAcceptabilityCheck(variable, ref.EntityType, set) {
			continue
		}

		testRangeBindings := base.solver.SolveRelationSet(set, mentalese.InitBindingSet(refBinding))
		if testRangeBindings.GetLength() == 1 {
			newBindings = testRangeBindings
			break
		}
	}

	return newBindings
}

func (base *SystemNestedStructureBase) SolveDefiniteReference(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	variable := relation.Arguments[0].TermValue
	set := relation.Arguments[1].TermValueRelationSet

	newBindings := base.SolveBackReference(relation, binding)

	if newBindings.IsEmpty() {
		newBindings = base.solver.SolveRelationSet(set, mentalese.InitBindingSet(binding))

		if newBindings.GetLength() > 1 {
			rangeIndex, found := base.rangeIndexClarification(newBindings, variable)
			if found {
				newBindings = mentalese.InitBindingSet(newBindings.Get(rangeIndex))
			} else {
				return mentalese.NewBindingSet()
			}
		}
	}

	return newBindings
}