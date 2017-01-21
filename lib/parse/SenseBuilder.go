package parse

import (
	"fmt"
	"nli-go/lib/mentalese"
)

type SenseBuilder struct {
	varIndexCounter map[string]int
}

// Joins the senses of a parent node with those of its children.
//
// parentSense: declaration(S1) object(S1, E1)
// childSenses: [isa(E13, horse), isa(V24, fall) specifier(V24, S9) isa(S9, the)]
// childVariables: {NP: E13, VP: V24}
// return: isa(E13, horse) isa(V24, fall) declaration(V24) object(V24, E13)
func (builder SenseBuilder) Join(parentSense mentalese.RelationSet, childSenses []mentalese.RelationSet, childVariables map[string] string) (mentalese.RelationSet, string) {

	join := mentalese.RelationSet{}
	parentVariable := ""

	return join, parentVariable
}

// Returns a new variable name
func (builder SenseBuilder) getNewVariable(formalVariable string) string {

	initial := formalVariable[0:1]

	_, present := builder.varIndexCounter[initial]
	if !present {
		builder.varIndexCounter[initial] = 1
	} else {
		builder.varIndexCounter[initial]++
	}

	return fmt.Sprint(initial, builder.varIndexCounter[initial])
}

// Creates a map of formal variables to actual variables (new variables are created)
func (builder SenseBuilder) createVariableMap(actualAntecedent string, formalVariables []string) map[string]string {

	m := map[string]string{}
	antecedentVariable := formalVariables[0]

	for i := 1; i < len(formalVariables); i++ {

		consequentVariable := formalVariables[i]

		if consequentVariable == antecedentVariable {

			// the consequent variable matches the antecedent variable, inherit its actual variable
			m[consequentVariable] = actualAntecedent

		} else {

			// we're going to add a new actual variable, unless we already have
			_, present := m[consequentVariable]
			if !present {
				m[consequentVariable] = builder.getNewVariable(consequentVariable)
			}
		}
	}

	return m
}


// Create actual relations given a set of templates and a variable map (formal to actual variables)
func (builder SenseBuilder) createGrammarRuleRelations(relationTemplates []mentalese.Relation, variableMap map[string]string) []mentalese.Relation {

	relations := []mentalese.Relation{}

	for _, relation := range relationTemplates {
		for a, argument := range relation.Arguments {

			relation.Arguments[a].TermType = mentalese.Term_variable
			relation.Arguments[a].TermValue = variableMap[argument.TermValue]
		}

		relations = append(relations, relation)
	}

	return relations
}

// Create actual relations given a set of templates and an actual variable to replace any * positions
func (builder SenseBuilder) createLexItemRelations(relationTemplates []mentalese.Relation, variable string) []mentalese.Relation {

	relations := []mentalese.Relation{}

	for _, relationTemplate := range relationTemplates {

		relation := mentalese.Relation{}
		relation.Predicate = relationTemplate.Predicate

		for _, argument := range relationTemplate.Arguments {

			relationArgument := argument

			if argument.TermType == mentalese.Term_predicateAtom && argument.TermValue == "this" {

				relationArgument.TermType = mentalese.Term_variable
				relationArgument.TermValue = variable
			}

			relation.Arguments = append(relation.Arguments, relationArgument)
		}

		relations = append(relations, relation)
	}

	return relations
}
