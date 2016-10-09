package knowledge

import (
	"nli-go/lib/mentalese"
)

type SimpleRuleBase struct {
	rules []mentalese.SimpleRule
}

func NewSimpleRuleBase(rules []mentalese.SimpleRule) *SimpleRuleBase {
	return &SimpleRuleBase{rules: rules}
}

func (ruleBase *SimpleRuleBase) Bind(goal mentalese.SimpleRelation) ([][]mentalese.SimpleRelation, []mentalese.SimpleBinding) {

	matcher := mentalese.NewSimpleRelationMatcher()
	subgoalRelationSets := [][]mentalese.SimpleRelation{}
	subgoalBindings := []mentalese.SimpleBinding{}

	for _, rule := range ruleBase.rules {

		binding := mentalese.SimpleBinding{}

		// match goal
		simpleBinding, success := matcher.MatchNeedleToHaystack(goal, rule.Goal, binding)
		if !success {
			continue
		}

		// create relation set from the goal conditions
		subgoalRelationSet := []mentalese.SimpleRelation{}

		for _, condition := range rule.Pattern {
			subgoalRelationSet = append(subgoalRelationSet, matcher.BindSingleRelationSingleBinding(condition, simpleBinding))
		}

		subgoalRelationSets = append(subgoalRelationSets, subgoalRelationSet)
		subgoalBindings = append(subgoalBindings, simpleBinding)
	}

	return subgoalRelationSets, subgoalBindings
}
