package function

import (
	"nli-go/lib/api"
	"nli-go/lib/central"
	"nli-go/lib/mentalese"
)

const contextVariableAtom = "$$context$$main"

// go:context_set(time, P1, $time_modifier)
func (base *SystemSolverFunctionBase) contextSet(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	slotName := bound.Arguments[0].TermValue
	mainEntityVar := bound.Arguments[1]
	relations := bound.Arguments[2].TermValueRelationSet

	boundRelations := relations.ReplaceTerm(mainEntityVar, mentalese.NewTermAtom(contextVariableAtom))

	if slotName == central.DeixisTime {
		base.deicticCenter.SetTime(boundRelations)
	}

	return mentalese.InitBindingSet(binding)
}

// go:context_extend(time, P1, $time_modifier)
func (base *SystemSolverFunctionBase) contextExtend(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	slotName := bound.Arguments[0].TermValue
	mainEntityVar := bound.Arguments[1]
	relations := bound.Arguments[2].TermValueRelationSet

	slotRelations := mentalese.RelationSet{}

	if slotName == central.DeixisTime {
		slotRelations = base.deicticCenter.GetTime()
	}

	boundRelations := relations.ReplaceTerm(mainEntityVar, mentalese.NewTermAtom(contextVariableAtom))

	if slotName == central.DeixisTime {
		base.deicticCenter.SetTime(slotRelations.Merge(boundRelations))
	}

	return mentalese.InitBindingSet(binding)
}

// go:context_clear(time)
func (base *SystemSolverFunctionBase) contextClear(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	slotName := bound.Arguments[0].TermValue

	if slotName == central.DeixisTime {
		base.deicticCenter.SetTime(mentalese.RelationSet{})
	}

	return mentalese.InitBindingSet(binding)
}

// go:context_get(time, P1, Time)
func (base *SystemSolverFunctionBase) contextCall(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	slotName := bound.Arguments[0].TermValue
	mainEntityVar := bound.Arguments[1]

	slotRelations := mentalese.RelationSet{}

	if slotName == central.DeixisTime {
		slotRelations = base.deicticCenter.GetTime()
	}

	unboundRelations := slotRelations.ReplaceTerm(mentalese.NewTermAtom(contextVariableAtom), mainEntityVar)

	cursor := messenger.GetCursor()
	state := cursor.GetState("state", 0)
	cursor.SetState("state", 1)

	newBindings := mentalese.NewBindingSet()

	if state == 0 {
		messenger.CreateChildStackFrame(unboundRelations, mentalese.InitBindingSet(binding))
	} else {
		newBindings = cursor.GetChildFrameResultBindings()
	}

	return newBindings
}

func (base *SystemSolverFunctionBase) dialogReadBindings(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	someBindingVar := input.Arguments[0].TermValue

	responseBinding := (*base.discourseEntities).Copy()

	newBinding := binding.Copy()
	newBinding.Set(someBindingVar, mentalese.NewTermJson(responseBinding.ToRaw()))
	return mentalese.InitBindingSet(newBinding)
}

func (base *SystemSolverFunctionBase) dialogWriteBindings(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
	bound := input.BindSingle(binding)

	someBindings := mentalese.NewBindingSet()
	someBindingsRaw := []map[string]mentalese.Term{}
	bound.Arguments[0].GetJsonValue(&someBindingsRaw)
	someBindings.FromRaw(someBindingsRaw)

	// todo multiple bindings should not be merged into a single binding
	for _, someBinding := range someBindings.GetAll() {
		for key, value := range someBinding.GetAll() {
			base.discourseEntities.Set(key, value)
		}
	}

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) dialogAnaphoraQueueLast(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)
	variable := bound.Arguments[0].TermValue
	description := bound.Arguments[1].TermValueRelationSet
	loading := false

	cursor := messenger.GetCursor()
	cursor.SetState("childIndex", 0)

	newBindings := mentalese.NewBindingSet()

	for _, group := range *base.anaphoraQueue {

		if len(group) > 1 { continue }

		ref := group[0]

		//b := binding.Copy()
		b := mentalese.NewBinding()
		b.Set(variable, mentalese.NewTermId(ref.Id, ref.Sort))

		testRangeBindings := mentalese.BindingSet{}
		testRangeBindings, loading = messenger.ExecuteChildStackFrameAsync(description, mentalese.InitBindingSet(b))
		if loading {
			return mentalese.NewBindingSet()
		}
		if testRangeBindings.GetLength() > 0 {
			newBindings = testRangeBindings
			break
		}
	}

	return newBindings
}