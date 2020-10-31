package knowledge

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type InMemoryRuleBase struct {
	KnowledgeBaseCore
	rules []mentalese.Rule
	log   *common.SystemLog
}

func NewInMemoryRuleBase(name string, rules []mentalese.Rule, log *common.SystemLog) *InMemoryRuleBase {
	return &InMemoryRuleBase{
		KnowledgeBaseCore: KnowledgeBaseCore{ Name: name},
		rules: rules,
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

func (ruleBase *InMemoryRuleBase) Bind(goal mentalese.Relation, binding mentalese.Binding) ([]mentalese.RelationSet, mentalese.BindingSet) {

	matcher := central.NewRelationMatcher(ruleBase.log)
	subgoalRelationSets := []mentalese.RelationSet{}
	subgoalBindings := mentalese.NewBindingSet()

	for _, rule := range ruleBase.rules {

		// match goal
		aBinding, match := matcher.MatchTwoRelations(goal, rule.Goal, binding)
		if match {
			bBinding, _ := matcher.MatchTwoRelations(rule.Goal, goal, mentalese.NewBinding())
			boundRule := rule.BindSingle(bBinding)
			boundRule = boundRule.InstantiateUnboundVariables(aBinding)
			subgoalRelationSets = append(subgoalRelationSets, boundRule.Pattern)
			subgoalBindings.Add(aBinding)
		}
	}

	return subgoalRelationSets, subgoalBindings
}

func (ruleBase *InMemoryRuleBase) Assert(rule mentalese.Rule) {
	ruleBase.rules = append(ruleBase.rules, rule)
}