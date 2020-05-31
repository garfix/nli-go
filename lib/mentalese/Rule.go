package mentalese

import "nli-go/lib/common"

type Rule struct {
	Goal    Relation
	Pattern RelationSet
}

func (rule Rule) BindSingle(binding Binding) Rule {
	return Rule{
		Goal: rule.Goal.BindSingle(binding),
		Pattern: rule.Pattern.BindSingle(binding),
	}
}

func (rule Rule) Equals(otherRule Rule) bool {
	return rule.Goal.Equals(otherRule.Goal) && rule.Pattern.Equals(otherRule.Pattern)
}

func (rule Rule) Copy() Rule {
	newRule := Rule{}
	newRule.Goal = rule.Goal.Copy()
	newRule.Pattern = rule.Pattern.Copy()
	return newRule
}

func (rule Rule) ConvertVariablesToConstants() Rule {
	newRule := Rule{}
	newRule.Goal = rule.Goal.ConvertVariablesToConstants()
	newRule.Pattern = rule.Pattern.ConvertVariablesToConstants()
	return newRule
}

func (rule Rule) GetVariableNames() []string {
	var names []string

	names = append(names, rule.Goal.GetVariableNames()...)
	names = append(names, rule.Pattern.GetVariableNames()...)

	return common.StringArrayDeduplicate(names)
}

func (rule Rule) String() string {
	s := rule.Goal.String() + " :- " + rule.Pattern.String()
	return s
}