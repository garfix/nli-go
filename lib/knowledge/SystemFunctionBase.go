package knowledge

import (
	"nli-go/lib/mentalese"
)

type SystemFunctionBase struct {
}

func NewSystemFunctionBase() *SystemFunctionBase {
	return &SystemFunctionBase{}
}

func (base *SystemFunctionBase) Execute(input mentalese.Relation, binding mentalese.Binding) (mentalese.Term, bool) {

	result := mentalese.Term{ TermType: mentalese.Term_stringConstant }
	found := true

	if input.Predicate == "join" {

		sep := ""

		for _, argument := range input.Arguments[2:] {
			variable := argument.TermValue
			variableValue, variableFound := binding[variable]
			if variableFound {
				result.TermValue += sep + variableValue.TermValue
				sep = input.Arguments[1].TermValue
			}
		}

	} else {
		found = false
	}

	return result, found
}
