package knowledge

import (
	"nli-go/lib/mentalese"
)

type FactBase interface {
	KnowledgeBase
	MatchRelationToDatabase(needleRelation mentalese.Relation) mentalese.Bindings
	Assert(relation mentalese.Relation)
	Retract(relation mentalese.Relation)
	GetMappings() []mentalese.RelationTransformation
	GetWriteMappings() []mentalese.RelationTransformation
	GetStatistics() mentalese.DbStats
	GetEntities() mentalese.Entities
}

const worst_cost = 100000000.0

func getFactBaseMatchingGroups(matcher *mentalese.RelationMatcher, set mentalese.RelationSet, factBase FactBase, keyCabinet *mentalese.KeyCabinet) []RelationGroup {

	matchingGroups := []RelationGroup{}

	matchingGroups = append(matchingGroups, getFactBaseReadGroups(matcher, set, factBase, keyCabinet)...)

	matchingGroups = append(matchingGroups, getFactBaseWriteGroups(matcher, set, factBase, keyCabinet, mentalese.PredicateAssert)...)
	matchingGroups = append(matchingGroups, getFactBaseWriteGroups(matcher, set, factBase, keyCabinet, mentalese.PredicateRetract)...)

	return matchingGroups
}

func getFactBaseReadGroups(matcher *mentalese.RelationMatcher, set mentalese.RelationSet, factBase FactBase, keyCabinet *mentalese.KeyCabinet) []RelationGroup {

	matchingGroups := []RelationGroup{}

	for _, mapping := range factBase.GetMappings() {

		bindings, _, indexesPerNode, match := matcher.MatchSequenceToSetWithIndexes(mapping.Pattern, set, mentalese.Binding{})

		if match {

			for i := range bindings {

				binding := bindings[i]
				indexes := indexesPerNode[i].Indexes

				matchingRelations := mentalese.RelationSet{}
				for _, i := range indexes {
					matchingRelations = append(matchingRelations, set[i])
				}

				boundReplacement := mapping.Replacement.BindSingle(binding)

				keyBoundReplacement := keyCabinet.BindToRelationSet(boundReplacement, factBase.GetName())

				cost := CalculateCost(keyBoundReplacement, factBase.GetStatistics())

				matchingGroups = append(matchingGroups, RelationGroup{matchingRelations, factBase.GetName(), cost})
			}
		}
	}

	return matchingGroups
}

func getFactBaseWriteGroups(matcher *mentalese.RelationMatcher, set mentalese.RelationSet, factBase FactBase, keyCabinet *mentalese.KeyCabinet, predicate string) []RelationGroup {

	matchingGroups := []RelationGroup{}

	for _, relation := range set {
		if relation.Predicate == mentalese.PredicateAssert || relation.Predicate == mentalese.PredicateRetract {
			content := relation.Arguments[0].TermValueRelationSet

			for _, mapping := range factBase.GetWriteMappings() {

				_, _, indexesPerNode, match := matcher.MatchSequenceToSetWithIndexes(mapping.Pattern, content, mentalese.Binding{})

				if match {

					indexes := indexesPerNode[0].Indexes

					matchingRelations := mentalese.RelationSet{}
					for _, i := range indexes {
						matchingRelations = append(matchingRelations, set[i])
					}

					cost := worst_cost

					matchingGroups = append(matchingGroups, RelationGroup{matchingRelations, factBase.GetName(), cost})
				}
			}
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