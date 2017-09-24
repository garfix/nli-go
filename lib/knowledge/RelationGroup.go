package knowledge

import "nli-go/lib/mentalese"

// A relation group is a small set of relations that, together, form the input to a single knowledge base
// The knowledge base index refers to the array allKnowledgeBases in ProblemSolver.
// It has a cost when applied to its knowledge base. This cost depends on the number of records visited.
// An array of relation groups can be sorted (by increasing cost)

type RelationGroup struct {
	Relations mentalese.RelationSet
	KnowledgeBaseIndex int
	Cost float64
}

type RelationGroups []RelationGroup

func (s RelationGroups) Len() int {
	return len(s)
}

func (s RelationGroups) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s RelationGroups) Less(i, j int) bool {
	return s[i].Cost < s[j].Cost
}

func (s RelationGroups) GetCombinedRelations() mentalese.RelationSet {

	relations := mentalese.RelationSet{}

	for _, group := range s {
		relations = append(relations, group.Relations...)
	}

	return relations
}

func (s RelationGroups) GetTotalRelationCount() int {

	count := 0

	for _, group := range s {
		count += len(group.Relations)
	}

	return count
}