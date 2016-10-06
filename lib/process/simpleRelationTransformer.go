package process

import (
	"fmt"
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
)

type simpleRelationTransformer struct {
	transformations []mentalese.SimpleRelationTransformation
	matcher simpleRelationMatcher
}

// using transformations
func NewSimpleRelationTransformer(transformations[]mentalese.SimpleRelationTransformation) *simpleRelationTransformer {
	return &simpleRelationTransformer{transformations: transformations, matcher: simpleRelationMatcher{}}
}

// using rules (subset of transformations)
func NewSimpleRelationTransformer2(rules[]mentalese.SimpleRule) *simpleRelationTransformer {

	transformations := []mentalese.SimpleRelationTransformation{}

	for _, rule := range rules {

		transformation := mentalese.SimpleRelationTransformation{Replacement: []mentalese.SimpleRelation{ rule.Goal }, Pattern: rule.Pattern }
		transformations = append(transformations, transformation)
	}

	return &simpleRelationTransformer{transformations: transformations, matcher: simpleRelationMatcher{}}
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
func (transformer *simpleRelationTransformer) Extract(relationSet []mentalese.SimpleRelation) ([][]mentalese.SimpleRelation, []mentalese.SimpleBinding) {

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
func (transformer *simpleRelationTransformer) matchAllTransformations(relations []mentalese.SimpleRelation) ([][]int, [][]mentalese.SimpleRelation, []mentalese.SimpleBinding){

	matchedIndexes := [][]int{}
	replacements := [][]mentalese.SimpleRelation{}
	bindings := []mentalese.SimpleBinding{}

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
func (transformer *simpleRelationTransformer) matchSingleTransformation(relations []mentalese.SimpleRelation, transformation mentalese.SimpleRelationTransformation) ([]int, []mentalese.SimpleRelation, mentalese.SimpleBinding){

	fmt.Printf("Matching: %v / %v\n", transformation.Pattern, relations)

	matchedIndexes, oldBinding := transformer.matcher.matchSubjectsToPatterns(transformation.Pattern, relations, true)
	_, binding := transformer.matcher.matchSubjectsToPatterns(relations, transformation.Pattern, true)

fmt.Printf("Matched: %v %v %v\n", matchedIndexes, oldBinding, binding)

	replacements := []mentalese.SimpleRelation{}
	if len(matchedIndexes) > 0 {
		replacements = transformer.createReplacements(transformation.Replacement, oldBinding)
	}

	return matchedIndexes, replacements, binding
}

func (transformer *simpleRelationTransformer) createReplacements(relations []mentalese.SimpleRelation, bindings mentalese.SimpleBinding) []mentalese.SimpleRelation {

	replacements := []mentalese.SimpleRelation{}

	fmt.Printf("Replace! %v %v\n", relations, bindings)

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