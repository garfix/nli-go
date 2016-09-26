package example3

type SimpleRuleBase struct {
	rules []SimpleRule
}

func NewSimpleRuleBase(rules []SimpleRule) *SimpleRuleBase {
	return &SimpleRuleBase{rules: rules}
}

func (ruleBase *SimpleRuleBase) Bind(goal SimpleRelation) ([][]SimpleRelation, []SimpleBinding) {

	matcher := NewSimpleRelationMatcher()
	subgoalRelationSets := [][]SimpleRelation{}
	subgoalBindings := []SimpleBinding{}

	for _, rule := range ruleBase.rules {

		binding := SimpleBinding{}

		// match goal
		simpleBinding, success := matcher.matchRelationToRelation(goal, rule.Goal, binding)
		if !success {
			continue
		}

		// create relation set from the goal conditions
		subgoalRelationSet := []SimpleRelation{}

		for _, condition := range rule.Pattern {
			subgoalRelationSet = append(subgoalRelationSet, matcher.bindSingleRelationSingleBinding(condition, simpleBinding))
		}

		subgoalRelationSets = append(subgoalRelationSets, subgoalRelationSet)
		subgoalBindings = append(subgoalBindings, simpleBinding)
	}

	return subgoalRelationSets, subgoalBindings
}
