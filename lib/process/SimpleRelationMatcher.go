package process

import "nli-go/lib/mentalese"

type simpleRelationMatcher struct {

}

func NewSimpleRelationMatcher() *simpleRelationMatcher {
	return &simpleRelationMatcher{}
}

func (matcher *simpleRelationMatcher) Match(subject *mentalese.SimpleRelationSet, pattern *mentalese.SimpleRelationSet) bool {
	matchedIndexes, _ := matcher.matchSubjectsToPatterns(subject.Relations, pattern.Relations, false)
	return len(matchedIndexes) > 0
}

func (matcher *simpleRelationMatcher) matchSubjectsToPatterns(subjectRelations []mentalese.SimpleRelation, patternRelations []mentalese.SimpleRelation, allowPartial bool) ([]int, mentalese.SimpleBinding){

	matchedIndexes := []int{}
	boundVariables := mentalese.SimpleBinding{}

	for _, patternRelation := range patternRelations {

		index, newBoundVariables, found := matcher.matchSubjectsToPattern(subjectRelations, patternRelation, boundVariables)
		if found {

			boundVariables = newBoundVariables
			matchedIndexes = append(matchedIndexes, index)

		} else {

			if !allowPartial {
				return []int{}, mentalese.SimpleBinding{}
			}
		}
	}

	return matchedIndexes, boundVariables
}

// Attempts to match a single pattern relation to a series of relations
func (matcher *simpleRelationMatcher) matchSubjectsToPattern(subjectRelations []mentalese.SimpleRelation, patternRelation mentalese.SimpleRelation, boundVariables mentalese.SimpleBinding) (int, mentalese.SimpleBinding, bool) {

	for index, subjectRelation := range subjectRelations {

		newBoundVariables, matched := matcher.MatchSubjectToPattern(subjectRelation, patternRelation, boundVariables)

		if matched {
			return index, newBoundVariables, true
		}
	}

	return 0, mentalese.SimpleBinding{}, false
}

func (matcher *simpleRelationMatcher) MatchSubjectToPattern(subjectRelation mentalese.SimpleRelation, patternRelation mentalese.SimpleRelation, binding mentalese.SimpleBinding) (mentalese.SimpleBinding, bool) {

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
func (matcher *simpleRelationMatcher) bindArgument(subjectArgument mentalese.SimpleTerm, patternArgument mentalese.SimpleTerm, binding mentalese.SimpleBinding) (mentalese.SimpleBinding, bool) {

	success := false

	if subjectArgument.IsAnonymousVariable() || patternArgument.IsAnonymousVariable() {

		// anonymous variables always match, but do not bind

		success = true

	} else if subjectArgument.IsVariable() {

		// variable

		value := mentalese.SimpleTerm{}

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

func (matcher *simpleRelationMatcher) BindSingleRelationSingleBinding(relation mentalese.SimpleRelation, binding mentalese.SimpleBinding) mentalese.SimpleRelation {

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

func (matcher *simpleRelationMatcher) bindMultipleRelationsSingleBinding(relations []mentalese.SimpleRelation, binding mentalese.SimpleBinding) []mentalese.SimpleRelation {

	for i, relation:= range relations {
		relations[i] = matcher.BindSingleRelationSingleBinding(relation, binding)
	}

	return relations
}

func (matcher *simpleRelationMatcher) bindMultipleRelationsMultipleBindings(relations []mentalese.SimpleRelation, bindings []mentalese.SimpleBinding) [][]mentalese.SimpleRelation {

	relationSets := [][]mentalese.SimpleRelation{}

	for _, binding := range bindings {
		relationSets = append(relationSets, matcher.bindMultipleRelationsSingleBinding(relations, binding))
	}

	return relationSets
}