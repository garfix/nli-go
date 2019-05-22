package knowledge

import (
	"nli-go/lib/mentalese"
)

type FactBase interface {
	KnowledgeBase
	MatchRelationToDatabase(needleRelation mentalese.Relation) []mentalese.Binding
	GetMappings() []mentalese.RelationTransformation
	GetStatistics() mentalese.DbStats
	GetEntities() mentalese.Entities
}

const worst_cost = 100000000.0

func getFactBaseMatchingGroups(matcher *mentalese.RelationMatcher, set mentalese.RelationSet, factBase FactBase, nameStore *mentalese.ResolvedNameStore) []RelationGroup {

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

			boundReplacement := mapping.Replacement.BindRelationSetSingleBinding(binding)

			keyBoundReplacement := nameStore.BindToRelationSet(boundReplacement, factBase.GetName())

			cost := CalculateCost(keyBoundReplacement, factBase.GetStatistics())

			matchingGroups = append(matchingGroups, RelationGroup{matchingRelations, factBase.GetName(), cost})
		}
	}

	return matchingGroups
}

// The cost of a relation set that is to be resolved by a fact base. The fact base brings in the stats.
// The higher the cost, the later it should be executed. Lower cost is better.
// If no cost can be calculated, the cost is assumed to be high.
// The function was found in "Efficient processing of interactive relational database queries expressed in logic" by David Warren
func CalculateCost(boundReplacement mentalese.RelationSet, stats mentalese.DbStats) float64 {

	cost := float64(0.0)

	for _, replacementRelation := range boundReplacement {

		relationStats, usedInFactBase := stats[replacementRelation.Predicate]

		if usedInFactBase {
			product := 1
			for columnIndex, distinctValues := range relationStats.DistinctValues {

				if !replacementRelation.Arguments[columnIndex].IsVariable() && !replacementRelation.Arguments[columnIndex].IsAnonymousVariable() {
					product *= distinctValues
				}
			}

			// the cost of a single relation is its domain size divided ("softened") by the product
			// of the domain sizes of its instantiated arguments.
			cost += float64(relationStats.Size) / float64(product)
		} else {
			// no cost available: presume high cost
			cost += worst_cost
		}
	}

	return cost
}