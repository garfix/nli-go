package function

import (
	"nli-go/lib/api"
	"nli-go/lib/central"
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

	unscopedSense := base.getSense(messenger).UnScope()

	if base.discourseEntities.ContainsVariable(variable) {
		newBinding := mentalese.NewBinding()
		newBinding.Set(variable, base.discourseEntities.MustGet(variable))
		return mentalese.InitBindingSet(newBinding), false
	}

	for _, group := range *base.anaphoraQueue {

		ref := group[0]

		b := mentalese.NewBinding()
		b.Set(variable, mentalese.NewTermId(ref.Id, ref.Sort))

		refBinding := binding.Merge(b)

		if base.isReflexive(unscopedSense, variable, ref) {
			continue
		}

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

func (base *SystemSolverFunctionBase) getSense(messenger api.ProcessMessenger) mentalese.RelationSet {
	term, found := messenger.GetProcessSlot(mentalese.SlotSense)
	if !found {
		base.log.AddError("Slot 'sense' not found")
	}

	return term.TermValueRelationSet
}

// checks if a (irreflexive) pronoun does not refer to another element in a same relation
func (base *SystemSolverFunctionBase) isReflexive(unscopedSense mentalese.RelationSet, referenceVariable string, antecedent central.EntityReference) bool {

	antecedentvariable := antecedent.Variable

	if antecedentvariable == "" {
		return false
	}

	reflexive := false
	for _, relation := range unscopedSense {
		ref := false
		ante := false
		for _, argument := range relation.Arguments {
			if argument.IsVariable() {
				if argument.TermValue == antecedentvariable {
					ante = true
				}
				if argument.TermValue == referenceVariable {
					ref = true
				}
			}
		}
		if ref && ante {
			reflexive = true
		}
	}

	return reflexive
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
			mentalese.NewRelation(false, mentalese.PredicateAssert, []mentalese.Term{
				mentalese.NewTermRelationSet(mentalese.RelationSet{
					mentalese.NewRelation(false, mentalese.PredicateOutput, []mentalese.Term{
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