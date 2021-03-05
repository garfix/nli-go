package knowledge

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type InMemoryRuleBase struct {
	KnowledgeBaseCore
	rules []mentalese.Rule
	writeList []string
	log   *common.SystemLog
}

func NewInMemoryRuleBase(name string, rules []mentalese.Rule, writeList []string, log *common.SystemLog) *InMemoryRuleBase {
	return &InMemoryRuleBase{
		KnowledgeBaseCore: KnowledgeBaseCore{ Name: name},
		rules: rules,
		writeList: writeList,
		log: log,
	}
}

func (ruleBase *InMemoryRuleBase) GetPredicates() []string {
	predicates := []string{}
	for _, rule := range ruleBase.rules {
		predicates = append(predicates, rule.Goal.Predicate)
	}
	return predicates
}

func (ruleBase *InMemoryRuleBase) GetRules() []mentalese.Rule {
	return ruleBase.rules
}

func (ruleBase *InMemoryRuleBase) GetWritablePredicates() []string {
	return ruleBase.writeList
}

func (ruleBase *InMemoryRuleBase) GetRulesForRelation(goal mentalese.Relation, binding mentalese.Binding) []mentalese.Rule {

	matcher := central.NewRelationMatcher(ruleBase.log)
	rules := []mentalese.Rule{}

	for _, rule := range ruleBase.rules {

		// match goal
		_, match := matcher.MatchTwoRelations(goal, rule.Goal, binding)
		if match {
			rules = append(rules, rule)
		}
	}

	return rules
}

func (ruleBase *InMemoryRuleBase) Assert(rule mentalese.Rule) {
	ruleBase.rules = append(ruleBase.rules, rule)
}