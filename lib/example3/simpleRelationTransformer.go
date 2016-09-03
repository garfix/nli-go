package example3

type simpleRelationTransformer struct {
	transformations []SimpleRelationTransformation
	matcher simpleRelationMatcher
}

func NewSimpleRelationTransformer(transformations[]SimpleRelationTransformation) *simpleRelationTransformer {
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

	matchedIndexes, boundVariables := transformer.matcher.matchRelations(relations, transformation.Pattern)

	replacements := []SimpleRelation{}
	if len(matchedIndexes) > 0 {
		replacements = append(replacements, transformer.createReplacements(transformation.Replacement, boundVariables)...)
	}

	return matchedIndexes, replacements
}

func (transformer *simpleRelationTransformer) createReplacements(relations []SimpleRelation, boundVariables map[string]SimpleTerm) []SimpleRelation {

	replacements := []SimpleRelation{}

	for _, relation := range relations {

		for i, argument := range relation.Arguments {

			if argument.IsVariable() {
				value, found := boundVariables[argument.AsKey()]
				if found {
					relation.Arguments[i] = value
				} else {
					// replacement could not be bound!
				}
			}
		}

		replacements = append(replacements, relation)
	}

	return replacements
}