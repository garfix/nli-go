package example3

type SimpleKnowledgeBase interface {
	Bind(goal SimpleRelation) map[string]SimpleTerm
}
