package knowledge

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type InMemoryRuleBase struct {
	KnowledgeBaseCore
	originalRules mentalese.Rules
	rules mentalese.Rules
	writeList []string
	storage *common.FileStorage
	log   *common.SystemLog
	changed bool
}

func NewInMemoryRuleBase(name string, rules mentalese.Rules, writeList []string, storage *common.FileStorage, log *common.SystemLog) *InMemoryRuleBase {
	ruleBase := InMemoryRuleBase{
		KnowledgeBaseCore: KnowledgeBaseCore{ Name: name},
		originalRules: rules,
		rules: rules.Copy(),
		writeList: writeList,
		storage: storage,
		log: log,
		changed: false,
	}

	if storage != nil {
		storage.Read(&ruleBase.rules)
	}

	return &ruleBase
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
	ruleBase.changed = true
}

func (ruleBase *InMemoryRuleBase) ResetSession() {
	ruleBase.rules = ruleBase.originalRules.Copy()
	ruleBase.changed = true
}

func (ruleBase *InMemoryRuleBase) Persist() {
	if ruleBase.storage != nil {
		if ruleBase.changed {
			ruleBase.storage.Write(ruleBase.rules)
		}
	}
}
