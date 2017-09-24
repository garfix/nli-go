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

// TODO: quant / HIERARCHISCHE CONSTRUCTIES (OR)

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
