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
