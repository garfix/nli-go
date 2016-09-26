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

func (matcher *simpleRelationMatcher) matchRelations(relations []SimpleRelation, pattern []SimpleRelation) ([]int, SimpleBinding){

	matchedIndexes := []int{}
	boundVariables := SimpleBinding{}

	for _, patternRelation := range pattern {

		index, newBoundVariables, found := matcher.matchSingleRelation(relations, patternRelation, boundVariables)
		if found {

			boundVariables = newBoundVariables
			matchedIndexes = append(matchedIndexes, index)

		} else {
			return []int{}, SimpleBinding{}
		}
	}

	return matchedIndexes, boundVariables
}

// Attempts to match a single pattern relation to a series of relations
func (matcher *simpleRelationMatcher) matchSingleRelation(relations []SimpleRelation, patternRelation SimpleRelation, boundVariables SimpleBinding) (int, SimpleBinding, bool) {

	for index, relation := range relations {

		newBoundVariables, matched := matcher.matchRelationToRelation(relation, patternRelation, boundVariables)

		if matched {
			return index, newBoundVariables, true
		}
	}

	return 0, SimpleBinding{}, false
}

func (matcher *simpleRelationMatcher) matchRelationToRelation(relation SimpleRelation, patternRelation SimpleRelation, boundVariables SimpleBinding) (SimpleBinding, bool) {

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

func (matcher *simpleRelationMatcher) bindArgument(argument SimpleTerm, patternRelationArgument SimpleTerm, boundVariables SimpleBinding) (SimpleBinding, bool) {

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

func (matcher *simpleRelationMatcher) bindSingleRelationSingleBinding(relation SimpleRelation, binding SimpleBinding) SimpleRelation {

	for i, argument := range relation.Arguments {

		if argument.IsVariable() {
			newValue, found := binding[argument.TermValue]
			if found {
				relation.Arguments[i] = newValue
			}
		}
	}

	return relation
}

func (matcher *simpleRelationMatcher) bindMultipleRelationsSingleBinding(relations []SimpleRelation, binding SimpleBinding) []SimpleRelation {

	for i, relation:= range relations {
		relations[i] = matcher.bindSingleRelationSingleBinding(relation, binding)
	}

	return relations
}

func (matcher *simpleRelationMatcher) bindMultipleRelationsMultipleBindings(relations []SimpleRelation, bindings []SimpleBinding) [][]SimpleRelation {

	relationSets := [][]SimpleRelation{}

	for _, binding := range bindings {
		relationSets = append(relationSets, matcher.bindMultipleRelationsSingleBinding(relations, binding))
	}

	return relationSets
}