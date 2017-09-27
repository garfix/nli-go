package knowledge

import "nli-go/lib/mentalese"

type FactBase interface {
	KnowledgeBase
	Bind(goal []mentalese.Relation) ([]mentalese.Binding, bool)
	GetMappings() []mentalese.RelationTransformation
	GetStatistics() mentalese.DbStats
}

const worst_cost = 100000000.0

func getFactBaseMatchingGroups(matcher *mentalese.RelationMatcher, set mentalese.RelationSet, factBase FactBase, knowledgeBaseIndex int) []RelationGroup {

	matchingGroups := []RelationGroup{}

	for _, mapping := range factBase.GetMappings() {

		bindings, _, indexesPerNode, match := matcher.MatchSequenceToSetWithIndexes(mapping.Pattern, set, mentalese.Binding{})

		if match {

			binding := bindings[0]
			indexes := indexesPerNode[0].Indexes

			matchingRelations := mentalese.RelationSet{}
			for _, i := range indexes {
				matchingRelations = append(matchingRelations, set[i])
			}

			boundReplacement := matcher.BindRelationSetSingleBinding(mapping.Replacement, binding)

			cost := float64(0.0)
			stats := factBase.GetStatistics()
			for _, replacementRelation := range boundReplacement {
				product := 1
				relationStats, usedInFactBase := stats[replacementRelation.Predicate]
				if usedInFactBase {
					for columnIndex, distinctValues := range relationStats.DistinctValues {
						if !replacementRelation.Arguments[columnIndex].IsVariable() && !replacementRelation.Arguments[columnIndex].IsAnonymousVariable() {
							product *= distinctValues
						}
					}
					cost += float64(relationStats.Size) / float64(product)
				} else {
					cost += worst_cost
				}
			}

			matchingGroups = append(matchingGroups, RelationGroup{matchingRelations, knowledgeBaseIndex, cost})
		}
	}

	return matchingGroups
}
