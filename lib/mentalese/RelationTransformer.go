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
func (transformer *RelationTransformer) Replace(rules []Rule, relationSet RelationSet) RelationSet {

	// replace the relations embedded in quants
	replacedSet := transformer.replaceRelations(rules, relationSet, Binding{})

	return replacedSet
}

func (transformer *RelationTransformer) replaceRelations(transformations []Rule, relationSet RelationSet, binding Binding) RelationSet {

	replacedSet := RelationSet{}
	for _, relation := range relationSet {

		// replace inside hierarchical relations
		deepRelation := NewRelation(relation.Predicate, relation.Arguments)

		for i, argument := range deepRelation.Arguments {
			if argument.IsRelationSet() {
				deepRelation.Arguments[i] = NewRelationSet(transformer.replaceRelations(transformations, argument.TermValueRelationSet, binding))
			}
		}

		// replace according to rules
		newRelations := RelationSet{ }

		found := false
		for _, rule := range transformations {
			aBinding, ok := transformer.matcher.MatchTwoRelations(rule.Goal, deepRelation, Binding{})
			if  ok {
				boundRelations := rule.Pattern.BindSingle(aBinding)
				newRelations = append(newRelations, boundRelations...)
				found = true
			}
		}

		if !found {
			newRelations = append(newRelations, deepRelation)
		}

		replacedSet = append(replacedSet, newRelations...)
	}

	return replacedSet
}
//
//// Attempts all transformations on all relations
//// Returns the Indexes of the matched relations, and the replacements that were created, each in a single set
//func (transformer *RelationTransformer) matchAllTransformations(transformations []Rule, haystackSet RelationSet, unScopedRelations RelationSet) ([]int, RelationSet) {
//
//	transformer.log.StartDebug("matchAllTransformations", transformations)
//
//	var matchedIndexes []int
//	var replacements RelationSet
//
//	for _, transformation := range transformations {
//
//		conditionBindings := transformer.bindCondition(transformation, unScopedRelations)
//
//		if len(conditionBindings) > 0 {
//			for _, conditionBinding := range conditionBindings {
//
//				// each transformation application is completely independent from the others
//				bindings, newIndexes, _, match := transformer.matcher.MatchSequenceToSetWithIndexes(transformation.Pattern, haystackSet, conditionBinding)
//				if match {
//					matchedIndexes = append(matchedIndexes, newIndexes...)
//					for _, binding := range bindings {
//						replacements = append(replacements, transformer.createReplacements(transformation.Pattern, binding)...)
//					}
//				}
//			}
//		}
//	}
//
//	matchedIndexes = common.IntArrayDeduplicate(matchedIndexes)
//
//	transformer.log.EndDebug("matchAllTransformations", matchedIndexes, replacements)
//
//	return matchedIndexes, replacements
//}
//
//func (transformer *RelationTransformer) bindCondition(transformation Rule, unScopedRelations RelationSet) Bindings {
//
//	bindings, _ := transformer.matcher.MatchRelationToSet(transformation.Goal, unScopedRelations, Binding{})
//
//	return bindings
//}
//
//func (transformer *RelationTransformer) createReplacements(relation Relation, binding Binding) RelationSet {
//
//	replacements := RelationSet{}
//
//	transformer.log.StartDebug("createReplacements", relation, binding)
//
//	newRelation := Relation{}
//	newRelation.Predicate = relation.Predicate
//
//	for _, argument := range relation.Arguments {
//
//		arg := argument.Copy()
//
//		if argument.IsRelationSet() {
//
//			arg.TermValueRelationSet = transformer.createReplacements(argument.TermValueRelationSet, binding)
//
//		} else if argument.IsVariable() {
//			value, found := binding[argument.String()]
//			if found {
//				arg = value
//			} else {
//				// replacement could not be bound, leave variable unchanged
//			}
//		}
//
//		newRelation.Arguments = append(newRelation.Arguments, arg)
//	}
//
//	replacements = append(replacements, newRelation)
//
//	transformer.log.EndDebug("createReplacements", replacements)
//
//	return replacements
//}
