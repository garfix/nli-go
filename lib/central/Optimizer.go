package central

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/knowledge"
	"sort"
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
func (optimizer Optimizer) CreateSolutionRoutes(set mentalese.RelationSet, knowledgeBases []knowledge.KnowledgeBase, nameStore *ResolvedNameStore) (knowledge.SolutionRoutes, mentalese.RelationSet, bool) {

	routes := knowledge.SolutionRoutes{}

	allRoutes := optimizer.findSolutionRoutes(knowledge.SolutionRoute{}, set, knowledgeBases, nameStore)

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
			sort.Sort(knowledge.SolutionRoute(route))

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

func (optimizer Optimizer) findSolutionRoutes(baseRoute knowledge.SolutionRoute, set mentalese.RelationSet, knowledgeBases []knowledge.KnowledgeBase, nameStore *ResolvedNameStore) knowledge.SolutionRoutes {

	routes := knowledge.SolutionRoutes{}

	if len(set) == 0 {
		return routes
	}

	for i, factBase := range knowledgeBases {
		for _, factBaseGroup := range factBase.GetMatchingGroups(set, i) {

			restOfSet := set.RemoveRelations(factBaseGroup.Relations)

			route := baseRoute
			route = append(baseRoute, factBaseGroup)
			routes = append(routes, route)

			restRoutes := optimizer.findSolutionRoutes(route, restOfSet, knowledgeBases, nameStore)
			for _, restRoute := range restRoutes {
				routes = append(routes, restRoute)
			}
		}
	}

	return routes
}
