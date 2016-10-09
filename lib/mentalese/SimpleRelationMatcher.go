package mentalese

import "nli-go/lib/common"

// This class matches relations to other relations and reports their bindings
// These concepts are used:
//
// sequence: a set of relations that is matched as a whole and shares a single binding
// set: a set of unordered relations
//
// needle: the active subject, whose variables are to be bound
// haystack: the base of relations that serve as matching candidates

type SimpleRelationMatcher struct {

}

func NewSimpleRelationMatcher() *SimpleRelationMatcher {
	return &SimpleRelationMatcher{}
}

// matches a sequence to a set
func (matcher *SimpleRelationMatcher) Match(needleSequence SimpleRelationSet, haystackSet SimpleRelationSet) bool {
	matchedIndexes, _ := matcher.MatchSequenceToSet(needleSequence, haystackSet)
	return len(matchedIndexes) > 0
}

// matches a sequence to a set
func (matcher *SimpleRelationMatcher) MatchSequenceToSet(sequenceNeedle SimpleRelationSet, haystackSet SimpleRelationSet) ([]int, SimpleBinding){

	matchedIndexes := []int{}
	binding := SimpleBinding{}

	common.Logf("MatchSequenceToSet: %v / %v\n", sequenceNeedle, haystackSet)

	for _, patternRelation := range haystackSet {

		index, newBoundVariables, found := matcher.matchSubjectsToPattern(sequenceNeedle, patternRelation, binding)
		if found {

			binding = newBoundVariables
			matchedIndexes = append(matchedIndexes, index)

		}
	}

	common.Logf("MatchSequenceToSet end: %v / %v\n", matchedIndexes, binding)

	return matchedIndexes, binding
}

// Attempts to match a single pattern relation to a series of relations
func (matcher *SimpleRelationMatcher) matchSubjectsToPattern(subjectRelations SimpleRelationSet, patternRelation SimpleRelation, boundVariables SimpleBinding) (int, SimpleBinding, bool) {

	common.Logf("matchSubjectsToPattern: %v / %v\n", subjectRelations, patternRelation)

	for index, subjectRelation := range subjectRelations {

		newBoundVariables, matched := matcher.MatchSubjectToPattern(subjectRelation, patternRelation, boundVariables)

		if matched {

			common.Logf("matchSubjectsToPattern end: %d %v\n", index, newBoundVariables)

			return index, newBoundVariables, true
		}
	}

	common.Log("matchSubjectsToPattern end: failed\n")

	return 0, SimpleBinding{}, false
}

func (matcher *SimpleRelationMatcher) MatchSubjectToPattern(subjectRelation SimpleRelation, patternRelation SimpleRelation, binding SimpleBinding) (SimpleBinding, bool) {

	success := true

	common.Logf("MatchSubjectToPattern: %v / %v\n", subjectRelation, patternRelation)

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

	common.Logf("MatchSubjectToPattern: %v / %v\n", binding, success)

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

func (matcher *SimpleRelationMatcher) bindMultipleRelationsSingleBinding(relations SimpleRelationSet, binding SimpleBinding) SimpleRelationSet {

	for i, relation:= range relations {
		relations[i] = matcher.BindSingleRelationSingleBinding(relation, binding)
	}

	return relations
}

func (matcher *SimpleRelationMatcher) BindMultipleRelationsMultipleBindings(relations SimpleRelationSet, bindings []SimpleBinding) []SimpleRelationSet {

	relationSets := []SimpleRelationSet{}

	for _, binding := range bindings {
		relationSets = append(relationSets, matcher.bindMultipleRelationsSingleBinding(relations, binding))
	}

	return relationSets
}