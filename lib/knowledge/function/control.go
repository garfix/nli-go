package function

import (
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"os/exec"
	"strconv"
)

func (base *SystemSolverFunctionBase) let(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "**", base.log) { return mentalese.NewBindingSet() }

	variable := relation.Arguments[0].TermValue
	value := bound.Arguments[1]
	variables := base.solver.GetCurrentScope().GetVariables()

	if !relation.Arguments[0].IsVariable() {
		base.log.AddError("Let: variable already in use. Value: " + variable)
		return mentalese.NewBindingSet()
	}

	(*variables).Set(variable, value)

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) ifThenElse(ifThenElse mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	condition := ifThenElse.Arguments[0].TermValueRelationSet
	action := ifThenElse.Arguments[1].TermValueRelationSet
	alternative := ifThenElse.Arguments[2].TermValueRelationSet

	newBindings := base.solver.SolveRelationSet(condition, mentalese.InitBindingSet(binding))

	if !newBindings.IsEmpty() {
		newBindings = base.solver.SolveRelationSet(action, newBindings )
	} else {
		newBindings = base.solver.SolveRelationSet(alternative, mentalese.InitBindingSet(binding))
	}

	return newBindings
}

func (base *SystemSolverFunctionBase) call(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	child := relation.Arguments[0].TermValueRelationSet

	newBindings := base.solver.SolveRelationSet(child, mentalese.InitBindingSet(binding))

	return newBindings
}

func (base *SystemSolverFunctionBase) ignore(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	child := relation.Arguments[0].TermValueRelationSet

	newBindings := base.solver.SolveRelationSet(child, mentalese.InitBindingSet(binding))

	if newBindings.IsEmpty() {
		return mentalese.InitBindingSet(binding)
	} else {
		return newBindings
	}
}

func (base *SystemSolverFunctionBase) rangeForEach(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

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

	for i := start; i <= end; i++ {
		scopedBinding := binding.Copy()
		if !variableTerm.IsAnonymousVariable() {
			scopedBinding.Set(variable, mentalese.NewTermString(strconv.Itoa(i)))
		}
		elementBindings := base.solver.SolveRelationSet(children, mentalese.InitBindingSet(scopedBinding))
		if !variableTerm.IsAnonymousVariable() {
			elementBindings = elementBindings.FilterOutVariablesByName([]string{ variable })
		}
		newBindings.AddMultiple(elementBindings)
		if base.solver.GetCurrentScope().IsBreaked() {
			scope.SetBreaked(false)
			break
		}
	}

	return newBindings
}

func (base *SystemSolverFunctionBase) doBreak(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "", base.log) { return mentalese.NewBindingSet() }

	base.solver.GetCurrentScope().SetBreaked(true)

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) exec(input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

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


func (base *SystemSolverFunctionBase) execResponse(input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

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
