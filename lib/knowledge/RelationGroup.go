package knowledge

import (
	"nli-go/lib/mentalese"
)

// A relation group is a small set of relations that, together, form the input to a single knowledge base
// The knowledge base index refers to the array allKnowledgeBases in ProblemSolver.

type RelationGroup struct {
	Relations mentalese.RelationSet
	KnowledgeBaseName string
}


func (s RelationGroup) String() string {

	str := s.Relations.String() + "@" + s.KnowledgeBaseName

	return str
}


func (s RelationGroup) Equals(t RelationGroup) bool {
	return s.KnowledgeBaseName == t.KnowledgeBaseName &&
		s.Relations.Equals(t.Relations)
}