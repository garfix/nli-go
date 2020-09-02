package central

import (
	"nli-go/lib/mentalese"
	"strconv"
)

type SystemGenerator struct {
}

func NewSystemGenerator() *SystemGenerator {
	return &SystemGenerator{}
}

func (gen SystemGenerator) generate(template mentalese.Relation, bindings mentalese.Bindings) (mentalese.RelationSet, bool) {

	result := mentalese.RelationSet{}
	found := false

	predicate := template.Predicate

	if predicate == mentalese.PredicateMakeAnd {
		found = true
		result = gen.makeAnd(template, bindings)
	}

	return result, found
}

// Creates a tree of 'and' relations, that grows to the right, and connects the entities from the first variable
func (gen SystemGenerator) makeAnd(template mentalese.Relation, bindings mentalese.Bindings) mentalese.RelationSet {

	result := mentalese.RelationSet{}
	entityVar := template.Arguments[0].TermValue
	rootTerm := template.Arguments[1]

	parentValue := rootTerm

	for i := 0; i < len(bindings)-2; i++ {

		binding := bindings[i]
		rightValue := mentalese.Term{TermType: mentalese.TermTypeVariable, TermValue: rootTerm.TermValue + strconv.Itoa(i+1)}

		relation := mentalese.NewRelation(true, mentalese.PredicateAnd, []mentalese.Term{
			parentValue,
			binding[entityVar],
			rightValue,
		})

		result = append(result, relation)
		parentValue = rightValue
	}

	if len(bindings) > 1 {

		beforeLastBinding := bindings[len(bindings)-2]
		lastBinding := bindings[len(bindings)-1]

		relation := mentalese.NewRelation(true, mentalese.PredicateAnd, []mentalese.Term{
			parentValue,
			beforeLastBinding[entityVar],
			lastBinding[entityVar],
		})

		result = append(result, relation)

	} else if len(bindings) == 1 {

		onlyBinding := bindings[0]

		relation := mentalese.NewRelation(true,mentalese.PredicateAnd, []mentalese.Term{
			parentValue,
			onlyBinding[entityVar],
			onlyBinding[entityVar],
		})

		result = append(result, relation)
	}

	return result
}
