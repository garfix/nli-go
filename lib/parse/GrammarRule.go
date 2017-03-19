package parse

import "nli-go/lib/mentalese"

type GrammarRule struct {
	SyntacticCategories []string
	EntityVariables     []string
	Sense               mentalese.RelationSet
}

func NewGrammarRule(syntacticCategories []string, entityVariables []string, sense mentalese.RelationSet) GrammarRule {
	return GrammarRule{
		SyntacticCategories: syntacticCategories,
		EntityVariables: entityVariables,
		Sense: sense,
	}
}

func (rule GrammarRule) GetAntecedent() string {
	return rule.SyntacticCategories[0]
}

func (rule GrammarRule) GetConsequents() []string {
	return rule.SyntacticCategories[1:]
}

func (rule GrammarRule) GetConsequent(i int) string {
	return rule.SyntacticCategories[i + 1]
}

func (rule GrammarRule) GetConsequentCount() int {
	return len(rule.SyntacticCategories) - 1
}

func (rule GrammarRule) Equals(otherRule GrammarRule) bool {

	if len(rule.SyntacticCategories) != len(otherRule.SyntacticCategories) {
		return false;
	}

	for i, v := range rule.SyntacticCategories {
		if v != otherRule.SyntacticCategories[i] { return false }
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

func (rule GrammarRule) String() string {

	s := ""

	s += rule.SyntacticCategories[0] + "(" + rule.EntityVariables[0] + ")"

	s += " :- "

	sep := ""
	for i := 1; i < len(rule.SyntacticCategories); i++  {
		s += sep + rule.SyntacticCategories[i] + "(" + rule.EntityVariables[i] + ")"
		sep = " "
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