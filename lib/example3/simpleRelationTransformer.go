package example3

type simpleRelationTransformer struct {
	transformations []SimpleRelationTransformation
}

func NewSimpleRelationTransformer(transformations[]SimpleRelationTransformation) *simpleRelationTransformer {
	return &simpleRelationTransformer{transformations: transformations}
}

// return the original relations, but replace the ones that have matched with their replacements
func (transformer *simpleRelationTransformer) Replace(relations []SimpleRelation) []SimpleRelation {
	return []SimpleRelation{}
}

// like replace, but attempt replacement recursively
func (transformer *simpleRelationTransformer) ReplaceRecursively(relations []SimpleRelation) []SimpleRelation {
	return []SimpleRelation{}
}

// return only the replacements
func (transformer *simpleRelationTransformer) Extract(relations []SimpleRelation) []SimpleRelation {

	_, replacements := transformer.matchAllTransformations(relations)
	return replacements
}

// only add the replacements to the original relations
func (transformer *simpleRelationTransformer) Append(relations []SimpleRelation) []SimpleRelation {
	return []SimpleRelation{}
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

	return matchedIndexes, replacements
}

// Attempts to match a single transformation
// Returns the indexes of matched relations, and the replacements
func (transformer *simpleRelationTransformer) matchSingleTransformation(relations []SimpleRelation, transformation SimpleRelationTransformation) ([]int, []SimpleRelation){

	matchedIndexes := []int{}
	replacements := []SimpleRelation{}

	boundVariables := map[string]SimpleTerm{}

	for _, patternRelation := range transformation.Pattern {

		index, newBoundVariables, found := transformer.matchSingleRelation(relations, patternRelation, boundVariables)
		if found {

			boundVariables = newBoundVariables
			matchedIndexes = append(matchedIndexes, index)

		} else {
			return []int{}, []SimpleRelation{}
		}
	}

	replacements = append(replacements, transformer.createReplacements(transformation.Replacement, boundVariables)...)

	return matchedIndexes, replacements
}

// Attempts to match a single pattern relation to a series of relations
func (transformer *simpleRelationTransformer) matchSingleRelation(relations []SimpleRelation, patternRelation SimpleRelation, boundVariables map[string]SimpleTerm) (int, map[string]SimpleTerm, bool) {

	for index, relation := range relations {

		newBoundVariables, matched := transformer.matchRelationToRelation(relation, patternRelation, boundVariables)

		if matched {
			return index, newBoundVariables, true
		}
	}

	return 0, map[string]SimpleTerm{}, false
}

func (transformer *simpleRelationTransformer) matchRelationToRelation(relation SimpleRelation, patternRelation SimpleRelation, boundVariables map[string]SimpleTerm) (map[string]SimpleTerm, bool) {

	success := true

	// predicate
	if relation.Predicate != patternRelation.Predicate {
		success = false
	} else {

		// arguments
		for i, argument := range relation.Arguments {
			newBoundVariables, ok := transformer.bindArgument(argument, patternRelation.Arguments[i], boundVariables)

			if ok {
				boundVariables = newBoundVariables
			} else {
				success = false
				break;
			}
		}
	}

	return boundVariables, success
}

func (transformer *simpleRelationTransformer) bindArgument(argument SimpleTerm, patternRelationArgument SimpleTerm, boundVariables map[string]SimpleTerm) (map[string]SimpleTerm, bool) {

	success := false

	if patternRelationArgument.IsVariable() {

		// variable

		value := SimpleTerm{}

		// does patternRelationArgument occur in boundVariables?
		value, match := boundVariables[patternRelationArgument.AsKey()]
		if match {
			// it does, use the bound variable
			if argument.Equals(value) {
				success = true
			}
		} else {
			// it does not, just assign the actual argument
			boundVariables[patternRelationArgument.AsKey()] = argument
			success = true
		}

	} else {

		// atom, constant

		if argument.Equals(patternRelationArgument) {
			success = true
		}
	}

	return boundVariables, success
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