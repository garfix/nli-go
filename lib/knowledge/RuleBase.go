package knowledge

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type RuleBase struct {
	rules []mentalese.Rule
	log   *common.SystemLog
}

func NewRuleBase(rules []mentalese.Rule, log *common.SystemLog) RuleBase {
	return RuleBase{rules: rules, log: log}
}

func (ruleBase *RuleBase) Bind(goal mentalese.Relation) ([]mentalese.RelationSet, []mentalese.Binding) {

	ruleBase.log.StartDebug("RuleBase Bind", goal)

	matcher := mentalese.NewRelationMatcher(ruleBase.log)
	subgoalRelationSets := []mentalese.RelationSet{}
	subgoalBindings := []mentalese.Binding{}

	for _, rule := range ruleBase.rules {

		// match goal
		aBinding, match := matcher.MatchTwoRelations(rule.Goal, goal, mentalese.Binding{})
		if match {
			subgoalRelationSets = append(subgoalRelationSets, rule.Pattern)
			subgoalBindings = append(subgoalBindings, aBinding)
		}
	}

	ruleBase.log.EndDebug("RuleBase Bind", subgoalRelationSets, subgoalBindings)

	return subgoalRelationSets, subgoalBindings
}
