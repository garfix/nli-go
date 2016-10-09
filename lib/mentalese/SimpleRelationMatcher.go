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

	for _, needleRelation := range needleSequence {

		//index, newBoundVariables, found := matcher.matchSequenceToRelation(needleSequence, haystackRelation, binding)
		index, newBinding, found := matcher.matchRelationToSet(needleRelation, haystackSet, binding)

		if found {
			binding = newBinding
			matchedIndexes = append(matchedIndexes, index)
		}
	}

	common.Logf("MatchSequenceToSet end: %v / %v\n", matchedIndexes, binding)

	return matchedIndexes, binding
}


// Attempts to match a single pattern relation to a single relation
func (matcher *SimpleRelationMatcher) matchRelationToSet(needleRelation SimpleRelation, haystackSet SimpleRelationSet, binding SimpleBinding) (int, SimpleBinding, bool) {

	common.Logf("matchRelationToSet: %v / %v\n", needleRelation, haystackSet)

	for index, haystackRelation := range haystackSet {

		newBinding, matched := matcher.MatchNeedleToHaystack(needleRelation, haystackRelation, binding)

		if matched {

			common.Logf("matchRelationToSet end: %d %v\n", index, newBinding)

			return index, newBinding, true
		}
	}

	common.Log("matchRelationToSet end: failed\n")

	return 0, SimpleBinding{}, false
}

func (matcher *SimpleRelationMatcher) MatchNeedleToHaystack(needleRelation SimpleRelation, haystackRelation SimpleRelation, binding SimpleBinding) (SimpleBinding, bool) {

	success := true

	common.Logf("MatchSubjectToPattern: %v / %v\n", needleRelation, haystackRelation)

	// predicate
	if needleRelation.Predicate != haystackRelation.Predicate {
		success = false
	} else {

		// arguments
		for i, subjectArgument := range needleRelation.Arguments {
			newBinding, ok := matcher.bindArgument(subjectArgument, haystackRelation.Arguments[i], binding)

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
