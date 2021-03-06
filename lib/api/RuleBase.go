package api

import "nli-go/lib/mentalese"

// A knowledge bases that uses rules to process relations
type RuleBase interface {
	KnowledgeBase
	GetPredicates() []string
	GetWritablePredicates() []string
	GetRules() []mentalese.Rule
	GetRulesForRelation(goal mentalese.Relation, binding mentalese.Binding) []mentalese.Rule
	Assert(rule mentalese.Rule)
}