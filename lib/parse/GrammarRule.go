package parse

import "nli-go/lib/mentalese"

type GrammarRule struct {
	SyntacticCategories []string
	EntityVariables     []string
	Sense               mentalese.RelationSet
	PopVariableList		VariableList
	PushVariableList	VariableList
}

func NewGrammarRule(syntacticCategories []string, entityVariables []string, sense mentalese.RelationSet, popVariableList []string, pushvariableList []string) GrammarRule {
	return GrammarRule{
		SyntacticCategories: syntacticCategories,
		EntityVariables:     entityVariables,
		Sense:               sense,
		PopVariableList:	 popVariableList,
		PushVariableList:    pushvariableList,
	}
}

func (rule GrammarRule) GetAntecedent() string {
	return rule.SyntacticCategories[0]
}

func (rule GrammarRule) GetAntecedentVariable() string {
	return rule.EntityVariables[0]
}

func (rule GrammarRule) GetConsequents() []string {
	return rule.SyntacticCategories[1:]
}

func (rule GrammarRule) GetConsequent(i int) string {
	return rule.SyntacticCategories[i+1]
}

func (rule GrammarRule) GetConsequentVariables() []string {
	return rule.EntityVariables[1:]
}

func (rule GrammarRule) GetConsequentVariable(i int) string {
	return rule.EntityVariables[i+1]
}

func (rule GrammarRule) GetConsequentCount() int {
	return len(rule.SyntacticCategories) - 1
}

func (rule GrammarRule) GetPopVariableList() VariableList {
	return rule.PopVariableList
}

func (rule GrammarRule) GetPushVariableList() VariableList {
	return rule.PushVariableList
}

func (rule GrammarRule) Equals(otherRule GrammarRule) bool {

	if len(rule.SyntacticCategories) != len(otherRule.SyntacticCategories) {
		return false
	}

	for i, v := range rule.SyntacticCategories {
		if v != otherRule.SyntacticCategories[i] {
			return false
		}
	}

	if !rule.PopVariableList.Equals(otherRule.PopVariableList) {
		return  false
	}

	if !rule.PushVariableList.Equals(otherRule.PushVariableList) {
		return false
	}

	return true
}

func (rule GrammarRule) GetConsequentIndexByVariable(variable string) (int, bool) {
	for i, entityVariable := range rule.EntityVariables[1:] {
		if entityVariable == variable {
			return i, true
		}
	}

	return 0, false
}

func (rule GrammarRule) Copy() GrammarRule {
	return GrammarRule{
		SyntacticCategories: rule.SyntacticCategories,
		EntityVariables:     rule.EntityVariables,
		Sense:               rule.Sense.Copy(),
		PopVariableList:     rule.PopVariableList.Copy(),
		PushVariableList:    rule.PushVariableList.Copy(),
	}
}

func (rule GrammarRule) String() string {

	s := ""

	s += rule.SyntacticCategories[0] + "(" + rule.EntityVariables[0] + ")"

	if !rule.PopVariableList.Empty() {
		s += " " + rule.PopVariableList.String()
	}

	s += " :- "

	sep := ""
	for i := 1; i < len(rule.SyntacticCategories); i++ {
		s += sep + rule.SyntacticCategories[i] + "(" + rule.EntityVariables[i] + ")"
		sep = " "
	}

	if !rule.PushVariableList.Empty() {
		s += " " + rule.PushVariableList.String()
	}

	s += " { "

	sep = ""
	for _, senseRelation := range rule.Sense {
		s += sep + senseRelation.String()
		sep = ", "
	}

	s += " }"

	return s
}
