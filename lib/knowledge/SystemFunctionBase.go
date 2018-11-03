package knowledge

import (
	"nli-go/lib/mentalese"
	"strings"
)

type SystemFunctionBase struct {
	KnowledgeBaseCore
}

func NewSystemFunctionBase() *SystemFunctionBase {
	return &SystemFunctionBase{}
}

func (base *SystemFunctionBase) GetMatchingGroups(set mentalese.RelationSet, knowledgeBaseIndex int) []RelationGroup {

	matchingGroups := []RelationGroup{}
	predicates := []string{"join", "split"}

	for _, setRelation := range set {
		for _, predicate:= range predicates {
			if predicate == setRelation.Predicate {
				// TODO calculate real cost
				matchingGroups = append(matchingGroups, RelationGroup{mentalese.RelationSet{setRelation}, knowledgeBaseIndex, worst_cost})
				break
			}
		}
	}

	return matchingGroups
}

func (base *SystemFunctionBase) Execute(input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	newBinding := binding
	found := false

	if input.Predicate == "split" {

		sourceValue, sourceFound := binding[input.Arguments[0].TermValue]
		if sourceFound {

			parts := strings.Split(sourceValue.TermValue, input.Arguments[1].TermValue)
			newBinding = binding.Copy()

			for i, argument := range input.Arguments[2:] {
				variable := argument.TermValue
				newTerm := mentalese.Term{}
				newTerm.TermType = mentalese.Term_stringConstant
				newTerm.TermValue = parts[i]
				newBinding[variable] = newTerm
			}

			found = true
		}
	}

	if input.Predicate == "join" {

		sep := ""
		result := mentalese.Term{}

		for _, argument := range input.Arguments[2:] {
			variable := argument.TermValue
			variableValue, variableFound := binding[variable]
			if variableFound {
				result.TermType = mentalese.Term_stringConstant
				result.TermValue += sep + variableValue.TermValue
				sep = input.Arguments[1].TermValue
			}
		}

		newBinding = binding.Copy()
		newBinding[input.Arguments[0].TermValue] = result

		found = true

	}

	return newBinding, found
}
