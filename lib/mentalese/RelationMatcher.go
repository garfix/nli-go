package mentalese

import (
	"nli-go/lib/common"
)

// This class matches relations to other relations and reports their bindings
// These concepts are used:
//
// sequence: a set of relations that is matched as a whole and shares a single binding
// set: a set of unordered relations
//
// needle: the active subject, whose variables are to be bound
// haystack: the base of relations that serve as matching candidates

type RelationMatcher struct {
	functionBases []FunctionBase
}

func NewRelationMatcher() *RelationMatcher {
	return &RelationMatcher{}
}

func (matcher *RelationMatcher) AddFunctionBase(functionBase FunctionBase) {
	matcher.functionBases = append(matcher.functionBases, functionBase)
}

type solutionNode struct {
	binding Binding
	indexes []int
}

// Matches a relation sequence to a set
// Returns multiple bindings
func (matcher *RelationMatcher) MatchSequenceToSet(needleSequence RelationSet, haystackSet RelationSet, binding Binding) ([]Binding, []int, bool){

	common.LogTree("MatchSequenceToSet", needleSequence, haystackSet, binding)

	newBindings := []Binding{}
	matchedIndexes := []int{}
	match := true

	nodes := []solutionNode{
		{binding, []int{}},
	}

	for _, needleRelation := range needleSequence {

		newNodes := []solutionNode{}

		for _, node := range nodes {


			// functions like join(N, ' ', F, I, L)

			for _, functionBase := range matcher.functionBases {
				functionBinding := node.binding.Copy()
				returnValue, ok := functionBase.Execute(needleRelation, functionBinding)
				if ok {
					functionBinding[needleRelation.Arguments[0].TermValue] = returnValue
					newIndexes := append(node.indexes, 0)
					newNodes = append(newNodes, solutionNode{functionBinding, newIndexes})
				}
			}



			someBindings, someIndexes := matcher.MatchRelationToSet(needleRelation, haystackSet, node.binding)
			for i, someBinding := range someBindings {
				someIndex := someIndexes[i]
				newIndexes := append(node.indexes, someIndex)
				newNodes = append(newNodes, solutionNode{someBinding, newIndexes})
			}
		}

		nodes = newNodes
	}

	for _, node := range nodes {
		newBindings = append(newBindings, node.binding)
		matchedIndexes = append(matchedIndexes, node.indexes...)
	}

	matchedIndexes = common.IntArrayDeduplicate(matchedIndexes)
	match = len(needleSequence) == 0 || len(matchedIndexes) > 0

	common.LogTree("MatchSequenceToSet", newBindings, matchedIndexes, match)

	return newBindings, matchedIndexes, match
}

// Matches a single relation to a relation set
// Returns multiple bindings
func (matcher *RelationMatcher) MatchRelationToSet(needleRelation Relation, haystackSet RelationSet, binding Binding) ([]Binding, []int) {

	common.LogTree("matchRelationToSet", needleRelation, haystackSet, binding)

	newBindings := []Binding{}
	indexes := []int{}

	for i, haystackRelation := range haystackSet {

		newBinding, match := matcher.MatchTwoRelations(needleRelation, haystackRelation, binding)

		if match {
			newBindings = append(newBindings, newBinding)
			indexes = append(indexes, i)
		}
	}

	common.LogTree("matchRelationToSet", newBindings, indexes)

	return newBindings, indexes
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
