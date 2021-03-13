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

	cursor := messenger.GetCursor()
	cursor.SetState("childIndex", 0)

	result, _ := base.doBackReference(messenger, relation, binding)
	return result
}

func (base *SystemSolverFunctionBase) doBackReference(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) (mentalese.BindingSet, bool) {

	variable := relation.Arguments[0].TermValue
	set := relation.Arguments[1].TermValueRelationSet
	loading := false

	newBindings := mentalese.NewBindingSet()

	for _, group := range *base.anaphoraQueue {

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
		testRangeBindings, loading = messenger.ExecuteChildStackFrameAsync(set, mentalese.InitBindingSet(refBinding))
		if loading {
			return mentalese.NewBindingSet(), true
		}
		if testRangeBindings.GetLength() == 1 {
			newBindings = testRangeBindings
			break
		}
	}

	return newBindings, loading
}

func (base *SystemSolverFunctionBase) definiteReference(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	set := relation.Arguments[1].TermValueRelationSet

	cursor := messenger.GetCursor()
	cursor.SetState("childIndex", 0)

	newBindings, loading := base.doBackReference(messenger, relation, binding)
	if loading { return mentalese.NewBindingSet() }

	if newBindings.IsEmpty() {
		newBindings, loading = messenger.ExecuteChildStackFrameAsync(set, mentalese.InitBindingSet(binding))
		if loading { return mentalese.NewBindingSet() }

		if newBindings.GetLength() > 1 {
			loading = base.rangeIndexClarification(messenger)
			newBindings = mentalese.NewBindingSet()
		}
	}

	return newBindings
}

// ask the user which of the specified entities he/she means
func (base *SystemSolverFunctionBase) rangeIndexClarification(messenger api.ProcessMessenger) bool {

	cursor := messenger.GetCursor()
	state := cursor.GetState("state", 0)
	cursor.SetState("state", state + 1)

	if state == 0 {

		set := mentalese.RelationSet{
			mentalese.NewRelation(true, mentalese.PredicateAssert, []mentalese.Term{
				mentalese.NewTermRelationSet(mentalese.RelationSet{
					mentalese.NewRelation(true, mentalese.PredicateOutput, []mentalese.Term{
						mentalese.NewTermString("I don't understand which one you mean"),
					})}),
			}),
		}
		messenger.CreateChildStackFrame(set, mentalese.InitBindingSet(mentalese.NewBinding()))
		return true

	} else {
		return false
	}
}

func (base *SystemSolverFunctionBase) sortalBackReference(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	variable := relation.Arguments[0].TermValue
	newBindings := mentalese.NewBindingSet()
	loading := false

	cursor := messenger.GetCursor()
	cursor.SetState("childIndex", 0)

	for _, group := range *base.anaphoraQueue {

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

		newBindings, loading = messenger.ExecuteChildStackFrameAsync(sortRelationSet, mentalese.InitBindingSet(binding))
		if loading { return mentalese.NewBindingSet() }
		break
	}

	return newBindings
}