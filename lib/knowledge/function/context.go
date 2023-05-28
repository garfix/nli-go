package function

import (
	"nli-go/lib/api"
	"nli-go/lib/mentalese"
)

const contextVariableAtom = "$$context$$main"

// go:context_set(time, P1, $time_modifier)
func (base *SystemSolverFunctionBase) contextSet(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	slotName := bound.Arguments[0].TermValue
	firstEventVar := input.Arguments[1]
	secondEvent := bound.Arguments[2]
	relations := input.Arguments[3].TermValueRelationSet

	boundRelations := relations.ReplaceTerm(firstEventVar, mentalese.NewTermAtom(contextVariableAtom))

	if slotName == mentalese.DeixisTimeRelation {
		base.dialogContext.DeicticCenter.SetTimeEvent(secondEvent)
		base.dialogContext.DeicticCenter.SetTime(boundRelations)
	}

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) contextGet(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	slotName := bound.Arguments[0].TermValue
	variable := input.Arguments[1].TermValue

	newBinding := mentalese.NewBinding()

	value, found := base.dialogContext.DeicticCenter.Binding.Get(slotName)
	if found {
		newBinding.Set(variable, value)
	}

	return mentalese.InitBindingSet(newBinding)
}

// go:context_extend(time, P1, $time_modifier)
func (base *SystemSolverFunctionBase) contextExtend(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	slotName := bound.Arguments[0].TermValue
	mainEntityVar := bound.Arguments[1]
	relations := bound.Arguments[2].TermValueRelationSet

	slotRelations := mentalese.RelationSet{}

	if slotName == mentalese.DeixisTimeRelation {
		slotRelations = base.dialogContext.DeicticCenter.GetTime()
	}

	boundRelations := relations.ReplaceTerm(mainEntityVar, mentalese.NewTermAtom(contextVariableAtom))

	if slotName == mentalese.DeixisTimeRelation {
		base.dialogContext.DeicticCenter.SetTime(slotRelations.Merge(boundRelations))
	}

	return mentalese.InitBindingSet(binding)
}

// go:context_clear(time)
func (base *SystemSolverFunctionBase) contextClear(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	slotName := bound.Arguments[0].TermValue

	if slotName == mentalese.DeixisTimeRelation {
		base.dialogContext.DeicticCenter.SetTime(mentalese.RelationSet{})
	}

	return mentalese.InitBindingSet(binding)
}

// go:context_get(time, P1, Time)
func (base *SystemSolverFunctionBase) contextCall(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	slotName := bound.Arguments[0].TermValue
	mainEntityVar := input.Arguments[1]

	slotRelations := mentalese.RelationSet{}

	if slotName == mentalese.DeixisTimeRelation {
		slotRelations = base.dialogContext.DeicticCenter.GetTime()
	}

	unboundRelations := slotRelations.ReplaceTerm(mentalese.NewTermAtom(contextVariableAtom), mainEntityVar)

	newBindings := messenger.ExecuteChildStackFrame(unboundRelations, mentalese.InitBindingSet(binding))

	return newBindings
}

func (base *SystemSolverFunctionBase) createGoal(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	processType := bound.Arguments[0].TermValue
	set := bound.Arguments[1].TermValueRelationSet

	// add it to the list; run it (async); remove it from the list
	result := messenger.StartProcess(processType, set, binding)
	if !result {
		base.log.AddError("A process for " + processType + " is already active")
		return mentalese.NewBindingSet()
	}

	return mentalese.InitBindingSet(binding)
}
