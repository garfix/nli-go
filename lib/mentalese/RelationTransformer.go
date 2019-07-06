package mentalese

import (
	"nli-go/lib/common"
)

type RelationTransformer struct {
	matcher *RelationMatcher
	log     *common.SystemLog
}

// using transformations
func NewRelationTransformer(matcher *RelationMatcher, log *common.SystemLog) *RelationTransformer {
	return &RelationTransformer{
		matcher: matcher,
		log:     log,
	}
}

// return the original relations, but replace the ones that have matched with their replacements
func (transformer *RelationTransformer) Replace(transformations []RelationTransformation, relationSet RelationSet) RelationSet {

	// prepare flat set of relations where all scopes are removed
	unScopedRelationSet := relationSet.UnScope()

	// replace the relations embedded in quants
	replacedSet := transformer.replaceEmbeddedRelations(transformations, relationSet)

	matchedIndexes, replacements := transformer.matchAllTransformations(transformations, replacedSet, unScopedRelationSet)
	newRelations := RelationSet{}

	for i, oldRelation := range replacedSet {
		if !common.IntArrayContains(matchedIndexes, i) {
			newRelations = append(newRelations, oldRelation)
		}
	}

	newRelations = append(newRelations, replacements...)

	return newRelations
}

func (transformer *RelationTransformer) replaceEmbeddedRelations(transformations []RelationTransformation, relationSet RelationSet) RelationSet {

	// replace inside hierarchical relations
	replacedSet := RelationSet{}
	for _, relation := range relationSet {

		if relation.Predicate == PredicateQuant {
			replacedRelation := relation.Copy()
			replacedRelation.Arguments[QuantRangeIndex].TermValueRelationSet =
				transformer.Replace(transformations, relation.Arguments[QuantRangeIndex].TermValueRelationSet)
			replacedRelation.Arguments[QuantQuantifierIndex].TermValueRelationSet =
				transformer.Replace(transformations, relation.Arguments[QuantQuantifierIndex].TermValueRelationSet)
			replacedRelation.Arguments[QuantScopeIndex].TermValueRelationSet =
				transformer.Replace(transformations, relation.Arguments[QuantScopeIndex].TermValueRelationSet)
			replacedSet = append(replacedSet, replacedRelation)
		} else {
			replacedSet = append(replacedSet, relation)
		}
	}

	return replacedSet
}

// Attempts all transformations on all relations
// Returns the Indexes of the matched relations, and the replacements that were created, each in a single set
func (transformer *RelationTransformer) matchAllTransformations(transformations []RelationTransformation, haystackSet RelationSet, unScopedRelations RelationSet) ([]int, RelationSet) {

	transformer.log.StartDebug("matchAllTransformations", transformations)

	var matchedIndexes []int
	var replacements RelationSet

	for _, transformation := range transformations {

		conditionBindings, ok := transformer.bindCondition(transformation, unScopedRelations)

		if ok {
			for _, conditionBinding := range conditionBindings {

				// each transformation application is completely independent from the others
				bindings, newIndexes, _, match := transformer.matcher.MatchSequenceToSetWithIndexes(transformation.Pattern, haystackSet, conditionBinding)
				if match {
					matchedIndexes = append(matchedIndexes, newIndexes...)
					for _, binding := range bindings {
						replacements = append(replacements, transformer.createReplacements(transformation.Replacement, binding)...)
					}
				}
			}
		}
	}

	matchedIndexes = common.IntArrayDeduplicate(matchedIndexes)

	transformer.log.EndDebug("matchAllTransformations", matchedIndexes, replacements)

	return matchedIndexes, replacements
}

func (transformer *RelationTransformer) bindCondition(transformation RelationTransformation, unScopedRelations RelationSet) (Bindings, bool) {

	bindings := Bindings{{}}

	ok := true

	if len(transformation.Condition) > 0 {

		bindings, ok = transformer.matcher.MatchSequenceToSet(transformation.Condition, unScopedRelations, Binding{})
	}

	return bindings, ok
}

func (transformer *RelationTransformer) createReplacements(relations RelationSet, bindings Binding) RelationSet {

	replacements := RelationSet{}

	transformer.log.StartDebug("createReplacements", relations, bindings)

	for _, relation := range relations {

		newRelation := Relation{}
		newRelation.Predicate = relation.Predicate

		for _, argument := range relation.Arguments {

			arg := argument.Copy()

			if argument.IsRelationSet() {

				arg.TermValueRelationSet = transformer.createReplacements(argument.TermValueRelationSet, bindings)

			} else if argument.IsVariable() {
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

	transformer.log.EndDebug("createReplacements", replacements)

	return replacements
}
