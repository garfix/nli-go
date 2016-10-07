package mentalese

type SimpleRelationMatcher struct {

}

func NewSimpleRelationMatcher() *SimpleRelationMatcher {
	return &SimpleRelationMatcher{}
}

func (matcher *SimpleRelationMatcher) Match(subject *SimpleRelationSet, pattern *SimpleRelationSet) bool {
	matchedIndexes, _ := matcher.MatchSubjectsToPatterns(subject.Relations, pattern.Relations, false)
	return len(matchedIndexes) > 0
}

func (matcher *SimpleRelationMatcher) MatchSubjectsToPatterns(subjectRelations []SimpleRelation, patternRelations []SimpleRelation, allowPartial bool) ([]int, SimpleBinding){

	matchedIndexes := []int{}
	boundVariables := SimpleBinding{}

	for _, patternRelation := range patternRelations {

		index, newBoundVariables, found := matcher.matchSubjectsToPattern(subjectRelations, patternRelation, boundVariables)
		if found {

			boundVariables = newBoundVariables
			matchedIndexes = append(matchedIndexes, index)

		} else {

			if !allowPartial {
				return []int{}, SimpleBinding{}
			}
		}
	}

	return matchedIndexes, boundVariables
}

// Attempts to match a single pattern relation to a series of relations
func (matcher *SimpleRelationMatcher) matchSubjectsToPattern(subjectRelations []SimpleRelation, patternRelation SimpleRelation, boundVariables SimpleBinding) (int, SimpleBinding, bool) {

	for index, subjectRelation := range subjectRelations {

		newBoundVariables, matched := matcher.MatchSubjectToPattern(subjectRelation, patternRelation, boundVariables)

		if matched {
			return index, newBoundVariables, true
		}
	}

	return 0, SimpleBinding{}, false
}

func (matcher *SimpleRelationMatcher) MatchSubjectToPattern(subjectRelation SimpleRelation, patternRelation SimpleRelation, binding SimpleBinding) (SimpleBinding, bool) {

	success := true

	// predicate
	if subjectRelation.Predicate != patternRelation.Predicate {
		success = false
	} else {

		// arguments
		for i, subjectArgument := range subjectRelation.Arguments {
			newBinding, ok := matcher.bindArgument(subjectArgument, patternRelation.Arguments[i], binding)

			if ok {
				binding = newBinding
			} else {
				success = false
				break;
			}
		}
	}

	return binding, success
}

// Extends the binding with new variable bindings for the variables of subjectArgument
func (matcher *SimpleRelationMatcher) bindArgument(subjectArgument SimpleTerm, patternArgument SimpleTerm, binding SimpleBinding) (SimpleBinding, bool) {

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

func (matcher *SimpleRelationMatcher) BindSingleRelationSingleBinding(relation SimpleRelation, binding SimpleBinding) SimpleRelation {

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

func (matcher *SimpleRelationMatcher) bindMultipleRelationsSingleBinding(relations []SimpleRelation, binding SimpleBinding) []SimpleRelation {

	for i, relation:= range relations {
		relations[i] = matcher.BindSingleRelationSingleBinding(relation, binding)
	}

	return relations
}

func (matcher *SimpleRelationMatcher) BindMultipleRelationsMultipleBindings(relations []SimpleRelation, bindings []SimpleBinding) [][]SimpleRelation {

	relationSets := [][]SimpleRelation{}

	for _, binding := range bindings {
		relationSets = append(relationSets, matcher.bindMultipleRelationsSingleBinding(relations, binding))
	}

	return relationSets
}