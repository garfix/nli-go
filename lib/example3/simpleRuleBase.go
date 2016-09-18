package example3

type SimpleRuleBase struct {
	rules []SimpleRule
}

func NewSimpleRuleBase(rules []SimpleRule) *SimpleRuleBase {
	return &SimpleRuleBase{rules: rules}
}

func (ruleBase *SimpleRuleBase) Bind(goal SimpleRelation) map[string]SimpleTerm {

}
