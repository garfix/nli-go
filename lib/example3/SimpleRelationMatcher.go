package example3

type simpleRelationMatcher struct {

}

func NewSimpleRelationMatcher() *simpleRelationMatcher {
	return &simpleRelationMatcher{}
}

func (matcher *simpleRelationMatcher) Match(pattern *SimpleRelationSet, subject *SimpleRelationSet) bool {
	matchedIndexes, _ := matcher.matchRelations(subject.relations, pattern.relations)
	return len(matchedIndexes) > 0
}

func (matcher *simpleRelationMatcher) matchRelations(relations []SimpleRelation, pattern []SimpleRelation) ([]int, map[string]SimpleTerm){

	matchedIndexes := []int{}
	boundVariables := map[string]SimpleTerm{}

	for _, patternRelation := range pattern {

		index, newBoundVariables, found := matcher.matchSingleRelation(relations, patternRelation, boundVariables)
		if found {

			boundVariables = newBoundVariables
			matchedIndexes = append(matchedIndexes, index)

		} else {
			return []int{}, map[string]SimpleTerm{}
		}
	}

	return matchedIndexes, boundVariables
}

// Attempts to match a single pattern relation to a series of relations
func (matcher *simpleRelationMatcher) matchSingleRelation(relations []SimpleRelation, patternRelation SimpleRelation, boundVariables map[string]SimpleTerm) (int, map[string]SimpleTerm, bool) {

	for index, relation := range relations {

		newBoundVariables, matched := matcher.matchRelationToRelation(relation, patternRelation, boundVariables)

		if matched {
			return index, newBoundVariables, true
		}
	}

	return 0, map[string]SimpleTerm{}, false
}

func (matcher *simpleRelationMatcher) matchRelationToRelation(relation SimpleRelation, patternRelation SimpleRelation, boundVariables map[string]SimpleTerm) (map[string]SimpleTerm, bool) {

	success := true

	// predicate
	if relation.Predicate != patternRelation.Predicate {
		success = false
	} else {

		// arguments
		for i, argument := range relation.Arguments {
			newBoundVariables, ok := matcher.bindArgument(argument, patternRelation.Arguments[i], boundVariables)

			if ok {
				boundVariables = newBoundVariables
			} else {
				success = false
				break;
			}
		}
	}

	return boundVariables, success
}

func (matcher *simpleRelationMatcher) bindArgument(argument SimpleTerm, patternRelationArgument SimpleTerm, boundVariables map[string]SimpleTerm) (map[string]SimpleTerm, bool) {

	success := false

	if patternRelationArgument.IsVariable() {

		// variable

		value := SimpleTerm{}

		// does patternRelationArgument occur in boundVariables?
		value, match := boundVariables[patternRelationArgument.AsKey()]
		if match {
			// it does, use the bound variable
			if argument.Equals(value) {
				success = true
			}
		} else {
			// it does not, just assign the actual argument
			boundVariables[patternRelationArgument.AsKey()] = argument
			success = true
		}

	} else {

		// atom, constant

		if argument.Equals(patternRelationArgument) {
			success = true
		}
	}

	return boundVariables, success
}