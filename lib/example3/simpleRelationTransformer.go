package example3

import "fmt"

type simpleRelationTransformer struct {
	transformations []SimpleRelationTransformation
	matcher simpleRelationMatcher
}

// using transformations
func NewSimpleRelationTransformer(transformations[]SimpleRelationTransformation) *simpleRelationTransformer {
	return &simpleRelationTransformer{transformations: transformations, matcher: simpleRelationMatcher{}}
}

// using rules (subset of transformations)
func NewSimpleRelationTransformer2(rules[]SimpleRule) *simpleRelationTransformer {

	transformations := []SimpleRelationTransformation{}

	for _, rule := range rules {

		transformation := SimpleRelationTransformation{Replacement: []SimpleRelation{ rule.Goal }, Pattern: rule.Pattern }
		transformations = append(transformations, transformation)
	}

	return &simpleRelationTransformer{transformations: transformations, matcher: simpleRelationMatcher{}}
}

// return the original relations, but replace the ones that have matched with their replacements
func (transformer *simpleRelationTransformer) Replace(relationSet *SimpleRelationSet) *SimpleRelationSet {

	matchedIndexes, replacements := transformer.matchAllTransformations(relationSet.relations)
	newRelations := NewSimpleRelationSet()

	for i, oldRelation := range relationSet.GetRelations()  {
		if !intArrayContains(matchedIndexes, i) {
			newRelations.AddRelation(oldRelation)
		}
	}

	newRelations.AddRelations(replacements)

	return newRelations
}

// return only the replacements
func (transformer *simpleRelationTransformer) Extract(relationSet *SimpleRelationSet) *SimpleRelationSet {

	_, replacements := transformer.matchAllTransformations(relationSet.relations)
	return NewSimpleRelationSet2(replacements)
}

// only add the replacements to the original relations
func (transformer *simpleRelationTransformer) Append(relationSet *SimpleRelationSet) *SimpleRelationSet {

	_, replacements := transformer.matchAllTransformations(relationSet.relations)

	newRelations := NewSimpleRelationSet2(relationSet.GetRelations())
	newRelations.AddRelations(replacements)

	return newRelations
}

// Attempts all transformations on all relations
// Returns the indexes of the matched relations, and the replacements that were created
func (transformer *simpleRelationTransformer) matchAllTransformations(relations []SimpleRelation) ([]int, []SimpleRelation){

	matchedIndexes := []int{}
	replacements := []SimpleRelation{}

	for _, transformation := range transformer.transformations {

		newMatchedIndexes, newReplacements := transformer.matchSingleTransformation(relations, transformation)
		matchedIndexes = append(matchedIndexes, newMatchedIndexes...)
		replacements = append(replacements, newReplacements...)
	}

	return intArrayDeduplicate(matchedIndexes), replacements
}

// Attempts to match a single transformation
// Returns the indexes of matched relations, and the replacements
func (transformer *simpleRelationTransformer) matchSingleTransformation(relations []SimpleRelation, transformation SimpleRelationTransformation) ([]int, []SimpleRelation){

	matchedIndexes, binding := transformer.matcher.matchSubjectsToPatterns(relations, transformation.Pattern)

	replacements := []SimpleRelation{}
	if len(matchedIndexes) > 0 {
		replacements = append(replacements, transformer.createReplacements(transformation.Replacement, binding)...)
	}

	return matchedIndexes, replacements
}

func (transformer *simpleRelationTransformer) createReplacements(relations []SimpleRelation, bindings SimpleBinding) []SimpleRelation {

	replacements := []SimpleRelation{}

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