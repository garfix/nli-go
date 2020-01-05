package earley

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
)

// Create a new child sense by applying a rule that contains a child sense template, and inherit the parent sense
// Example:
// 		parentSSelection: person
// 		{ rule: relative_clause(E1) -> np(E2) aux_be(C1) holding(P1),           sense: hold(P1, E2, E1) }
// 		return E1 = person, E2 = object; or: false in case of conflict
func combineSSelection(predicates mentalese.Predicates, parentType string, rule parse.GrammarRule) (parse.SSelection, bool) {

	sSelection := parse.SSelection{parentType}
	antecedent := rule.GetAntecedentVariable()

	for _, variable := range rule.GetConsequentVariables() {

		sType := ""

		if variable == antecedent {
			sType = parentType
		} else {
			sType = getTypeFromSense(predicates, variable, rule.Sense)
		}

		sSelection = append(sSelection, sType)
	}

	return sSelection, true
}

func getTypeFromSense(predicates mentalese.Predicates, variable string, sense mentalese.RelationSet) string {

	sType := ""

	for _, relation := range sense {
		for i, argument := range relation.Arguments {
			if argument.IsVariable() && argument.TermValue == variable {

				sTypes, found := predicates[relation.Predicate]
				if found {
					sType = sTypes.EntityTypes[i]
					goto end
				}
			}
		}
	}

end:

	return sType
}