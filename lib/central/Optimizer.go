package central

import (
	"nli-go/lib/mentalese"
	"sort"
)

const worst_cost = 100000000.0

// The optimizer reorders the relations in a set to minimize the number of tuples retrieved from the fact bases
//
// It implements ideas from "Efficient Processing of Interactive Relational Database Queries Expressed in Logic" by David H.D. Warren
type Optimizer struct {

}

func (optimizer Optimizer) Optimize(set mentalese.RelationSet, factBases []mentalese.FactBase) mentalese.RelationSet {

	costs := []RelationCost{}

	// edge case
	if len(set) <= 1 {
		return set
	}

	// calculate cost per relation
	for relationIndex, relation := range set {

		relationCost := float32(0.0)
		usedInAnyFactBase := false

		for _, factBase := range factBases {

			factBaseCost := float32(0.0)

			stats := factBase.GetStatistics()

			relationStats, usedInFactBase := stats[relation.Predicate]

			if usedInFactBase {

				usedInAnyFactBase = true
				product := 1

				for columnIndex, distinctValues := range relationStats.DistinctValues {
					if !relation.Arguments[columnIndex].IsVariable() && !relation.Arguments[columnIndex].IsAnonymousVariable() {
						product *= distinctValues
					}
				}

				factBaseCost = float32(relationStats.Size) / float32(product)

			} else {

				// no stats for relation, but does it occur at all?
				for _, mapping := range factBase.GetMappings() {
					if (mapping.DsSource.Predicate == relation.Predicate) {

						// yes it does, and since we don't know the cost, we must presume the worst
						factBaseCost = worst_cost
					}
				}
			}

			relationCost += factBaseCost
		}

		// relations that are not used in any fact base are placed last. thus, they may benefit from variable bindings applied on the fact base relations
		if !usedInAnyFactBase {
			relationCost = worst_cost
		}

		costs = append(costs, RelationCost{relationCost, relationIndex })
	}

	// sort costs
	sort.Sort(RelationCosts(costs))

	// created new set sorted by cost
	orderedSet := mentalese.RelationSet{}
	for _, cost := range costs {
		orderedSet = append(orderedSet, set[cost.relationIndex])
	}

	return orderedSet
}

type RelationCost struct {
	Cost float32
	relationIndex int
}

type RelationCosts []RelationCost

func (s RelationCosts) Len() int {
	return len(s)
}
func (s RelationCosts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s RelationCosts) Less(i, j int) bool {
	return s[i].Cost < s[j].Cost
}