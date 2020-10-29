package api

// Generic interface for any system that processes a given set of predicates
type KnowledgeBase interface {

	HandlesPredicate(predicate string) bool
	GetName() string
}
