package function

import (
	"nli-go/lib/api"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"os/exec"
	"strconv"
	"strings"
)

func (base *SystemSolverFunctionBase) assign(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "**", base.log) {
		return mentalese.NewBindingSet()
	}

	variables := []mentalese.Term{}
	values := []mentalese.Term{}

	if relation.Arguments[0].IsList() {
		variables = relation.Arguments[0].TermValueList
		values = bound.Arguments[1].TermValueList
	} else {
		variables = []mentalese.Term{relation.Arguments[0]}
		values = []mentalese.Term{bound.Arguments[1]}
	}

	for i, variable := range variables {
		value := values[i]

		if value.IsVariable() {
			base.log.AddError("Value of " + value.TermValue + " is unassigned")
			return mentalese.NewBindingSet()
		}

		if variable.IsMutableVariable() {
			messenger.SetMutableVariable(variable.TermValue, value)
		} else {
			existingValue, found := binding.Get(variable.TermValue)
			if found {
				if !existingValue.Equals(value) {
					base.log.AddError("Attempt to assign new value to " + variable.TermValue + "(" + existingValue.String() + " -> " + value.String() + ")")
					return mentalese.NewBindingSet()
				}
			}
		}
		binding.Set(variable.TermValue, value)
	}

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) append(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "v*", base.log) {
		return mentalese.NewBindingSet()
	}

	variable := relation.Arguments[0].TermValue
	value := bound.Arguments[1]

	if value.IsVariable() {
		base.log.AddError("Value of array " + variable + " is unassigned")
		return mentalese.NewBindingSet()
	}

	if !relation.Arguments[0].IsMutableVariable() {
		existingValue, found := binding.Get(variable)
		if found {
			if !existingValue.Equals(value) {
				base.log.AddError("Attempt to assign new value to " + variable + "(" + existingValue.String() + " -> " + value.String() + ")")
				return mentalese.NewBindingSet()
			}
		}
	}

	termValue, found := messenger.GetMutableVariable(variable)
	if !found {
		base.log.AddError("Attempt to append to initialized variable " + variable)
		return mentalese.NewBindingSet()
	}

	list := termValue.TermValueList
	list = append(list, value)
	newTermValue := mentalese.NewTermList((list))

	messenger.SetMutableVariable(variable, newTermValue)

	binding.Set(variable, newTermValue)

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) ifThen(messenger api.ProcessMessenger, ifThenElse mentalese.Relation,
	binding mentalese.Binding) mentalese.BindingSet {

	conditionArg := ifThenElse.Arguments[0]
	condition := ifThenElse.Arguments[0].TermValueRelationSet
	action := ifThenElse.Arguments[1].TermValueRelationSet

	var newBindings mentalese.BindingSet

	if conditionArg.IsAtom() {
		if conditionArg.TermValue == "true" {
			newBindings = messenger.ExecuteChildStackFrame(action, mentalese.InitBindingSet(binding))
		} else {
			newBindings = mentalese.InitBindingSet(binding)
		}
		return newBindings
	}

	conditionBindings := messenger.ExecuteChildStackFrame(condition, mentalese.InitBindingSet(binding))
	if conditionBindings.IsEmpty() {
		newBindings = mentalese.InitBindingSet(binding)
	} else {
		newBindings = messenger.ExecuteChildStackFrame(action, conditionBindings)
	}

	return newBindings
}

func (base *SystemSolverFunctionBase) ifThenElse(messenger api.ProcessMessenger, ifThenElse mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	conditionArg := ifThenElse.Arguments[0]
	condition := ifThenElse.Arguments[0].TermValueRelationSet
	action := ifThenElse.Arguments[1].TermValueRelationSet
	alternative := ifThenElse.Arguments[2].TermValueRelationSet

	var newBindings mentalese.BindingSet

	if conditionArg.IsAtom() {
		if conditionArg.TermValue == "true" {
			newBindings = messenger.ExecuteChildStackFrame(action, mentalese.InitBindingSet(binding))
		} else {
			newBindings = messenger.ExecuteChildStackFrame(alternative, mentalese.InitBindingSet(binding))
		}
		return newBindings
	}

	conditionBindings := messenger.ExecuteChildStackFrame(condition, mentalese.InitBindingSet(binding))
	if !conditionBindings.IsEmpty() {
		newBindings = messenger.ExecuteChildStackFrame(action, conditionBindings)
	} else {
		newBindings = messenger.ExecuteChildStackFrame(alternative, mentalese.InitBindingSet(binding))
	}

	return newBindings
}

func (base *SystemSolverFunctionBase) fail(messenger api.ProcessMessenger, ifThenElse mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	return mentalese.NewBindingSet()
}

func (base *SystemSolverFunctionBase) returnStatement(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "", base.log) {
		return mentalese.NewBindingSet()
	}

	messenger.AddProcessInstruction(mentalese.ProcessInstructionReturn, "")

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) call(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	child := relation.Arguments[0].TermValueRelationSet

	newBindings := messenger.ExecuteChildStackFrame(child, mentalese.InitBindingSet(binding))

	return newBindings
}

// apply(function, arg1, arg2, ...)
func (base *SystemSolverFunctionBase) apply(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	function := bound.Arguments[0].TermValueRelationSet[0]
	arguments := bound.Arguments[1:]

	functionBody := function.Arguments[len(function.Arguments)-1].TermValueRelationSet
	functionVariables := function.Arguments[0 : len(function.Arguments)-1]

	functionBinding := binding
	for i, variable := range functionVariables {
		functionBinding.Set(variable.TermValue, arguments[i])
	}

	boundFunctionBody := functionBody.BindSingle(functionBinding)

	result := messenger.ExecuteChildStackFrame(boundFunctionBody, mentalese.InitBindingSet(binding))
	return result
}

func (base *SystemSolverFunctionBase) ignore(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	child := relation.Arguments[0].TermValueRelationSet

	childBindings := messenger.ExecuteChildStackFrame(child, mentalese.InitBindingSet(binding))
	if childBindings.IsEmpty() {
		return mentalese.InitBindingSet(binding)
	} else {
		return childBindings
	}
}

func (base *SystemSolverFunctionBase) rangeForEach(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	startTerm := bound.Arguments[0].TermValue
	endTerm := bound.Arguments[1].TermValue
	variableTerm := relation.Arguments[2]
	variable := variableTerm.TermValue
	children := relation.Arguments[3].TermValueRelationSet

	newBindings := mentalese.NewBindingSet()

	start, err := strconv.Atoi(startTerm)
	if err != nil {
		return newBindings
	}

	end, err := strconv.Atoi(endTerm)
	if err != nil {
		return newBindings
	}

	cursor := messenger.GetCursor()
	cursor.SetType(mentalese.FrameTypeLoop)

	for index := start; index <= end; index++ {
		scopedBinding := binding.Copy()
		if !variableTerm.IsAnonymousVariable() {
			scopedBinding.Set(variable, mentalese.NewTermString(strconv.Itoa(index)))
		}
		childBindings := messenger.ExecuteChildStackFrame(children, mentalese.InitBindingSet(scopedBinding))
		newBindings.AddMultiple(childBindings)
	}

	return newBindings
}

func (base *SystemSolverFunctionBase) forRelations(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "rr", base.log) {
		return mentalese.NewBindingSet()
	}

	cursor := messenger.GetCursor()
	cursor.SetType(mentalese.FrameTypeLoop)

	forRelations := relation.Arguments[0].TermValueRelationSet
	bodyRelations := relation.Arguments[1].TermValueRelationSet

	forBindings := messenger.ExecuteChildStackFrameMutable(forRelations, mentalese.NewBinding())

	for _, forBinding := range forBindings.GetAll() {
		messenger.ExecuteChildStackFrame(bodyRelations, mentalese.InitBindingSet(forBinding))
	}

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) forIndexValue(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "vvlr", base.log) {
		return mentalese.NewBindingSet()
	}

	cursor := messenger.GetCursor()
	cursor.SetType(mentalese.FrameTypeLoop)

	indexVar := relation.Arguments[0].TermValue
	valueVar := relation.Arguments[1].TermValue
	list := bound.Arguments[2].TermValueList
	bodyRelations := relation.Arguments[3].TermValueRelationSet

	for i, element := range list {
		messenger.SetMutableVariable(indexVar, mentalese.NewTermString(strconv.Itoa(i)))
		messenger.SetMutableVariable(valueVar, element)
		messenger.ExecuteChildStackFrame(bodyRelations, mentalese.InitBindingSet(binding))
	}

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) forRange(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "viir", base.log) {
		return mentalese.NewBindingSet()
	}

	cursor := messenger.GetCursor()
	cursor.SetType(mentalese.FrameTypeLoop)

	elementVar := relation.Arguments[0].TermValue
	startValue, _ := bound.Arguments[1].GetIntValue()
	endValue, _ := bound.Arguments[2].GetIntValue()
	bodyRelations := relation.Arguments[3].TermValueRelationSet

	for i := startValue; i <= endValue; i++ {
		messenger.SetMutableVariable(elementVar, mentalese.NewTermString(strconv.Itoa(i)))
		messenger.ExecuteChildStackFrame(bodyRelations, mentalese.InitBindingSet(binding))
	}

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) listIndex2(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "liv", base.log) {
		return mentalese.NewBindingSet()
	}

	index, _ := bound.Arguments[1].GetIntValue()
	list := bound.Arguments[0].TermValueList
	returnVar := relation.Arguments[2].TermValue

	if index > len(list) {
		return mentalese.NewBindingSet()
	}
	if index < 0 {
		return mentalese.NewBindingSet()
	}

	element := list[index]
	newBinding := mentalese.NewBinding()
	newBinding.Set(returnVar, element)

	return mentalese.InitBindingSet(newBinding)
}

func (base *SystemSolverFunctionBase) doBreak(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "", base.log) {
		return mentalese.NewBindingSet()
	}

	messenger.AddProcessInstruction(mentalese.ProcessInstructionBreak, "")

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) cancel(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "", base.log) {
		return mentalese.NewBindingSet()
	}

	messenger.AddProcessInstruction(mentalese.ProcessInstructionCancel, "")

	return mentalese.NewBindingSet()
}

func (base *SystemSolverFunctionBase) waitFor(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	command := relation.Arguments[0].TermValue
	resultVar := relation.Arguments[1].TermValue

	bound := relation.BindSingle(binding)

	parameters := []interface{}{}
	for i, param := range bound.Arguments {
		if i > 1 {
			parameters = append(parameters, param.AsSimple())
		}
	}

	choiceKey := ""
	if command == mentalese.MessageChoose {
		choiceKey = bound.Arguments[2].TermValue
		for _, option := range bound.Arguments[3].TermValueList {
			choiceKey += "|" + option.TermValue
		}
	}

	if command == mentalese.MessageChoose {
		answer, found := base.choices[choiceKey]
		if found {
			binding.Set(resultVar, mentalese.NewTermString(answer))
			return mentalese.InitBindingSet(binding)
		}
	}

	base.clientConnector.SendToClient(messenger.GetProcess().GetType(), command, parameters)

	message := <-messenger.GetProcess().GetChannel()

	binding.Set(resultVar, mentalese.NewTermString(message.Message.(string)))

	if command == mentalese.MessageChoose {
		base.choices[choiceKey] = message.Message.(string)
	}

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) exec(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !knowledge.Validate(bound, "S", base.log) {
		return mentalese.NewBindingSet()
	}

	command := bound.Arguments[0].TermValue
	args := []string{}
	for i := range bound.Arguments {
		if i == 0 {
			continue
		}
		args = append(args, bound.Arguments[i].TermValue)
	}
	cmd := exec.Command(command, args...)
	_, err := cmd.Output()
	if err != nil {
		base.log.AddError(err.Error())
	}

	newBinding := binding.Copy()

	return mentalese.InitBindingSet(newBinding)
}

func (base *SystemSolverFunctionBase) execResponse(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)
	responseVar := input.Arguments[0].TermValue

	if !knowledge.Validate(bound, "vS", base.log) {
		return mentalese.NewBindingSet()
	}

	command := bound.Arguments[1].TermValue
	args := []string{}
	for i := range bound.Arguments {
		if i < 2 {
			continue
		}
		args = append(args, bound.Arguments[i].TermValue)
	}
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		base.log.AddError(err.Error())
	}

	newBinding := binding.Copy()

	newBinding.Set(responseVar, mentalese.NewTermString(strings.TrimSpace(string(output))))

	return mentalese.InitBindingSet(newBinding)
}

func (base *SystemSolverFunctionBase) slot(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	slotName := bound.Arguments[0].TermValue
	slotField := bound.Arguments[1]

	newBinding := mentalese.NewBinding()

	if slotField.IsVariable() {

		value, found := messenger.GetProcessSlot(slotName)
		if found {
			newBinding.Set(slotField.TermValue, value)
		} else {
			base.log.AddError("Slot not found: " + slotName)
		}

	} else {

		messenger.SetProcessSlot(slotName, slotField)

	}

	return mentalese.InitBindingSet(newBinding)
}
