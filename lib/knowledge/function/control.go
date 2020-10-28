package function

import (
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"strconv"
)

func (base *SystemSolverFunctionBase) let(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "**", base.log) { return mentalese.NewBindingSet() }

	variable := relation.Arguments[0].TermValue
	value := bound.Arguments[1]
	variables := base.solver.GetCurrentScope().GetVariables()
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

func (base *SystemSolverFunctionBase) rangeForEach(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

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
	}

	return newBindings
}