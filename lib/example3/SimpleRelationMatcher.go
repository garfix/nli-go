package example3

type simpleRelationMatcher struct {

}

func NewSimpleRelationMatcher() *simpleRelationMatcher {
	return &simpleRelationMatcher{}
}

func (matcher *simpleRelationMatcher) Match(subject *SimpleRelationSet, pattern *SimpleRelationSet) bool {
	matchedIndexes, _ := matcher.matchSubjectsToPatterns(subject.relations, pattern.relations)
	return len(matchedIndexes) > 0
}

func (matcher *simpleRelationMatcher) matchSubjectsToPatterns(subjectRelations []SimpleRelation, patternRelations []SimpleRelation) ([]int, SimpleBinding){

	matchedIndexes := []int{}
	boundVariables := SimpleBinding{}

	for _, patternRelation := range patternRelations {

		index, newBoundVariables, found := matcher.matchSubjectsToPattern(subjectRelations, patternRelation, boundVariables)
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
func (matcher *simpleRelationMatcher) matchSubjectsToPattern(subjectRelations []SimpleRelation, patternRelation SimpleRelation, boundVariables SimpleBinding) (int, SimpleBinding, bool) {

	for index, subjectRelation := range subjectRelations {

		newBoundVariables, matched := matcher.matchSubjectToPattern(subjectRelation, patternRelation, boundVariables)

		if matched {
			return index, newBoundVariables, true
		}
	}

	return 0, SimpleBinding{}, false
}

func (matcher *simpleRelationMatcher) matchSubjectToPattern(subjectRelation SimpleRelation, patternRelation SimpleRelation, boundVariables SimpleBinding) (SimpleBinding, bool) {

	success := true

	// predicate
	if subjectRelation.Predicate != patternRelation.Predicate {
		success = false
	} else {

		// arguments
		for i, subjectArgument := range subjectRelation.Arguments {
			newBoundVariables, ok := matcher.bindArgument(subjectArgument, patternRelation.Arguments[i], boundVariables)

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

// Extends the binding with new variable bindings for the variables of subjectArgument
func (matcher *simpleRelationMatcher) bindArgument(subjectArgument SimpleTerm, patternArgument SimpleTerm, binding SimpleBinding) (SimpleBinding, bool) {

	success := false

	if subjectArgument.IsAnonymousVariable() || patternArgument.IsAnonymousVariable() {

		// anonymous variables always match, but do not bind

		success = true

	} else if subjectArgument.IsVariable() {

		// variable

		value := SimpleTerm{}

		// does patternRelationArgument occur in boundVariables?
		value, match := binding[subjectArgument.String()]
		if match {
			// it does, use the bound variable
			if patternArgument.Equals(value) {
				success = true
			}
		} else {
			// it does not, just assign the actual argument
			binding[subjectArgument.String()] = patternArgument
			success = true
		}

	} else {

		// subject is atom, constant

		if patternArgument.IsVariable() {
			// note: no binding is made
			success = true
		} else if patternArgument.Equals(subjectArgument) {
			success = true
		}
	}

	return binding, success
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