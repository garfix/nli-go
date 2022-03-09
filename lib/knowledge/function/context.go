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
		base.dialogContext.DeicticCenter.SetTime(boundRelations)
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
		slotRelations = base.dialogContext.DeicticCenter.GetTime()
	}

	boundRelations := relations.ReplaceTerm(mainEntityVar, mentalese.NewTermAtom(contextVariableAtom))

	if slotName == central.DeixisTime {
		base.dialogContext.DeicticCenter.SetTime(slotRelations.Merge(boundRelations))
	}

	return mentalese.InitBindingSet(binding)
}

// go:context_clear(time)
func (base *SystemSolverFunctionBase) contextClear(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	slotName := bound.Arguments[0].TermValue

	if slotName == central.DeixisTime {
		base.dialogContext.DeicticCenter.SetTime(mentalese.RelationSet{})
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
		slotRelations = base.dialogContext.DeicticCenter.GetTime()
	}

	unboundRelations := slotRelations.ReplaceTerm(mentalese.NewTermAtom(contextVariableAtom), mainEntityVar)

	newBindings := messenger.ExecuteChildStackFrame(unboundRelations, mentalese.InitBindingSet(binding))

	return newBindings
}

func (base *SystemSolverFunctionBase) dialogReadBindings(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	someBindingVar := input.Arguments[0].TermValue

	responseBinding := (*base.dialogContext.DiscourseEntities).Copy()

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

	for _, someBinding := range someBindings.GetAll() {
		for key, value := range someBinding.GetAll() {
			if value.IsId() {
				existingValue, found := base.dialogContext.DiscourseEntities.Get(key)
				if found {
					if existingValue.IsList() {
						if !existingValue.ListContains(value) {
							list := existingValue.TermValueList
							list = append(list, value)
							base.dialogContext.DiscourseEntities.Set(key, mentalese.NewTermList(list))
						}
					} else if !existingValue.Equals(value) {
						list := []mentalese.Term{existingValue, value}
						base.dialogContext.DiscourseEntities.Set(key, mentalese.NewTermList(list))
					}
				} else {
					base.dialogContext.DiscourseEntities.Set(key, value)
				}
			}
		}
	}

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) createGoal(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	set := bound.Arguments[0].TermValueRelationSet

	// add it to the list; run it (async); remove it from the list
	messenger.StartProcess(set, binding)

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) dialogAddResponseClause(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
	bound := input.BindSingle(binding)

	essentialResponseBindings := mentalese.NewBindingSet()
	someBindingsRaw := []map[string]mentalese.Term{}
	bound.Arguments[0].GetJsonValue(&someBindingsRaw)
	essentialResponseBindings.FromRaw(someBindingsRaw)

	entities := []*mentalese.ClauseEntity{}
	for _, binding := range essentialResponseBindings.GetAll() {
		for _, variable := range binding.GetKeys() {
			entities = append(entities, mentalese.NewClauseEntity(variable, mentalese.AtomFunctionObject))
		}
	}

	clause := mentalese.NewClause(nil, true, entities)

	if len(entities) > 0 {
		clause.Center = entities[0]
	} else {
		previousClause := base.dialogContext.ClauseList.GetLastClause()
		if previousClause != nil {
			clause.Center = previousClause.Center
		}
	}

	base.dialogContext.ClauseList.AddClause(clause)

	return mentalese.InitBindingSet(binding)
}
