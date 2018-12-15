package knowledge

import (
	"nli-go/lib/mentalese"
)

// A relation group is a small set of relations that, together, form the input to a single knowledge base
// The knowledge base index refers to the array allKnowledgeBases in ProblemSolver.
// It has a cost when applied to its knowledge base. This cost depends on the number of records visited.
// An array of relation groups can be sorted (by increasing cost)

type RelationGroup struct {
	Relations mentalese.RelationSet
	KnowledgeBaseName string
	Cost float64
}


func (s RelationGroup) String() string {

	str := s.Relations.String() + "@" + s.KnowledgeBaseName

	return str
}


func (s RelationGroup) Equals(t RelationGroup) bool {
	return s.Cost == t.Cost &&
		s.KnowledgeBaseName == t.KnowledgeBaseName &&
		s.Relations.Equals(t.Relations)
}