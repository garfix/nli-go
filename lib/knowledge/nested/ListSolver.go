package nested

import (
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"strconv"
)

func (base *SystemNestedStructureBase) SolveListOrder(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "lav", base.log) { return mentalese.NewBindingSet() }

	list := bound.Arguments[0].TermValueList
	orderFunction := bound.Arguments[1].TermValue
	listVariable := bound.Arguments[2].TermValue

	orderedList := base.entityQuickSort(list, orderFunction)

	newBinding := binding.Copy()
	newBinding.Set(listVariable, mentalese.NewTermList(orderedList))

	return mentalese.InitBindingSet(newBinding)
}

func (base *SystemNestedStructureBase) SolveListForeach(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)
	newBindings := mentalese.NewBindingSet()

	if len(relation.Arguments) == 3 {

		list := bound.Arguments[0].TermValueList
		elementVar := relation.Arguments[1].TermValue
		scope := relation.Arguments[2].TermValueRelationSet

		for _, element := range list {
			scopedBinding := binding.Copy()
			scopedBinding.Set(elementVar, element)
			elementBindings := base.solver.SolveRelationSet(scope, mentalese.InitBindingSet(scopedBinding))
			newBindings.AddMultiple(elementBindings)
		}

	} else if len(relation.Arguments) == 4 {

		list := bound.Arguments[0].TermValueList
		indexVar := relation.Arguments[1].TermValue
		elementVar := relation.Arguments[2].TermValue
		scope := relation.Arguments[3].TermValueRelationSet

		for index, element := range list {
			scopedBinding := binding.Copy()
			scopedBinding.Set(indexVar, mentalese.NewTermString(strconv.Itoa(index)))
			scopedBinding.Set(elementVar, element)
			elementBindings := base.solver.SolveRelationSet(scope, mentalese.InitBindingSet(scopedBinding))
			elementBindings = elementBindings.FilterOutVariablesByName([]string{indexVar, elementVar})
			newBindings.AddMultiple(elementBindings)
		}
	}

	return newBindings
}

func (base *SystemNestedStructureBase) listDeduplicate(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "lv", base.log) { return mentalese.NewBindingSet() }

	list := bound.Arguments[0].TermValueList
	newlistVar := bound.Arguments[1].TermValue

	newList := list.Deduplicate()

	newBinding := binding.Copy()
	newBinding.Set(newlistVar, mentalese.NewTermList(newList))
	return mentalese.InitBindingSet(newBinding)
}

func (base *SystemNestedStructureBase) listSort(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "lv", base.log) { return mentalese.NewBindingSet() }

	list := bound.Arguments[0].TermValueList
	newlistVar := bound.Arguments[1].TermValue

	newList, ok := list.Sort()
	if !ok {
		base.log.AddError("Could not sort list (must be strings or integers): " + list.String())
		return mentalese.NewBindingSet()
	}

	newBinding := binding.Copy()
	newBinding.Set(newlistVar, mentalese.NewTermList(newList))
	return mentalese.InitBindingSet(newBinding)
}

func (base *SystemNestedStructureBase) listIndex(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "l*v", base.log) { return mentalese.NewBindingSet() }

	list := bound.Arguments[0].TermValueList
	term := bound.Arguments[1]
	indexVar := bound.Arguments[2].TermValue

	newBindings := mentalese.NewBindingSet()

	for i, element := range list {
		if element.Equals(term) {
			newBinding := binding.Copy()
			newBinding.Set(indexVar, mentalese.NewTermString(strconv.Itoa(i)))
			newBindings.Add(newBinding)
		}
	}

	return newBindings
}

func (base *SystemNestedStructureBase) listGet(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "li*", base.log) { return mentalese.NewBindingSet() }

	list := bound.Arguments[0].TermValueList
	index := bound.Arguments[1].TermValue
	termVar := relation.Arguments[2].TermValue

	i, err := strconv.Atoi(index)
	if err != nil {
		base.log.AddError("Index should be an integer: " + index)
		return mentalese.NewBindingSet()
	}

	if i < 0 || i >= len(list) {
		return mentalese.NewBindingSet()
	}

	term := list[i]

	newBinding := binding.Copy()
	newBinding.Set(termVar, term)
	return mentalese.InitBindingSet(newBinding)
}

func (base *SystemNestedStructureBase) listLength(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "lv", base.log) { return mentalese.NewBindingSet() }

	list := bound.Arguments[0].TermValueList
	lengthVar := bound.Arguments[1].TermValue

	length := len(list)

	newBinding := binding.Copy()
	newBinding.Set(lengthVar, mentalese.NewTermString(strconv.Itoa(length)))
	return mentalese.InitBindingSet(newBinding)
}

func (base *SystemNestedStructureBase) listExpand(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "lv", base.log) { return mentalese.NewBindingSet() }

	list := bound.Arguments[0].TermValueList
	termVar := bound.Arguments[1].TermValue

	newBindings := mentalese.NewBindingSet()

	for _, element := range list {
		newBinding := binding.Copy()
		newBinding.Set(termVar, element)
		newBindings.Add(newBinding)
	}

	return newBindings
}