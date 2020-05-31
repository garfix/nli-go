package parse

import (
	"fmt"
	"nli-go/lib/mentalese"
)

type SenseBuilder struct {
	varIndexCounter map[string]int
	constantCounter int
}

func NewSenseBuilder() SenseBuilder {
	return SenseBuilder{varIndexCounter: map[string]int{}, constantCounter: 1}
}

// Returns a new variable name
func (builder SenseBuilder) GetNewVariable(formalVariable string) string {

	initial := formalVariable[0:1]

	_, present := builder.varIndexCounter[initial]
	if !present {
		builder.varIndexCounter[initial] = 5
	} else {
		builder.varIndexCounter[initial]++
	}

	return fmt.Sprint(initial, builder.varIndexCounter[initial])
}

// Creates a map of formal variables to actual variables (new variables are created)
func (builder SenseBuilder) CreateVariableMap(actualAntecedents []string, antecedentVariables []string, allConsequentVariables [][]string) map[string]mentalese.Term {

	m := map[string]mentalese.Term{}

	if len(actualAntecedents) == 0 {
		return m
	}

	for i, antecedentVariable := range antecedentVariables {
		m[antecedentVariable] = mentalese.NewVariable(actualAntecedents[i])
	}

	for _, consequentVariables := range allConsequentVariables {
		for _, consequentVariable := range consequentVariables {

			_, present := m[consequentVariable]
			if !present {
				m[consequentVariable] = mentalese.NewVariable(builder.GetNewVariable(consequentVariable))
			}
		}
	}

	return m
}

func (builder SenseBuilder) ExtendVariableMap(sense mentalese.RelationSet, variableMap map[string]mentalese.Term) map[string]mentalese.Term {

	for _, relation := range sense {
		for _, argument := range relation.Arguments {
			if argument.IsVariable() {
				variable := argument.TermValue
				_, present := variableMap[variable]
				if !present {
					variableMap[variable] = mentalese.NewVariable(builder.GetNewVariable(variable))
				}
			} else if argument.IsRelationSet() {
				childMap := builder.ExtendVariableMap(argument.TermValueRelationSet, variableMap)
				for key, value := range childMap {
					variableMap[key] = value
				}
			}
		}
	}

	return variableMap
}

// Create actual relations given a set of templates and a variable map (formal to actual variables)
func (builder SenseBuilder) CreateGrammarRuleRelations(relationTemplates mentalese.RelationSet, variableMap map[string]mentalese.Term) mentalese.RelationSet {

	relations := mentalese.RelationSet{}

	for _, relation := range relationTemplates {

		newRelation := relation.Copy()

		for a, argument := range newRelation.Arguments {

			if argument.TermType == mentalese.TermVariable {

				newRelation.Arguments[a].TermType = variableMap[argument.TermValue].TermType
				newRelation.Arguments[a].TermValue = variableMap[argument.TermValue].TermValue

			} else if argument.TermType == mentalese.TermRelationSet {

				newRelation.Arguments[a].TermType = mentalese.TermRelationSet
				newRelation.Arguments[a].TermValueRelationSet = builder.CreateGrammarRuleRelations(argument.TermValueRelationSet, variableMap)
			}
		}

		relations = append(relations, newRelation)
	}

	return relations
}

// Create actual relations given a set of templates and an actual variable to replace any * positions
func (builder SenseBuilder) CreateLexItemRelations(relationTemplates mentalese.RelationSet, variable string) mentalese.RelationSet {

	from := mentalese.Term{TermType: mentalese.TermVariable, TermValue: "E"}
	to := mentalese.Term{TermType: mentalese.TermVariable, TermValue: variable}

	return builder.ReplaceTerm(relationTemplates, from, to)
}

// Replaces all occurrences in relationTemplates from from to to
func (builder SenseBuilder) ReplaceTerm(relationTemplates mentalese.RelationSet, from mentalese.Term, to mentalese.Term) mentalese.RelationSet {

	relations := mentalese.RelationSet{}

	for _, relationTemplate := range relationTemplates {

		arguments := []mentalese.Term{}
		predicate := relationTemplate.Predicate

		for _, argument := range relationTemplate.Arguments {

			relationArgument := argument

			if argument.IsRelationSet() {

				relationArgument.TermValueRelationSet = builder.ReplaceTerm(relationArgument.TermValueRelationSet, from, to)

			} else {

				if argument.TermType == from.TermType && argument.TermValue == from.TermValue {

					relationArgument.TermType = to.TermType
					relationArgument.TermValue = to.TermValue
				}
			}

			arguments = append(arguments, relationArgument)
		}

		relation := mentalese.NewRelation(true, predicate, arguments)
		relations = append(relations, relation)
	}

	return relations
}
