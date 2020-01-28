package central

import (
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
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
func (optimizer Optimizer) CreateSolutionRoutes(set mentalese.RelationSet, knowledgeBases []knowledge.KnowledgeBase, keyCabinet *mentalese.KeyCabinet) (knowledge.SolutionRoutes, mentalese.RelationSet, bool) {

	routes := knowledge.SolutionRoutes{}

	allRoutes := optimizer.findSolutionRoutes(knowledge.SolutionRoute{}, set, knowledgeBases, keyCabinet)

	remainingRelations := mentalese.RelationSet{}

	longestRoute := knowledge.SolutionRoute{}
	longestRouteRelationCount := 0
	for _, route := range allRoutes {

		relationCount := route.GetTotalRelationCount()

		// find the relation for which no relation group could be found
		if relationCount > longestRouteRelationCount {
			longestRoute = route
			longestRouteRelationCount = longestRoute.GetTotalRelationCount()
		}

		// collect and deduplicate full routes
		if relationCount == len(set) {

			// sort relation groups by cost
//			sort.Sort(knowledge.SolutionRoute(route))

			found := optimizer.isPresent(route, routes)
			if !found {
				routes = append(routes, route)
			}
		}
	}

	if longestRouteRelationCount < len(set) {
		remainingRelations = set.RemoveRelations(longestRoute.GetCombinedRelations())
	}

	ok := len(remainingRelations) == 0

	return routes, remainingRelations, ok
}

func (optimizer Optimizer) isPresent(route knowledge.SolutionRoute, routes []knowledge.SolutionRoute) bool {

	for _, aRoute := range routes {
		if route.Equals(aRoute) {
			return true
		}
	}

	return false
}

func (optimizer Optimizer) findSolutionRoutes(baseRoute knowledge.SolutionRoute, set mentalese.RelationSet, knowledgeBases []knowledge.KnowledgeBase, keyCabinet *mentalese.KeyCabinet) knowledge.SolutionRoutes {

	// find matching groups in all knowledge bases
	matchingGroupSets := [][]knowledge.RelationGroup{}
	for _, factBase := range knowledgeBases {
		matchingGroupSets = append(matchingGroupSets, factBase.GetMatchingGroups(set, keyCabinet))
	}

	// collect groups by relation (relation index => group set, group index)

	// initialize data
	indexedGroups := map[int]map[int]map[int]knowledge.RelationGroup{}
	indexedGroupRelationIndexes := map[int]map[int][]int{}

	for r := range set {
		indexedGroups[r] = map[int]map[int]knowledge.RelationGroup{}
		for s := range matchingGroupSets {
			indexedGroups[r][s] = map[int]knowledge.RelationGroup{}
		}
	}
	for s, groupSet := range matchingGroupSets {
		indexedGroupRelationIndexes[s] = map[int][]int{}
		for g := range groupSet {
			indexedGroupRelationIndexes[s][g] = []int{}
		}
	}

	// fill data structure
	for r, relation := range set {
		for s, groupSet := range matchingGroupSets {
			for g, group := range groupSet {
				for _, rel := range group.Relations {
					if relation.Equals(rel) {
						indexedGroups[r][s][g] = group
						indexedGroupRelationIndexes[s][g] = append(indexedGroupRelationIndexes[s][g], r)
					}
				}
			}
		}
	}

	// create routes
	routes := optimizer.createRoutes(set, 0, []int{}, indexedGroups, indexedGroupRelationIndexes, knowledge.SolutionRoute{})

	return routes
}

func (optimizer Optimizer) createRoutes(set mentalese.RelationSet, r int, handledRelationIndexes []int, indexedGroups map[int]map[int]map[int]knowledge.RelationGroup, indexedGroupRelationIndexes map[int]map[int][]int, solutionRoute knowledge.SolutionRoute) knowledge.SolutionRoutes {

	solutionRoutes := []knowledge.SolutionRoute{}

	if r == len(set) {
		solutionRoutes = append(solutionRoutes, solutionRoute)
	} else if common.IntArrayContains(handledRelationIndexes, r) {
		// relation already present in earlier group, skip it
		newSolutionRoutes := optimizer.createRoutes(set, r + 1, handledRelationIndexes, indexedGroups, indexedGroupRelationIndexes, solutionRoute)
		solutionRoutes = append(solutionRoutes, newSolutionRoutes...)
	} else {
		var found = false
		for s, groupSet := range indexedGroups[r] {
			for g, group := range groupSet {
				found = true
				newSolutionRoute := append(solutionRoute, group)
				newSolutionRelationIndexes := append(handledRelationIndexes, indexedGroupRelationIndexes[s][g]...)
				newSolutionRoutes := optimizer.createRoutes(set, r + 1, newSolutionRelationIndexes, indexedGroups, indexedGroupRelationIndexes, newSolutionRoute)
				solutionRoutes = append(solutionRoutes, newSolutionRoutes...)
			}
		}
		if !found {
			// relation not handled by any kb, skip it
			newSolutionRoutes := optimizer.createRoutes(set, r + 1, handledRelationIndexes, indexedGroups, indexedGroupRelationIndexes, solutionRoute)
			solutionRoutes = append(solutionRoutes, newSolutionRoutes...)
		}
	}

	return solutionRoutes
}
