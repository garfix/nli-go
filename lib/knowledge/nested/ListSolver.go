package nested

import (
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
)

func (base *SystemNestedStructureBase) SolveListOrder(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "lav", base.log) { return nil }

	list := bound.Arguments[0].TermValueList
	orderFunction := bound.Arguments[1].TermValue
	listVariable := bound.Arguments[2].TermValue

	orderedList := base.entityQuickSort(list, orderFunction)

	newBinding := binding.Copy()
	newBinding[listVariable] = mentalese.NewList(orderedList)

	return mentalese.Bindings{ newBinding }
}

func (base *SystemNestedStructureBase) SolveListForeach(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "lvr", base.log) { return nil }

	list := bound.Arguments[0].TermValueList
	variable := bound.Arguments[1].TermValue
	scope := bound.Arguments[2].TermValueRelationSet

	newBindings := mentalese.Bindings{}

	for _, element := range list {
		scopedBinding := binding.Copy()
		scopedBinding[variable] = element
		newBindings = append(newBindings, base.solver.SolveRelationSet(scope, mentalese.Bindings{ scopedBinding })...)
	}

	return newBindings
}