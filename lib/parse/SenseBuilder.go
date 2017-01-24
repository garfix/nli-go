package parse

import (
	"fmt"
	"nli-go/lib/mentalese"
)

type SenseBuilder struct {
	varIndexCounter map[string]int
}

func NewSenseBuilder() SenseBuilder {
	return SenseBuilder{varIndexCounter: map[string]int{}}
}

// Returns a new variable name
func (builder SenseBuilder) GetNewVariable(formalVariable string) string {

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
func (builder SenseBuilder) CreateVariableMap(actualAntecedent string, formalVariables []string) map[string]string {

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
				m[consequentVariable] = builder.GetNewVariable(consequentVariable)
			}
		}
	}

	return m
}


// Create actual relations given a set of templates and a variable map (formal to actual variables)
func (builder SenseBuilder) CreateGrammarRuleRelations(relationTemplates mentalese.RelationSet, variableMap map[string]string) mentalese.RelationSet {

	relations := mentalese.RelationSet{}

	for _, relation := range relationTemplates {
		for a, argument := range relation.Arguments {

			if argument.TermType == mentalese.Term_variable {

				relation.Arguments[a].TermType = mentalese.Term_variable
				relation.Arguments[a].TermValue = variableMap[argument.TermValue]
			}
		}

		relations = append(relations, relation)
	}

	return relations
}

// Create actual relations given a set of templates and an actual variable to replace any * positions
func (builder SenseBuilder) CreateLexItemRelations(relationTemplates mentalese.RelationSet, variable string) mentalese.RelationSet {

	relations := mentalese.RelationSet{}

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
