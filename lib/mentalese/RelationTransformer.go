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
	replacedSet := transformer.replaceRelations(rules, relationSet, NewBinding())

	return replacedSet
}

func (transformer *RelationTransformer) replaceRelations(transformations []Rule, relationSet RelationSet, binding Binding) RelationSet {

	replacedSet := RelationSet{}
	for _, relation := range relationSet {

		// replace inside hierarchical relations
		deepRelation := NewRelation(true, relation.Predicate, relation.Arguments)

		for i, argument := range deepRelation.Arguments {
			if argument.IsRelationSet() {
				deepRelation.Arguments[i] = NewTermRelationSet(transformer.replaceRelations(transformations, argument.TermValueRelationSet, binding))
			} else if argument.IsRule() {
				// no need for implementation
			} else if argument.IsList() {
				// no need for implementation
			}
		}

		// replace according to rules
		newRelations := RelationSet{ }

		found := false
		for _, rule := range transformations {
			aBinding, ok := transformer.matcher.MatchTwoRelations(rule.Goal, deepRelation, NewBinding())
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
