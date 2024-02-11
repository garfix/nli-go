package mentalese

import "nli-go/lib/common"

type Rule struct {
	Goal       Relation
	Pattern    RelationSet
	IsFunction bool
}

func (rule Rule) BindSingle(binding Binding) Rule {
	return Rule{
		Goal:    rule.Goal.BindSingle(binding),
		Pattern: rule.Pattern.BindSingle(binding),
	}
}

func (rule Rule) InstantiateUnboundVariables(binding Binding, variableGenerator *VariableGenerator) Rule {
	newRule := Rule{}
	newRule.Goal = rule.Goal
	newRule.Pattern = rule.Pattern.InstantiateUnboundVariables(binding, variableGenerator)
	newRule.IsFunction = rule.IsFunction
	return newRule
}

func (rule Rule) Equals(otherRule Rule) bool {
	return rule.Goal.Equals(otherRule.Goal) && rule.Pattern.Equals(otherRule.Pattern)
}

func (rule Rule) Copy() Rule {
	newRule := Rule{}
	newRule.Goal = rule.Goal.Copy()
	newRule.Pattern = rule.Pattern.Copy()
	newRule.IsFunction = rule.IsFunction
	return newRule
}

func (rule Rule) ConvertVariablesToConstants() Rule {
	newRule := Rule{}
	newRule.Goal = rule.Goal.ConvertVariablesToConstants()
	newRule.Pattern = rule.Pattern.ConvertVariablesToConstants()
	newRule.IsFunction = rule.IsFunction
	return newRule
}

func (rule Rule) ConvertToFunction() Rule {

	newRule := rule.Copy()
	newRule.Pattern = rule.Pattern.ConvertVariablesToMutables()

	// convert the variables in the body
	// assign the arguments to local variables at the beginning of the body
	for _, argument := range rule.Goal.Arguments[0 : len(rule.Goal.Arguments)-rule.Goal.ReturnVariables] {
		assignment := NewRelation(false, PredicateAssign, []Term{
			NewTermVariable(":" + argument.TermValue),
			NewTermVariable(argument.TermValue),
		}, 0)
		newRule.Pattern = append([]Relation{assignment}, newRule.Pattern...)
	}

	returnVariables := []string{}
	for i := 0; i < rule.Goal.ReturnVariables; i++ {
		returnVariables = append(returnVariables, rule.Goal.Arguments[len(rule.Goal.Arguments)-i-1].TermValue)
	}

	// assign the return value to the last argument
	for i := 0; i < len(returnVariables); i++ {
		returnVar := returnVariables[i]
		assignment := NewRelation(false, PredicateAssign, []Term{
			NewTermVariable(returnVar),
			NewTermVariable(":" + returnVar),
		}, 0)
		newRule.Pattern = append(newRule.Pattern, assignment)
	}

	return newRule
}

func (rule Rule) ConvertVariablesToImmutables() Rule {
	returnVar := rule.Goal.Arguments[len(rule.Goal.Arguments)-1].TermValue
	newRule := rule.Copy()
	// convert the variables in the body
	newRule.Pattern = rule.Pattern.ConvertVariablesToImmutables()
	// assign the arguments to local variables at the beginning of the body
	for _, argument := range rule.Goal.Arguments[0 : len(rule.Goal.Arguments)-1] {
		assignment := NewRelation(false, PredicateAssign, []Term{
			NewTermVariable(":" + argument.TermValue),
			NewTermVariable(argument.TermValue),
		}, 0)
		newRule.Pattern = append([]Relation{assignment}, newRule.Pattern...)
	}
	// assign the return value to the last argument
	assignment := NewRelation(false, PredicateAssign, []Term{
		NewTermVariable(returnVar),
		NewTermVariable(":" + returnVar),
	}, 0)
	newRule.Pattern = append(newRule.Pattern, assignment)

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
