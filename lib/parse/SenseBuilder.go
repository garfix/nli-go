package parse

import (
	"nli-go/lib/mentalese"
)

type SenseBuilder struct {
	variableGenerator *mentalese.VariableGenerator
}

func NewSenseBuilder(variableGenerator *mentalese.VariableGenerator) SenseBuilder {
	return SenseBuilder{
		variableGenerator: variableGenerator,
	}
}

// Creates a map of formal variables to actual variables (new variables are created)
func (builder SenseBuilder) CreateVariableMap(actualAntecedents []string, antecedentVariables []string, allConsequentVariables [][]string) map[string]mentalese.Term {

	m := map[string]mentalese.Term{}

	if len(actualAntecedents) == 0 {
		return m
	}

	for i, antecedentVariable := range antecedentVariables {
		m[antecedentVariable] = mentalese.NewTermVariable(actualAntecedents[i])
	}

	for _, consequentVariables := range allConsequentVariables {
		for _, consequentVariable := range consequentVariables {

			_, present := m[consequentVariable]
			if !present {
				m[consequentVariable] = mentalese.NewTermVariable(builder.variableGenerator.GenerateVariable(consequentVariable).TermValue)
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
					variableMap[variable] = mentalese.NewTermVariable(builder.variableGenerator.GenerateVariable(variable).TermValue)
				}
			} else if argument.IsRelationSet() {
				childMap := builder.ExtendVariableMap(argument.TermValueRelationSet, variableMap)
				for key, value := range childMap {
					variableMap[key] = value
				}
			} else if argument.IsRule() {
				childMap := builder.ExtendVariableMap(mentalese.RelationSet{ argument.TermValueRule.Goal }, variableMap)
				for key, value := range childMap {
					variableMap[key] = value
				}
				childMap = builder.ExtendVariableMap(argument.TermValueRule.Pattern, variableMap)
				for key, value := range childMap {
					variableMap[key] = value
				}
			} else if argument.IsList() {
				panic("to be implemented")
			}
		}
	}

	return variableMap
}

// Create actual relations given a set of templates and a variable map (formal to actual variables)
func (builder SenseBuilder) CreateGrammarRuleRelations(relationTemplates mentalese.RelationSet) mentalese.RelationSet {

	relations := mentalese.RelationSet{}

	for _, relation := range relationTemplates {

		newRelation := relation.Copy()

		for a, argument := range newRelation.Arguments {

			if argument.IsVariable() {

				newRelation.Arguments[a].TermType = argument.TermType
				newRelation.Arguments[a].TermValue = argument.TermValue

			} else if argument.IsRelationSet() {

				newRelation.Arguments[a].TermType = mentalese.TermTypeRelationSet
				newRelation.Arguments[a].TermValueRelationSet = builder.CreateGrammarRuleRelations(argument.TermValueRelationSet)

			} else if argument.IsRule() {

				newGoal := builder.CreateGrammarRuleRelations(mentalese.RelationSet{ argument.TermValueRule.Goal })
				newPattern := builder.CreateGrammarRuleRelations(argument.TermValueRule.Pattern)
				newRule := mentalese.Rule{ Goal: newGoal[0], Pattern: newPattern }
				newRelation.Arguments[a].TermType = mentalese.TermTypeRule
				newRelation.Arguments[a].TermValueRule = &newRule

			} else if argument.IsList() {
				panic("to be implemented")
			}
		}

		relations = append(relations, newRelation)
	}

	return relations
}

// Create actual relations given a set of templates and an actual variable to replace any * positions
func (builder SenseBuilder) CreateLexItemRelations(relationTemplates mentalese.RelationSet, variable string) mentalese.RelationSet {

	from := mentalese.Term{TermType: mentalese.TermTypeVariable, TermValue: "E"}
	to := mentalese.Term{TermType: mentalese.TermTypeVariable, TermValue: variable}

	return relationTemplates.ReplaceTerm(from, to)
}

