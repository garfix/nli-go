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

func (gen SystemGenerator) generate(template mentalese.Relation, bindings []mentalese.Binding) (mentalese.RelationSet, bool) {

	result := mentalese.RelationSet{}
	found := false

	predicate := template.Predicate

	if predicate == "make_and" {
		found = true
		result = gen.makeAnd(template, bindings)
	}

	return result, found
}

// Creates a tree of 'and' relations, that grows to the right, and connects the entities from the first variable
func (gen SystemGenerator) makeAnd(template mentalese.Relation, bindings []mentalese.Binding) mentalese.RelationSet {

	result := mentalese.RelationSet{}
	entityVar := template.Arguments[0].TermValue
	rootVar := template.Arguments[1].TermValue

	parentValue := mentalese.Term{TermType: mentalese.Term_variable, TermValue: rootVar}

	for i := 0; i < len(bindings)-2; i++ {

		binding := bindings[i]
		rightValue := mentalese.Term{TermType: mentalese.Term_variable, TermValue: rootVar + strconv.Itoa(i+1)}

		relation := mentalese.Relation{Predicate: "and", Arguments: []mentalese.Term{
			parentValue,
			binding[entityVar],
			rightValue,
		}}

		result = append(result, relation)
		parentValue = rightValue
	}

	if len(bindings) > 1 {

		beforeLastBinding := bindings[len(bindings)-2]
		lastBinding := bindings[len(bindings)-1]

		relation := mentalese.Relation{Predicate: "and", Arguments: []mentalese.Term{
			parentValue,
			beforeLastBinding[entityVar],
			lastBinding[entityVar],
		}}

		result = append(result, relation)

	} else if len(bindings) == 1 {

		onlyBinding := bindings[0]

		relation := mentalese.Relation{Predicate: "and", Arguments: []mentalese.Term{
			parentValue,
			onlyBinding[entityVar],
			onlyBinding[entityVar],
		}}

		result = append(result, relation)
	}

	return result
}
