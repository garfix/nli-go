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
	_, _, match := matcher.MatchSequenceToSet(needleSequence, haystackSet, SimpleBinding{})
	return match
}

// matches a sequence to a set
// NB!! should return multiple bindings
func (matcher *SimpleRelationMatcher) MatchSequenceToSet(needleSequence SimpleRelationSet, haystackSet SimpleRelationSet, binding SimpleBinding) ([]int, SimpleBinding, bool){

	matchedIndexes := []int{}
	match := true

	common.LogTree("MatchSequenceToSet", needleSequence, haystackSet)

	newBinding := SimpleBinding{}.Merge(binding)

	for _, needleRelation := range needleSequence {

		index, aBinding, found := matcher.matchRelationToSet(needleRelation, haystackSet, newBinding)

		if found {
			newBinding = aBinding
			matchedIndexes = append(matchedIndexes, index)
		} else {
			newBinding = binding
			matchedIndexes = []int{}
			match = false
			break
		}
	}

	common.LogTree("MatchSequenceToSet", matchedIndexes, binding, match)

	return matchedIndexes, newBinding, match
}


// Attempts to match a single pattern relation to a single relation
func (matcher *SimpleRelationMatcher) matchRelationToSet(needleRelation SimpleRelation, haystackSet SimpleRelationSet, binding SimpleBinding) (int, SimpleBinding, bool) {

	common.LogTree("matchRelationToSet", needleRelation, haystackSet, binding)

	aBinding := SimpleBinding{}.Merge(binding)
	newBinding := binding
	i := 0
	bound := false

	for index, haystackRelation := range haystackSet {

		aBinding, matched := matcher.MatchTwoRelations(needleRelation, haystackRelation, aBinding)

		if matched {

			i = index
			newBinding = aBinding
			bound = true
			break
		}
	}

	common.LogTree("matchRelationToSet", bound, i, newBinding)

	return i, newBinding, bound
}

func (matcher *SimpleRelationMatcher) MatchTwoRelations(needleRelation SimpleRelation, haystackRelation SimpleRelation, binding SimpleBinding) (SimpleBinding, bool) {

	success := true

	common.LogTree("MatchTwoRelations", needleRelation, haystackRelation, binding)

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

	common.LogTree("MatchTwoRelations", binding, success)

	return binding, success
}
