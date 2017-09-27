package knowledge

import (
	"nli-go/lib/mentalese"
	"strconv"
)

// A relation group is a small set of relations that, together, form the input to a single knowledge base
// The knowledge base index refers to the array allKnowledgeBases in ProblemSolver.
// It has a cost when applied to its knowledge base. This cost depends on the number of records visited.
// An array of relation groups can be sorted (by increasing cost)

type RelationGroup struct {
	Relations mentalese.RelationSet
	KnowledgeBaseIndex int
	Cost float64
}


func (s RelationGroup) String() string {

	str := s.Relations.String() + "@" + strconv.Itoa(s.KnowledgeBaseIndex)

	return str
}


func (s RelationGroup) Equals(t RelationGroup) bool {
	return s.Cost == t.Cost &&
		s.KnowledgeBaseIndex == t.KnowledgeBaseIndex &&
		s.Relations.Equals(t.Relations)
}