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

type RelationMatcher struct {

}

func NewRelationMatcher() *RelationMatcher {
	return &RelationMatcher{}
}

// matches a sequence to a set
func (matcher *RelationMatcher) Match(needleSequence RelationSet, haystackSet RelationSet) bool {
	_, _, match := matcher.MatchSequenceToSet(needleSequence, haystackSet, Binding{})
	return match
}

// matches a sequence to a set
// NB!! should return multiple bindings
func (matcher *RelationMatcher) MatchSequenceToSet(needleSequence RelationSet, haystackSet RelationSet, binding Binding) ([]int, Binding, bool){

	matchedIndexes := []int{}
	match := true

	common.LogTree("MatchSequenceToSet", needleSequence, haystackSet, binding)

	newBinding := Binding{}.Merge(binding)

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

	common.LogTree("MatchSequenceToSet", matchedIndexes, newBinding, match)

	return matchedIndexes, newBinding, match
}


// Attempts to match a single pattern relation to a single relation
func (matcher *RelationMatcher) matchRelationToSet(needleRelation Relation, haystackSet RelationSet, binding Binding) (int, Binding, bool) {

	common.LogTree("matchRelationToSet", needleRelation, haystackSet, binding)

	aBinding := Binding{}.Merge(binding)
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

func (matcher *RelationMatcher) MatchTwoRelations(needleRelation Relation, haystackRelation Relation, binding Binding) (Binding, bool) {

	newBinding := Binding{}.Merge(binding)
	success := true

	common.LogTree("MatchTwoRelations", needleRelation, haystackRelation, binding)

	// predicate
	if needleRelation.Predicate != haystackRelation.Predicate {
		success = false
	} else {

		// arguments
		for i, subjectArgument := range needleRelation.Arguments {
			aBinding, ok := matcher.bindArgument(subjectArgument, haystackRelation.Arguments[i], newBinding)

			if ok {
				newBinding = aBinding
			} else {
				success = false
				break;
			}
		}
	}

	common.LogTree("MatchTwoRelations", newBinding, success)

	return newBinding, success
}
