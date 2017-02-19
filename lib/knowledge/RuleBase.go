package knowledge

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
)

type RuleBase struct {
	rules []mentalese.Rule
}

func NewRuleBase(rules []mentalese.Rule) RuleBase {
	return RuleBase{rules: rules}
}

func (ruleBase *RuleBase) Bind(goal mentalese.Relation) ([]mentalese.RelationSet, []mentalese.Binding) {

	common.LogTree("RuleBase Bind", goal);

	matcher := mentalese.NewRelationMatcher()
	subgoalRelationSets := []mentalese.RelationSet{}
	subgoalBindings := []mentalese.Binding{}

	for _, rule := range ruleBase.rules {

		binding := mentalese.Binding{}

		// match goal
		aBinding, match := matcher.MatchTwoRelations(rule.Goal, goal, binding)
		if !match {
			continue
		}

		// create relation set from the goal conditions
		subgoalRelationSet := mentalese.RelationSet{}

		for _, condition := range rule.Pattern {
			//subgoalRelationSet = append(subgoalRelationSet, matcher.BindSingleRelationSingleBinding(condition, aBinding))
			subgoalRelationSet = append(subgoalRelationSet, condition)
		}

		subgoalRelationSets = append(subgoalRelationSets, subgoalRelationSet)
		subgoalBindings = append(subgoalBindings, aBinding)
	}

	common.LogTree("RuleBase Bind", subgoalRelationSets, subgoalBindings);

	return subgoalRelationSets, subgoalBindings
}
