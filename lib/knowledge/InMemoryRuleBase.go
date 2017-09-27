package knowledge


import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type InMemoryRuleBase struct {
	rules []mentalese.Rule
	log   *common.SystemLog
}

func NewRuleBase(rules []mentalese.Rule, log *common.SystemLog) RuleBase {
	return InMemoryRuleBase{rules: rules, log: log}
}

func (ruleBase InMemoryRuleBase) GetMatchingGroups(set mentalese.RelationSet, knowledgeBaseIndex int) []RelationGroup {

	matchingGroups := []RelationGroup{}

	for _, rule := range ruleBase.rules {
		for _, setRelation := range set {
			if rule.Goal.Predicate == setRelation.Predicate {
// TDOD calculate real costs
				matchingGroups = append(matchingGroups, RelationGroup{mentalese.RelationSet{setRelation}, knowledgeBaseIndex, worst_cost})
				break
			}
		}
	}

	return matchingGroups
}

func (ruleBase InMemoryRuleBase) Bind(goal mentalese.Relation) ([]mentalese.RelationSet, []mentalese.Binding) {

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
