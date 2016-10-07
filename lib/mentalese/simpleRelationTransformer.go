package mentalese

import (
	"nli-go/lib/common"
)

type simpleRelationTransformer struct {
	transformations []SimpleRelationTransformation
	matcher SimpleRelationMatcher
}

// using transformations
func NewSimpleRelationTransformer(transformations[]SimpleRelationTransformation) *simpleRelationTransformer {
	return &simpleRelationTransformer{transformations: transformations, matcher: SimpleRelationMatcher{}}
}

// using rules (subset of transformations)
func NewSimpleRelationTransformer2(rules[]SimpleRule) *simpleRelationTransformer {

	transformations := []SimpleRelationTransformation{}

	for _, rule := range rules {

		transformation := SimpleRelationTransformation{Replacement: SimpleRelationSet{ rule.Goal }, Pattern: rule.Pattern }
		transformations = append(transformations, transformation)
	}

	return &simpleRelationTransformer{transformations: transformations, matcher: SimpleRelationMatcher{}}
}

//// return the original relations, but replace the ones that have matched with their replacements
//func (transformer *simpleRelationTransformer) Replace(relationSet *SimpleRelationSet) *SimpleRelationSet {
//
//	matchedIndexes, replacements := transformer.matchAllTransformations(relationSet.relations)
//	newRelations := NewSimpleRelationSet()
//
//	for i, oldRelation := range relationSet.GetRelations()  {
//		if !intArrayContains(matchedIndexes, i) {
//			newRelations.AddRelation(oldRelation)
//		}
//	}
//
//	newRelations.AddRelations(replacements)
//
//	return newRelations
//}

// return only the replacements
func (transformer *simpleRelationTransformer) Extract(relationSet SimpleRelationSet) ([]SimpleRelationSet, []SimpleBinding) {

	_, replacements, bindings := transformer.matchAllTransformations(relationSet)
	return replacements, bindings
}

//// only add the replacements to the original relations
//func (transformer *simpleRelationTransformer) Append(relationSet *SimpleRelationSet) *SimpleRelationSet {
//
//	_, replacements := transformer.matchAllTransformations(relationSet.relations)
//
//	newRelations := NewSimpleRelationSet2(relationSet.GetRelations())
//	newRelations.AddRelations(replacements)
//
//	return newRelations
//}

// Attempts all transformations on all relations
// Returns the indexes of the matched relations, and the replacements that were created
func (transformer *simpleRelationTransformer) matchAllTransformations(relations SimpleRelationSet) ([][]int, []SimpleRelationSet, []SimpleBinding){

	matchedIndexes := [][]int{}
	replacements := []SimpleRelationSet{}
	bindings := []SimpleBinding{}

	for _, transformation := range transformer.transformations {

		newMatchedIndexes, newReplacements, newBinding := transformer.matchSingleTransformation(relations, transformation)
		if len(newMatchedIndexes) > 0 {
			matchedIndexes = append(matchedIndexes, common.IntArrayDeduplicate(newMatchedIndexes))
			replacements = append(replacements, newReplacements)
			bindings = append(bindings, newBinding)
		}
	}

	return matchedIndexes, replacements, bindings
}

// Attempts to match a single transformation
// Returns the indexes of matched relations, and the replacements
func (transformer *simpleRelationTransformer) matchSingleTransformation(relations SimpleRelationSet, transformation SimpleRelationTransformation) ([]int, SimpleRelationSet, SimpleBinding){

	common.Logf("Matching: %v / %v\n", transformation.Pattern, relations)

	matchedIndexes, oldBinding := transformer.matcher.MatchSubjectsToPatterns(transformation.Pattern, relations, true)
	_, binding := transformer.matcher.MatchSubjectsToPatterns(relations, transformation.Pattern, true)

	common.Logf("Matched: %v %v %v\n", matchedIndexes, oldBinding, binding)

	replacements := SimpleRelationSet{}
	if len(matchedIndexes) > 0 {
		replacements = transformer.createReplacements(transformation.Replacement, oldBinding)
	}

	return matchedIndexes, replacements, binding
}

func (transformer *simpleRelationTransformer) createReplacements(relations SimpleRelationSet, bindings SimpleBinding) SimpleRelationSet {

	replacements := SimpleRelationSet{}

	common.Logf("Replace! %v %v\n", relations, bindings)

	for _, relation := range relations {

		newRelation := relation

		for i, argument := range relation.Arguments {

			if argument.IsVariable() {
				value, found := bindings[argument.String()]
				if found {
					newRelation.Arguments[i] = value
				} else {
					// replacement could not be bound!
				}
			}
		}

		replacements = append(replacements, relation)
	}

	return replacements
}