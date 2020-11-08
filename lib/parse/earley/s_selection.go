package earley

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
)

// Create a new s-selection
// Inherit the types that were bound to the antecedent
// If not inherited, find a proper type from the sense
func combineSSelection(predicates *mentalese.Meta, parentTypes []string, rule parse.GrammarRule) (parse.SSelection, bool) {

	// start with the type of the antecedent
	sSelection := parse.SSelection{ parentTypes }

	// for each consequent
	for _, singleConsequentVariables := range rule.GetAllConsequentVariables() {

		// single consequent
		consequentType := []string{}

		// for each consequentVariable of a single consequent
		for _, consequentVariable := range singleConsequentVariables {

			singleArgumentType := ""

			// compare with each antecedent variable
			for a, antecedentVariable := range rule.GetAntecedentVariables() {
				if consequentVariable == antecedentVariable {
					singleArgumentType = parentTypes[a]
					break
				}
			}

			if singleArgumentType == "" {
				singleArgumentType = getTypeFromSense(predicates, consequentVariable, rule.Sense)
			}

			consequentType = append(consequentType, singleArgumentType)
		}
		sSelection = append(sSelection, consequentType)
	}

	return sSelection, true
}

func getTypeFromSense(meta *mentalese.Meta, variable string, sense mentalese.RelationSet) string {

	sType := ""

	for _, relation := range sense {
		for i, argument := range relation.Arguments {
			if argument.IsVariable() && argument.TermValue == variable {

				sType = meta.GetSort(relation.Predicate, i)
				if sType != "" {
					goto end
				}
			} else if argument.IsRelationSet() {
				sTypeRecursive := getTypeFromSense(meta, variable, argument.TermValueRelationSet)
				if sTypeRecursive != "" {
					sType = sTypeRecursive
					goto end
				}
			} else if argument.IsRule() {
				// no need to implement
			} else if argument.IsList() {
				// no need to implement
			}
		}
	}

end:

	return sType
}