package function

import (
	"nli-go/lib/api"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
)

func (base *SystemSolverFunctionBase) intent(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !knowledge.Validate(bound, "a*", base.log) {
		return mentalese.NewBindingSet()
	}

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) backReference(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	variable := relation.Arguments[0].TermValue
	set := relation.Arguments[1].TermValueRelationSet

	newBindings := mentalese.NewBindingSet()

	for _, group := range *base.dialogContext.AnaphoraQueue {

		ref := group[0]

		b := mentalese.NewBinding()
		b.Set(variable, mentalese.NewTermId(ref.Id, ref.Sort))

		refBinding := binding.Merge(b)

		// empty set ("it")
		if len(set) == 0 {
			newBindings = mentalese.InitBindingSet(refBinding)
			break
		}

		if !base.quickAcceptabilityCheck(variable, ref.Sort, set) {
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

func (base *SystemSolverFunctionBase) definiteReference(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	variable := relation.Arguments[0].TermValue
	set := relation.Arguments[1].TermValueRelationSet

	newBindings := base.backReference(messenger, relation, binding)

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

func (base *SystemSolverFunctionBase) sortalBackReference(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	variable := relation.Arguments[0].TermValue

	newBindings := mentalese.NewBindingSet()

	for _, group := range *base.dialogContext.AnaphoraQueue {

		sort := ""

		for _, ref := range group {
			if sort == "" {
				sort = ref.Sort
			} else if sort != ref.Sort {
				sort = ""
				break
			}
		}

		if sort == "" {
			continue
		}

		sortInfo, found := base.meta.GetSortInfo(sort)
		if !found {
			continue
		}

		if sortInfo.Entity.Equals(mentalese.RelationSet{}) {
			continue
		}

		sortRelationSet := sortInfo.Entity.ReplaceTerm(mentalese.NewTermVariable(mentalese.IdVar), mentalese.NewTermVariable(variable))

		newBindings = base.solver.SolveRelationSet(sortRelationSet, mentalese.InitBindingSet(binding))
		break
	}

	return newBindings
}