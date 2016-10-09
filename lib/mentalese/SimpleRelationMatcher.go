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
func (matcher *SimpleRelationMatcher) MatchSequenceToSet(needleSequence SimpleRelationSet, haystackSet SimpleRelationSet) ([]int, SimpleBinding){

	matchedIndexes := []int{}
	binding := SimpleBinding{}

	common.Logf("MatchSequenceToSet: %v / %v\n", needleSequence, haystackSet)

	for _, haystackRelation := range haystackSet {

		index, newBoundVariables, found := matcher.matchSequenceToRelation(needleSequence, haystackRelation, binding)
		if found {

			binding = newBoundVariables
			matchedIndexes = append(matchedIndexes, index)

		}
	}

	common.Logf("MatchSequenceToSet end: %v / %v\n", matchedIndexes, binding)

	return matchedIndexes, binding
}

// Attempts to match a single pattern relation to a single relation
func (matcher *SimpleRelationMatcher) matchSequenceToRelation(subjectRelations SimpleRelationSet, patternRelation SimpleRelation, boundVariables SimpleBinding) (int, SimpleBinding, bool) {

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
