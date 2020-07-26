package nested

import (
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
)

func (base *SystemNestedStructureBase) Call(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	child := relation.Arguments[0].TermValueRelationSet

	newBindings := base.solver.SolveRelationSet(child, mentalese.Bindings{ binding })

	return newBindings
}

func (base *SystemNestedStructureBase) QuantOrder(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	bound := relation.BindSingle(binding)

	if !knowledge.Validate(bound, "rav", base.log) {
		return nil
	}

	quant := bound.Arguments[0].TermValueRelationSet[0]
	orderFunction := bound.Arguments[1].TermValue
	quantVariable := bound.Arguments[2].TermValue

	orderedQuant := base.quantOrderSingle(quant, orderFunction)

	newBinding := binding.Copy()
	newBinding[quantVariable] = mentalese.NewRelationSet(orderedQuant)

	return mentalese.Bindings{ newBinding }
}

func (base *SystemNestedStructureBase) quantOrderSingle(quant mentalese.Relation, orderFunction string) mentalese.RelationSet {

	orderedQuant := quant.Copy()

	if quant.Predicate == mentalese.PredicateQuant {
		for len(orderedQuant.Arguments) < 3 {
			orderedQuant.Arguments = append(orderedQuant.Arguments, mentalese.NewAnonymousVariable())
		}
		orderedQuant.Arguments[2] = mentalese.NewPredicateAtom(orderFunction)
	} else {
		leftQuant := orderedQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
		rightQuant := orderedQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]

		orderedQuant.Arguments[mentalese.SeqFirstOperandIndex] = mentalese.NewRelationSet( base.quantOrderSingle(leftQuant, orderFunction) )
		orderedQuant.Arguments[mentalese.SeqSecondOperandIndex] = mentalese.NewRelationSet( base.quantOrderSingle(rightQuant, orderFunction) )
	}

	return mentalese.RelationSet{ orderedQuant }
}