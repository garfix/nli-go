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

func (matcher *RelationMatcher) MatchSequenceToSet(needleSequence RelationSet, haystackSet RelationSet, binding Binding) (BindingSet, bool) {

	matcher.log.StartDebug("MatchSequenceToSetWithIndexes", needleSequence, haystackSet, binding)

	newBindings := NewBindingSet()

	match := true

	nodes := []solutionNode{
		{binding, []int{}},
	}

	for _, needleRelation := range needleSequence {

		var newNodes []solutionNode

		nodeMatches := false

		for _, node := range nodes {

			// functions like join(N, ' ', F, I, L)
			functionBinding, functionFound, success := matcher.ExecuteFunction(needleRelation, node.Binding)
			if functionFound  && success {
				newIndexes := node.Indexes
				newNodes = append(newNodes, solutionNode{functionBinding, newIndexes})
				nodeMatches = true
			}

			someBindings, someIndexes := matcher.MatchRelationToSet(needleRelation, haystackSet, node.Binding)
			for i, someBinding := range someBindings.GetAll() {
				someIndex := someIndexes[i]
				newIndexes := append(node.Indexes, someIndex)
				newNodes = append(newNodes, solutionNode{someBinding, newIndexes})
			}
			if !someBindings.IsEmpty() {
				nodeMatches = true
			}
		}

		if !nodeMatches {
			match = false
		}

		nodes = newNodes
	}

	for _, node := range nodes {
		newBindings.Add(node.Binding)
	}

	matcher.log.EndDebug("MatchSequenceToSetWithIndexes", newBindings, match)

	return newBindings, match
}


// functions like join(N, ' ', F, I, L)
// returns a binding with only one variable
func (matcher *RelationMatcher) ExecuteFunction(needleRelation Relation, binding Binding) (Binding, bool, bool) {

	newBinding := NewBinding()
	resultBinding := NewBinding()
	functionFound := false
	success := false

	for _, functionBase := range matcher.functionBases {
		resultBinding, functionFound, success = functionBase.Execute(needleRelation, binding)
		if functionFound {
			newBinding = resultBinding
			break
		}
	}

	return newBinding, functionFound, success
}

// Matches a single relation to a relation set
// Returns multiple bindings
func (matcher *RelationMatcher) MatchRelationToSet(needleRelation Relation, haystackSet RelationSet, binding Binding) (BindingSet, []int) {

	matcher.log.StartDebug("matchRelationToSet", needleRelation, haystackSet, binding)

	newBindings := NewBindingSet()
	indexes := []int{}

	for i, haystackRelation := range haystackSet {

		newBinding, match := matcher.MatchTwoRelations(needleRelation, haystackRelation, binding)

		if match {
			newBindings.Add(newBinding)
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
	} else if needleRelation.Positive != haystackRelation.Positive {
		match = false
	} else if len(needleRelation.Arguments) != len(haystackRelation.Arguments) {
		match = false
	} else {

		// arguments
		for i, subjectArgument := range needleRelation.Arguments {
			newBinding, match = matcher.MatchTerm(subjectArgument, haystackRelation.Arguments[i], newBinding)

			if !match {
				break
			}
		}
	}

	matcher.log.EndDebug("MatchTwoRelations", newBinding, match)

	return newBinding, match
}


// Extends the Binding with new variable bindings for the variables of subjectArgument
func (matcher *RelationMatcher) MatchTerm(subjectArgument Term, patternArgument Term, subjectBinding Binding) (Binding, bool) {

	success := false

	if subjectArgument.IsAnonymousVariable() || patternArgument.IsAnonymousVariable() {

		// anonymous variables always match, but do not bind

		// A, _
		// _, A
		return subjectBinding, true

	} else if subjectArgument.IsVariable() {

		value, match := subjectBinding.Get(subjectArgument.String())
		if match {

			if patternArgument.IsVariable() {
				// A, B {A:C}
				// A, B {A:13}
				success = true
			} else {
				// A, 13 {A:B}
				// A, 13 {A:15}
				success = patternArgument.Equals(value)
			}

			return subjectBinding, success

		} else {

			// A, 13, {B:7} => {B:7, A:13}
			newBinding := subjectBinding.Copy()
			newBinding.Set(subjectArgument.String(), patternArgument)
			return newBinding, true
		}

	} else if subjectArgument.IsRelationSet() {

		newBinding := subjectBinding.Copy()

		if patternArgument.IsVariable() {
			// [ isa(E, very) ], V
			success = true

		} else if patternArgument.IsRelationSet() {

			subSetBindings, ok := matcher.MatchSequenceToSet(subjectArgument.TermValueRelationSet, patternArgument.TermValueRelationSet, newBinding)

			if ok {
				newBinding = subSetBindings.Get(0)
				success = true
			}
		}

		return newBinding, success

	} else if subjectArgument.IsRule() {

		panic("to be implemented")

	} else if subjectArgument.IsList() {

		panic("to be implemented")

	} else {

		// subject is atom, constant

		if patternArgument.IsVariable() {
			// 13, V
			success = true
		} else if patternArgument.Equals(subjectArgument) {
			// 13, 13
			// female, female
			// 'Jack', 'Jack'
			success = true
		}

		return subjectBinding, success
	}
}
