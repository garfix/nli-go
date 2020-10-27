package nested

import (
	"nli-go/lib/mentalese"
	"strconv"
)

func (base *SystemNestedStructureBase) Call(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	child := relation.Arguments[0].TermValueRelationSet

	newBindings := base.solver.SolveRelationSet(child, mentalese.InitBindingSet(binding))

	return newBindings
}

func (base *SystemNestedStructureBase) RangeForeach(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

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

func (base *SystemNestedStructureBase) quantOrderSingle(quant mentalese.Relation, orderFunction string) mentalese.RelationSet {

	orderedQuant := quant.Copy()

	if quant.Predicate == mentalese.PredicateQuant {
		for len(orderedQuant.Arguments) < 3 {
			orderedQuant.Arguments = append(orderedQuant.Arguments, mentalese.NewTermAnonymousVariable())
		}
		orderedQuant.Arguments[2] = mentalese.NewTermAtom(orderFunction)
	} else {
		leftQuant := orderedQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
		rightQuant := orderedQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]

		orderedQuant.Arguments[mentalese.SeqFirstOperandIndex] = mentalese.NewTermRelationSet( base.quantOrderSingle(leftQuant, orderFunction) )
		orderedQuant.Arguments[mentalese.SeqSecondOperandIndex] = mentalese.NewTermRelationSet( base.quantOrderSingle(rightQuant, orderFunction) )
	}

	return mentalese.RelationSet{ orderedQuant }
}