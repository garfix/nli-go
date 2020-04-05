package knowledge

import (
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

func (ruleBase *InMemoryRuleBase) Bind(goal mentalese.Relation, binding mentalese.Binding) ([]mentalese.RelationSet, mentalese.Bindings) {

	ruleBase.log.StartDebug("RuleBase BindSingle", goal)

	matcher := mentalese.NewRelationMatcher(ruleBase.log)
	subgoalRelationSets := []mentalese.RelationSet{}
	subgoalBindings := mentalese.Bindings{}

	for _, rule := range ruleBase.rules {

		// match goal
		aBinding, match := matcher.MatchTwoRelations(goal, rule.Goal, binding)
		if match {
			bBinding, _ := matcher.MatchTwoRelations(rule.Goal, goal, mentalese.Binding{})
			boundRule := rule.ImportBinding(bBinding)
			subgoalRelationSets = append(subgoalRelationSets, boundRule.Pattern)
			subgoalBindings = append(subgoalBindings, aBinding)
		}
	}

	ruleBase.log.EndDebug("RuleBase BindSingle", subgoalRelationSets, subgoalBindings)

	return subgoalRelationSets, subgoalBindings
}
