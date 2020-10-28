package api

type KnowledgeBase interface {

	HandlesPredicate(predicate string) bool
	GetName() string
}
