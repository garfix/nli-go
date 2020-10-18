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
	newBinding.Set(listVariable, mentalese.NewTermList(orderedList))

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
		scopedBinding.Set(variable, element)
		elementBindings := base.solver.SolveRelationSet(scope, mentalese.Bindings{ scopedBinding })
		if len(elementBindings) == 0 {
			newBindings = mentalese.Bindings{}
			break
		}
		newBindings = append(newBindings, elementBindings...)
	}

	return newBindings
}