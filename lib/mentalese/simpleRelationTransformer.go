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

// return the original relations, but replace the ones that have matched with their replacements
func (transformer *simpleRelationTransformer) Replace(relationSet SimpleRelationSet) SimpleRelationSet {

	matchedIndexes, replacements := transformer.matchAllTransformations(relationSet)
	newRelations := SimpleRelationSet{}

	for i, oldRelation := range relationSet  {
		if !common.IntArrayContains(matchedIndexes, i) {
			newRelations = append(newRelations, oldRelation)
		}
	}

	newRelations = append(newRelations, replacements...)

	return newRelations
}

// Try to match all transformations to relationSet, and return the replacements that resulted from the transformations
func (transformer *simpleRelationTransformer) Extract(relationSet SimpleRelationSet) (SimpleRelationSet) {

	common.LogTree("Extract", relationSet)

	_, replacements := transformer.matchAllTransformations(relationSet)

	common.LogTree("Extract", replacements)

	return replacements
}

// only add the replacements to the original relations
func (transformer *simpleRelationTransformer) Append(relationSet SimpleRelationSet) SimpleRelationSet {

	_, replacements := transformer.matchAllTransformations(relationSet)

	newRelations := SimpleRelationSet{}
	newRelations = append(newRelations, relationSet...)
	newRelations = append(newRelations, replacements...)

	return newRelations
}

// Attempts all transformations on all relations
// Returns the indexes of the matched relations, and the replacements that were created, each in a single set
func (transformer *simpleRelationTransformer) matchAllTransformations(haystackSet SimpleRelationSet) ([]int, SimpleRelationSet){

	common.LogTree("matchAllTransformations", haystackSet)

	matchedIndexes := []int{}
	replacements := SimpleRelationSet{}

	for _, transformation := range transformer.transformations {

		// each transformation application is completely independent from the others
		newIndexes, binding, match := transformer.matcher.MatchSequenceToSet(transformation.Pattern, haystackSet, SimpleBinding{})
		if match {
			matchedIndexes = append(matchedIndexes, common.IntArrayDeduplicate(newIndexes)...)
			replacements = append(replacements, transformer.createReplacements(transformation.Replacement, binding)...)
		}
	}

	common.LogTree("matchAllTransformations", matchedIndexes, replacements)

	return matchedIndexes, replacements
}

func (transformer *simpleRelationTransformer) createReplacements(relations SimpleRelationSet, bindings SimpleBinding) SimpleRelationSet {

	replacements := SimpleRelationSet{}

	common.LogTree("createReplacements", relations, bindings)

	for _, relation := range relations {

		newRelation := SimpleRelation{}
		newRelation.Predicate = relation.Predicate

		for _, argument := range relation.Arguments {

			arg := argument

			if argument.IsVariable() {
				value, found := bindings[argument.String()]
				if found {
					arg = value
				} else {
					// replacement could not be bound, leave variable unchanged
				}
			}

			newRelation.Arguments = append(newRelation.Arguments, arg)
		}

		replacements = append(replacements, newRelation)
	}

	common.LogTree("createReplacements", replacements)

	return replacements
}