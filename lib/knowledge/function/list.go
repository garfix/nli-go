package function

import (
	"nli-go/lib/api"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"strconv"
)

func (base *SystemSolverFunctionBase) listOrder(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "lav", base.log) { return mentalese.NewBindingSet() }

	cursor := messenger.GetCursor()
	cursor.SetState("childIndex", 0)

	list := bound.Arguments[0].TermValueList.Copy()
	orderFunction := bound.Arguments[1].TermValue
	listVariable := bound.Arguments[2].TermValue

	orderedList, loading := base.entityQuickSort(messenger, list, orderFunction)
	if loading {
		return mentalese.NewBindingSet()
	}

	newBinding := binding.Copy()
	newBinding.Set(listVariable, mentalese.NewTermList(orderedList))

	return mentalese.InitBindingSet(newBinding)
}

func (base *SystemSolverFunctionBase) listAppend(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "l*v", base.log) { return mentalese.NewBindingSet() }

	list := bound.Arguments[0].TermValueList
	element := bound.Arguments[1]
	listVariable := relation.Arguments[2].TermValue

	newList := list.Copy()
	newList = newList.Append(element)

	newBinding := binding.Copy()
	newBinding.Set(listVariable, mentalese.NewTermList(newList))
	return mentalese.InitBindingSet(newBinding)
}

func (base *SystemSolverFunctionBase) listForeach(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)
	newBindings := mentalese.NewBindingSet()

	cursor := messenger.GetCursor()
	index := cursor.GetState("index", 0)
	cursor.SetState("index", index + 1)

	if len(relation.Arguments) == 3 {

		list := bound.Arguments[0].TermValueList
		elementVar := relation.Arguments[1].TermValue
		children := relation.Arguments[2].TermValueRelationSet

		if index == 0 {
			cursor.SetType(mentalese.FrameTypeLoop)
		} else {
			newBindings.AddMultiple(cursor.GetChildFrameResultBindings())
		}

		if index < len(list) {

			element := list[index]

			scopedBinding := binding.Copy()
			scopedBinding.Set(elementVar, element)

			messenger.CreateChildStackFrame(children, mentalese.InitBindingSet(scopedBinding))
		}

	} else if len(relation.Arguments) == 4 {

		list := bound.Arguments[0].TermValueList
		indexVar := relation.Arguments[1].TermValue
		elementVar := relation.Arguments[2].TermValue
		children := relation.Arguments[3].TermValueRelationSet

		if index == 0 {
			cursor.SetType(mentalese.FrameTypeLoop)
		} else {
			newBindings.AddMultiple(cursor.GetChildFrameResultBindings())
		}

		if index < len(list) {

			element := list[index]

			scopedBinding := binding.Copy()
			scopedBinding.Set(indexVar, mentalese.NewTermString(strconv.Itoa(index)))
			scopedBinding.Set(elementVar, element)

			messenger.CreateChildStackFrame(children, mentalese.InitBindingSet(scopedBinding))
		}
	}

	return newBindings
}

func (base *SystemSolverFunctionBase) listDeduplicate(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "lv", base.log) { return mentalese.NewBindingSet() }

	list := bound.Arguments[0].TermValueList
	newlistVar := bound.Arguments[1].TermValue

	newList := list.Deduplicate()

	newBinding := binding.Copy()
	newBinding.Set(newlistVar, mentalese.NewTermList(newList))
	return mentalese.InitBindingSet(newBinding)
}

func (base *SystemSolverFunctionBase) listSort(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

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

func (base *SystemSolverFunctionBase) listIndex(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

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

func (base *SystemSolverFunctionBase) listExpand(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

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