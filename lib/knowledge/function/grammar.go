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

	if messenger != nil {
		cursor := messenger.GetCursor()
		cursor.SetState("childIndex", 0)
	}

	result, _ := base.doBackReference(messenger, relation, binding)
	return result
}

func (base *SystemSolverFunctionBase) doBackReference(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) (mentalese.BindingSet, bool) {

	variable := relation.Arguments[0].TermValue
	set := relation.Arguments[1].TermValueRelationSet
	loading := false

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

		testRangeBindings := mentalese.BindingSet{}
		if messenger == nil {
			testRangeBindings = base.solver.SolveRelationSet(set, mentalese.InitBindingSet(refBinding))
		} else {
			testRangeBindings, loading = messenger.ExecuteChildStackFrameAsync(set, mentalese.InitBindingSet(refBinding))
		}
		if testRangeBindings.GetLength() == 1 {
			newBindings = testRangeBindings
			break
		}
	}

	return newBindings, loading
}

func (base *SystemSolverFunctionBase) definiteReference(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	variable := relation.Arguments[0].TermValue
	set := relation.Arguments[1].TermValueRelationSet

	if messenger != nil {
		cursor := messenger.GetCursor()
		cursor.SetState("childIndex", 0)
	}

	newBindings, loading := base.doBackReference(messenger, relation, binding)
	if loading { return mentalese.NewBindingSet() }

	if newBindings.IsEmpty() {
		if messenger == nil {
			newBindings = base.solver.SolveRelationSet(set, mentalese.InitBindingSet(binding))
		} else {
			newBindings, loading = messenger.ExecuteChildStackFrameAsync(set, mentalese.InitBindingSet(binding))
			if loading { return mentalese.NewBindingSet() }
		}

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
	loading := false

	if messenger != nil {
		cursor := messenger.GetCursor()
		cursor.SetState("childIndex", 0)
	}

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

		if messenger == nil {
			newBindings = base.solver.SolveRelationSet(sortRelationSet, mentalese.InitBindingSet(binding))
		} else {
			newBindings, loading = messenger.ExecuteChildStackFrameAsync(sortRelationSet, mentalese.InitBindingSet(binding))
			if loading { return mentalese.NewBindingSet() }
		}
		break
	}

	return newBindings
}