package nested

import (
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"strconv"
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

func (base *SystemNestedStructureBase) listDeduplicate(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "lv", base.log) {
		return mentalese.Bindings{}
	}

	list := bound.Arguments[0].TermValueList
	newlistVar := bound.Arguments[1].TermValue

	newList := list.Deduplicate()

	newBinding := binding.Copy()
	newBinding.Set(newlistVar, mentalese.NewTermList(newList))
	return mentalese.Bindings{ newBinding }
}

func (base *SystemNestedStructureBase) listSort(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "lv", base.log) {
		return mentalese.Bindings{}
	}

	list := bound.Arguments[0].TermValueList
	newlistVar := bound.Arguments[1].TermValue

	newList, ok := list.Sort()
	if !ok {
		base.log.AddError("Could not sort list (must be strings or integers): " + list.String())
		return mentalese.Bindings{}
	}

	newBinding := binding.Copy()
	newBinding.Set(newlistVar, mentalese.NewTermList(newList))
	return mentalese.Bindings{ newBinding }
}

func (base *SystemNestedStructureBase) listIndex(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "l*v", base.log) {
		return mentalese.Bindings{}
	}

	list := bound.Arguments[0].TermValueList
	term := bound.Arguments[1]
	indexVar := bound.Arguments[2].TermValue

	newBindings := mentalese.Bindings{}

	for i, element := range list {
		if element.Equals(term) {
			newBinding := binding.Copy()
			newBinding.Set(indexVar, mentalese.NewTermString(strconv.Itoa(i)))
			newBindings = append(newBindings, newBinding)
		}
	}

	return newBindings
}

func (base *SystemNestedStructureBase) listGet(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "li*", base.log) {
		return mentalese.Bindings{}
	}

	list := bound.Arguments[0].TermValueList
	index := bound.Arguments[1].TermValue
	termVar := bound.Arguments[2].TermValue

	i, err := strconv.Atoi(index)
	if err != nil {
		base.log.AddError("Index should be an integer: " + index)
		return mentalese.Bindings{}
	}

	if i < 0 || i >= len(list) {
		return mentalese.Bindings{}
	}

	term := list[i]

	newBinding := binding.Copy()
	newBinding.Set(termVar, term)
	return mentalese.Bindings{ newBinding }
}

func (base *SystemNestedStructureBase) listLength(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "lv", base.log) {
		return mentalese.Bindings{}
	}

	list := bound.Arguments[0].TermValueList
	lengthVar := bound.Arguments[1].TermValue

	length := len(list)

	newBinding := binding.Copy()
	newBinding.Set(lengthVar, mentalese.NewTermString(strconv.Itoa(length)))
	return mentalese.Bindings{ newBinding }
}