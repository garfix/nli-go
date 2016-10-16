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

// matches a relation sequence to a set
// NB!! should return multiple bindings
func (matcher *RelationMatcher) MatchSequenceToSet(needleSequence RelationSet, haystackSet RelationSet, binding Binding) (Binding, []int, bool){

	common.LogTree("MatchSequenceToSet", needleSequence, haystackSet, binding)

	newBinding := binding.Copy()
	matchedIndexes := []int{}
	match := true
	index := 0

	for _, needleRelation := range needleSequence {

		newBinding, index, match = matcher.MatchRelationToSet(needleRelation, haystackSet, newBinding)

		if match {
			matchedIndexes = append(matchedIndexes, index)
		} else {
			break
		}
	}

	common.LogTree("MatchSequenceToSet", newBinding, matchedIndexes, match)

	return newBinding, matchedIndexes, match
}

// Attempts to match a single pattern relation to a single relation
func (matcher *RelationMatcher) MatchRelationToSet(needleRelation Relation, haystackSet RelationSet, binding Binding) (Binding, int, bool) {

	common.LogTree("matchRelationToSet", needleRelation, haystackSet, binding)

	newBinding := binding.Copy()
	match := false
	index := 0

	for i, haystackRelation := range haystackSet {

		newBinding, match = matcher.MatchTwoRelations(needleRelation, haystackRelation, newBinding)

		if match {
			index = i
			break
		}
	}

	common.LogTree("matchRelationToSet", newBinding, index, match)

	return newBinding, index, match
}

// Matches needleRelation to haystackRelation, using binding
func (matcher *RelationMatcher) MatchTwoRelations(needleRelation Relation, haystackRelation Relation, binding Binding) (Binding, bool) {

	newBinding := binding.Copy()
	match := true

	common.LogTree("MatchTwoRelations", needleRelation, haystackRelation, binding)

	// predicate
	if needleRelation.Predicate != haystackRelation.Predicate {
		match = false
	} else {

		// arguments
		for i, subjectArgument := range needleRelation.Arguments {
			newBinding, match = matcher.BindTerm(subjectArgument, haystackRelation.Arguments[i], newBinding)

			if !match {
				break;
			}
		}
	}

	common.LogTree("MatchTwoRelations", newBinding, match)

	return newBinding, match
}
