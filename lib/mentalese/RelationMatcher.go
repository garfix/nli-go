package mentalese

import (
	"nli-go/lib/common"
)

// This class matches relations to other relations and reports their bindings
// These concepts are used:
//
// sequence: a set of relations that is matched as a whole and shares a single Binding
// set: a set of unordered relations
//
// needle: the active subject, whose variables are to be bound
// haystack: the base of relations that serve as matching candidates

type RelationMatcher struct {
	functionBases []FunctionBase
	log           *common.SystemLog
}

func NewRelationMatcher(log *common.SystemLog) *RelationMatcher {
	return &RelationMatcher{log: log}
}

func (matcher *RelationMatcher) AddFunctionBase(functionBase FunctionBase) {
	matcher.functionBases = append(matcher.functionBases, functionBase)
}

type solutionNode struct {
	Binding Binding
	Indexes []int
}

func (matcher *RelationMatcher) MatchSequenceToSet(needleSequence RelationSet, haystackSet RelationSet, binding Binding) ([]Binding, bool) {

	bindings, _, _, match := matcher.MatchSequenceToSetWithIndexes(needleSequence, haystackSet, binding)
	return bindings, match
}

// Matches a relation sequence to a set
// Returns multiple bindings for variables in needleSequence
func (matcher *RelationMatcher) MatchSequenceToSetWithIndexes(needleSequence RelationSet, haystackSet RelationSet, binding Binding) ([]Binding, []int, []solutionNode, bool) {

	matcher.log.StartDebug("MatchSequenceToSetWithIndexes", needleSequence, haystackSet, binding)

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
			functionBinding, functionFound := matcher.MatchRelationToFunction(needleRelation, node.Binding)
			if functionFound {
				newIndexes := node.Indexes
				newNodes = append(newNodes, solutionNode{functionBinding, newIndexes})
			}

			someBindings, someIndexes := matcher.MatchRelationToSet(needleRelation, haystackSet, node.Binding)
			for i, someBinding := range someBindings {
				someIndex := someIndexes[i]
				newIndexes := append(node.Indexes, someIndex)
				newNodes = append(newNodes, solutionNode{someBinding, newIndexes})
			}
		}

		nodes = newNodes
	}

	for _, node := range nodes {
		newBindings = append(newBindings, node.Binding)
		matchedIndexes = append(matchedIndexes, node.Indexes...)
	}

	matchedIndexes = common.IntArrayDeduplicate(matchedIndexes)
	match = len(needleSequence) == 0 || len(matchedIndexes) > 0

	matcher.log.EndDebug("MatchSequenceToSetWithIndexes", newBindings, matchedIndexes, match)

	return newBindings, matchedIndexes, nodes, match
}

// functions like join(N, ' ', F, I, L)
// returns a binding with only one variable
func (matcher *RelationMatcher) MatchRelationToFunction(needleRelation Relation, binding Binding) (Binding, bool) {

	newBinding := Binding{}
	functionFound := false
	returnValue := Term{}

	for _, functionBase := range matcher.functionBases {
		returnValue, functionFound = functionBase.Execute(needleRelation, binding)
		if functionFound {
			newBinding = binding.Copy()
			newBinding[needleRelation.Arguments[0].TermValue] = returnValue
			break
		}
	}

	return newBinding, functionFound
}

// Matches a single relation to a relation set
// Returns multiple bindings
func (matcher *RelationMatcher) MatchRelationToSet(needleRelation Relation, haystackSet RelationSet, binding Binding) ([]Binding, []int) {

	matcher.log.StartDebug("matchRelationToSet", needleRelation, haystackSet, binding)

	newBindings := []Binding{}
	indexes := []int{}

	for i, haystackRelation := range haystackSet {

		newBinding, match := matcher.MatchTwoRelations(needleRelation, haystackRelation, binding)

		if match {
			newBindings = append(newBindings, newBinding)
			indexes = append(indexes, i)
		}
	}

	matcher.log.EndDebug("matchRelationToSet", newBindings, indexes)

	return newBindings, indexes
}

// Matches needleRelation to haystackRelation, using Binding
func (matcher *RelationMatcher) MatchTwoRelations(needleRelation Relation, haystackRelation Relation, binding Binding) (Binding, bool) {

	newBinding := binding.Copy()
	match := true

	matcher.log.StartDebug("MatchTwoRelations", needleRelation, haystackRelation, binding)

	// predicate
	if needleRelation.Predicate != haystackRelation.Predicate {
		match = false
	} else if len(needleRelation.Arguments) != len(haystackRelation.Arguments) {
		match = false
	} else {

		// arguments
		for i, subjectArgument := range needleRelation.Arguments {
			newBinding, match = matcher.BindTerm(subjectArgument, haystackRelation.Arguments[i], newBinding)

			if !match {
				break
			}
		}
	}

	matcher.log.EndDebug("MatchTwoRelations", newBinding, match)

	return newBinding, match
}
