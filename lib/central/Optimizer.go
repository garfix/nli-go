package central

import (
	"nli-go/lib/mentalese"
	"sort"
	"nli-go/lib/knowledge"
)

// The optimizer reorders the relations in a set to minimize the number of tuples retrieved from the fact bases
//
// It implements ideas from "Efficient Processing of Interactive Relational Database Queries Expressed in Logic" by David H.D. Warren
type Optimizer struct {
	matcher               *mentalese.RelationMatcher
}

func NewOptimizer(matcher *mentalese.RelationMatcher) Optimizer {
	return Optimizer{
		matcher: matcher,
	}
}

// Groups set into relation groups based on knowledge base input
// Relations that were not found are placed in the remaining set

func (optimizer Optimizer) CreateRelationGroups(set mentalese.RelationSet, knowledgeBases []knowledge.KnowledgeBase) (knowledge.RelationGroups, mentalese.RelationSet, bool) {

	groups := optimizer.findGroups(set, knowledgeBases)

	// find the relation for which no relation group could be found
	remainingRelations := mentalese.RelationSet{}
	if groups.GetTotalRelationCount() != len(set) {
		remainingRelations = set.RemoveRelations(groups.GetCombinedRelations())
	}

// TODO: quant

	// sort by cost
	sort.Sort(knowledge.RelationGroups(groups))

	ok := len(remainingRelations) == 0

	return groups, remainingRelations, ok
}

func (optimizer Optimizer) findGroups(set mentalese.RelationSet, knowledgeBases []knowledge.KnowledgeBase) knowledge.RelationGroups {

	groups := knowledge.RelationGroups{}

	for i, factBase := range knowledgeBases {
		for _, factBaseGroup := range factBase.GetMatchingGroups(set, i) {

			restOfSet := set.RemoveRelations(factBaseGroup.Relations)
			restGroups := optimizer.findGroups(restOfSet, knowledgeBases)

			groups = knowledge.RelationGroups{factBaseGroup}
			groups = append(groups, restGroups...)

			if groups.GetTotalRelationCount() == len(set) {
				goto end
			}
		}
	}

	end:

	return groups
}

//func (optimizer Optimizer) Optimize(set mentalese.RelationSet, factBases []knowledge.FactBase) mentalese.RelationSet {
//
//	// edge case
//	if len(set) <= 1 {
//		return set
//	}
//
//	// initialize costs
//	costs := []RelationCost{}
//	for i, _ := range set {
//		costs = append(costs, RelationCost{0.0, i})
//	}
//
//	// go through all fact bases
//	for _, factBase := range factBases {
//
//		stats := factBase.GetStatistics()
//
//		// go through all mappings of the fact base
//		for _, mapping := range factBase.GetMappings() {
//
//			// find the indexes of the matching relations
//			_, _, bindingIndexes, match := optimizer.matcher.MatchSequenceToSetWithIndexes(mapping.Pattern, set, mentalese.Binding{})
//
//			if match {
//
//				for _, bindingIndex := range bindingIndexes {
//					boundReplacement := optimizer.matcher.BindRelationSetSingleBinding(mapping.Replacement, bindingIndex.Binding)
//
//					// what is the cost of this mapping?
//					mappingCost := optimizer.getMappingCost(boundReplacement, stats)
//
//					// update costs for matching relations
//					for _, index := range bindingIndex.Indexes {
//						costs[index].Cost = costs[index].Cost + mappingCost
//					}
//				}
//			}
//		}
//	}
//
//	// unknown relations by any knowledge base will get high costs
//	for i, cost := range costs {
//		if cost.Cost == 0.0 {
//			costs[i].Cost = worst_cost
//		}
//	}
//
//	// sort costs
//	sort.Sort(RelationCosts(costs))
//
//	// created new set sorted by cost
//	orderedSet := mentalese.RelationSet{}
//	for _, cost := range costs {
//		orderedSet = append(orderedSet, set[cost.relationIndex])
//	}
//
//	return orderedSet
//}
//
//func (optimizer Optimizer) getMappingCost(boundReplacement mentalese.RelationSet, stats mentalese.DbStats) float64 {
//
//	cost := float64(0.0)
//
//	for _, relation := range boundReplacement {
//
//		relationStats, usedInFactBase := stats[relation.Predicate]
//
//		if usedInFactBase {
//
//			product := 1
//
//			for columnIndex, distinctValues := range relationStats.DistinctValues {
//				if !relation.Arguments[columnIndex].IsVariable() && !relation.Arguments[columnIndex].IsAnonymousVariable() {
//					product *= distinctValues
//				}
//			}
//
//			cost += float64(relationStats.Size) / float64(product)
//		}
//	}
//
//	return cost
//}
//
//type RelationCost struct {
//	Cost float64
//	relationIndex int
//}
//
//type RelationCosts []RelationCost
//
//func (s RelationCosts) Len() int {
//	return len(s)
//}
//func (s RelationCosts) Swap(i, j int) {
//	s[i], s[j] = s[j], s[i]
//}
//func (s RelationCosts) Less(i, j int) bool {
//	return s[i].Cost < s[j].Cost
//}