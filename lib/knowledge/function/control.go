package function

import (
	"nli-go/lib/api"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"os/exec"
	"strconv"
)

// todo: remove
func (base *SystemSolverFunctionBase) let(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "**", base.log) { return mentalese.NewBindingSet() }

	variable := relation.Arguments[0].TermValue
	value := bound.Arguments[1]
	variables := base.solver.GetCurrentScope().GetVariables()

	if !relation.Arguments[0].IsVariable() {
		base.log.AddError("Let: variable already in use. Value: " + variable)
		return mentalese.NewBindingSet()
	}

	if messenger == nil {
		(*variables).Set(variable, value)
	} else {
		messenger.AddProcessInstruction(mentalese.ProcessInstructionLet, variable)
		binding = binding.Copy()
		binding.Set(variable, value)
	}

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) ifThen(messenger api.ProcessMessenger, ifThenElse mentalese.Relation,
	binding mentalese.Binding) mentalese.BindingSet {

	condition := ifThenElse.Arguments[0].TermValueRelationSet
	action := ifThenElse.Arguments[1].TermValueRelationSet

	newBindings := mentalese.NewBindingSet()

	if messenger == nil {

		newBindings = base.solver.SolveRelationSet(condition, mentalese.InitBindingSet(binding))

		if !newBindings.IsEmpty() {
			newBindings = base.solver.SolveRelationSet(action, newBindings)
		} else {
			newBindings = mentalese.InitBindingSet(binding)
		}

	} else {

		cursor := messenger.GetCursor()
		state := cursor.GetState("state", 0)
		if state == 0 {

			cursor.SetState("state", 1)
			messenger.CreateChildStackFrame(condition, mentalese.InitBindingSet(binding))

		} else if state == 1 {

			cursor.SetState("state", 2)
			conditionBindings := cursor.GetChildFrameResultBindings()
			if conditionBindings.IsEmpty() {
				newBindings = mentalese.InitBindingSet(binding)
			} else {
				messenger.CreateChildStackFrame(action, conditionBindings)
			}

		} else {

			actionBindings := cursor.GetChildFrameResultBindings()
			newBindings = actionBindings

		}
	}

	return newBindings
}

func (base *SystemSolverFunctionBase) ifThenElse(messenger api.ProcessMessenger, ifThenElse mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	condition := ifThenElse.Arguments[0].TermValueRelationSet
	action := ifThenElse.Arguments[1].TermValueRelationSet
	alternative := ifThenElse.Arguments[2].TermValueRelationSet

	newBindings := mentalese.NewBindingSet()

	if messenger == nil {

		newBindings := base.solver.SolveRelationSet(condition, mentalese.InitBindingSet(binding))

		if !newBindings.IsEmpty() {
			newBindings = base.solver.SolveRelationSet(action, newBindings)
		} else {
			newBindings = base.solver.SolveRelationSet(alternative, mentalese.InitBindingSet(binding))
		}

	} else {

		cursor := messenger.GetCursor()
		state := cursor.GetState("state", 0)
		cursor.SetState("state", state + 1)

		if state == 0 {

			messenger.CreateChildStackFrame(condition, mentalese.InitBindingSet(binding))

		} else if state == 1 {

			conditionBindings := cursor.GetChildFrameResultBindings()
			if !conditionBindings.IsEmpty() {
				messenger.CreateChildStackFrame(action, conditionBindings)
			} else {
				messenger.CreateChildStackFrame(alternative, mentalese.InitBindingSet(binding))
			}

		} else {

			newBindings = cursor.GetChildFrameResultBindings()

		}

	}

	return newBindings
}

func (base *SystemSolverFunctionBase) fail(messenger api.ProcessMessenger, ifThenElse mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	return mentalese.NewBindingSet()
}

func (base *SystemSolverFunctionBase) call(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	child := relation.Arguments[0].TermValueRelationSet

	newBindings := mentalese.NewBindingSet()

	if messenger == nil {

		newBindings = base.solver.SolveRelationSet(child, mentalese.InitBindingSet(binding))

	} else {

		cursor := messenger.GetCursor()
		state := cursor.GetState("state", 0)
		cursor.SetState("state", 1)

		if state == 0 {
			messenger.CreateChildStackFrame(child, mentalese.InitBindingSet(binding))
		} else {
			newBindings = cursor.GetChildFrameResultBindings()
		}

	}

	return newBindings
}

func (base *SystemSolverFunctionBase) ignore(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	child := relation.Arguments[0].TermValueRelationSet

	if messenger == nil {

		newBindings := base.solver.SolveRelationSet(child, mentalese.InitBindingSet(binding))

		if newBindings.IsEmpty() {
			return mentalese.InitBindingSet(binding)
		} else {
			return newBindings
		}

	} else {
		cursor := messenger.GetCursor()
		state := cursor.GetState("state", 0)
		if state == 0 {
			cursor.SetState("state", 1)
			messenger.CreateChildStackFrame(child, mentalese.InitBindingSet(binding))
		} else {
			childBindings := messenger.GetCursor().GetChildFrameResultBindings()
			if childBindings.IsEmpty() {
				return mentalese.InitBindingSet(binding)
			} else {
				return childBindings
			}
		}
	}

	return mentalese.NewBindingSet()
}

func (base *SystemSolverFunctionBase) rangeForEach(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	startTerm := bound.Arguments[0].TermValue
	endTerm := bound.Arguments[1].TermValue
	variableTerm := relation.Arguments[2]
	variable := variableTerm.TermValue
	children := relation.Arguments[3].TermValueRelationSet
	scope := base.solver.GetCurrentScope()

	newBindings := mentalese.NewBindingSet()

	start, err := strconv.Atoi(startTerm)
	if err != nil {
		return newBindings
	}

	end, err := strconv.Atoi(endTerm)
	if err != nil {
		return newBindings
	}

	if messenger == nil {

		for i := start; i <= end; i++ {
			scopedBinding := binding.Copy()
			if !variableTerm.IsAnonymousVariable() {
				scopedBinding.Set(variable, mentalese.NewTermString(strconv.Itoa(i)))
			}
			elementBindings := base.solver.SolveRelationSet(children, mentalese.InitBindingSet(scopedBinding))
			if !variableTerm.IsAnonymousVariable() {
				elementBindings = elementBindings.FilterOutVariablesByName([]string{variable})
			}
			newBindings.AddMultiple(elementBindings)
			if base.solver.GetCurrentScope().IsBreaked() {
				scope.SetBreaked(false)
				break
			}
		}
	} else {

		cursor := messenger.GetCursor()
		index := cursor.GetState("index", start)
		cursor.SetState("index", index + 1)

		if index == start {
			cursor.SetType(mentalese.FrameTypeLoop)
		} else {
			newBindings.AddMultiple(cursor.GetChildFrameResultBindings())
		}

		if index <= end {
			scopedBinding := binding.Copy()
			if !variableTerm.IsAnonymousVariable() {
				scopedBinding.Set(variable, mentalese.NewTermString(strconv.Itoa(index)))
			}
			messenger.CreateChildStackFrame(children, mentalese.InitBindingSet(scopedBinding))
		}
	}

	return newBindings
}

func (base *SystemSolverFunctionBase) doBreak(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "", base.log) { return mentalese.NewBindingSet() }

	if messenger == nil {
		base.solver.GetCurrentScope().SetBreaked(true)
	} else {
		messenger.AddProcessInstruction(mentalese.ProcessInstructionBreak, "")
	}

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) waitFor(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	child := relation.Arguments[0].TermValueRelationSet

	newBindings := mentalese.NewBindingSet()

	cursor := messenger.GetCursor()
	state := cursor.GetState("state", 0)
	cursor.SetState("state", 1)

	if state == 0 {
		messenger.CreateChildStackFrame(child, mentalese.InitBindingSet(binding))
	} else {
		newBindings = cursor.GetChildFrameResultBindings()
		if newBindings.IsEmpty() {
			messenger.CreateChildStackFrame(child, mentalese.InitBindingSet(binding))
			messenger.AddProcessInstruction(mentalese.ProcessInstructionStop, "")
		}
	}

	return newBindings
}

func (base *SystemSolverFunctionBase) exec(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !knowledge.Validate(bound, "S", base.log) {
		return mentalese.NewBindingSet()
	}

	command := bound.Arguments[0].TermValue
	args := []string{}
	for i := range bound.Arguments {
		if i == 0 { continue }
		args = append(args, bound.Arguments[i].TermValue)
	}
	cmd := exec.Command(command, args...)
	_, err := cmd.Output()
	if err != nil {
		base.log.AddError(err.Error())
	}

	newBinding := binding.Copy()

	return mentalese.InitBindingSet( newBinding )
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
		if i < 2 { continue }
		args = append(args, bound.Arguments[i].TermValue)
	}
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		base.log.AddError(err.Error())
	}

	newBinding := binding.Copy()

	newBinding.Set(responseVar, mentalese.NewTermString(string(output)))

	return mentalese.InitBindingSet( newBinding )
}
