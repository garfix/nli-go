package api

import "nli-go/lib/mentalese"

// A knowledge bases that uses rules to process relations
type RuleBase interface {
	KnowledgeBase
	Bind(goal mentalese.Relation, binding mentalese.Binding) ([]mentalese.RelationSet, mentalese.BindingSet)
	Assert(rule mentalese.Rule)
}