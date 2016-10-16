package knowledge

import (
	"nli-go/lib/mentalese"
)

type RuleBase struct {
	rules []mentalese.Rule
}

func NewRuleBase(rules []mentalese.Rule) *RuleBase {
	return &RuleBase{rules: rules}
}

func (ruleBase *RuleBase) Bind(goal mentalese.Relation) ([]mentalese.RelationSet, []mentalese.Binding) {

	matcher := mentalese.NewRelationMatcher()
	subgoalRelationSets := []mentalese.RelationSet{}
	subgoalBindings := []mentalese.Binding{}

	for _, rule := range ruleBase.rules {

		binding := mentalese.Binding{}

		// match goal
		aBinding, match := matcher.MatchTwoRelations(goal, rule.Goal, binding)
		if !match {
			continue
		}

		// create relation set from the goal conditions
		subgoalRelationSet := []mentalese.Relation{}

		for _, condition := range rule.Pattern {
			subgoalRelationSet = append(subgoalRelationSet, matcher.BindSingleRelationSingleBinding(condition, aBinding))
		}

		subgoalRelationSets = append(subgoalRelationSets, subgoalRelationSet)
		subgoalBindings = append(subgoalBindings, aBinding)
	}

	return subgoalRelationSets, subgoalBindings
}
