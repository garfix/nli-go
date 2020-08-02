package parse

import "nli-go/lib/mentalese"

type GrammarRule struct {
	// relation -> relation word-form
	PositionTypes 		[]string
	// vp -> np vbar
	SyntacticCategories []string
	// (P1, E1) -> (E1) (P1, E1)
	EntityVariables     [][]string
	Sense               mentalese.RelationSet
}

const PosTypeRelation = "relation"
const PosTypeWordForm = "word-form"
const PosTypeRegExp = "reg-exp"

func NewGrammarRule(positionTypes []string, syntacticCategories []string, entityVariables [][]string, sense mentalese.RelationSet) GrammarRule {
	return GrammarRule{
		PositionTypes: 		 positionTypes,
		SyntacticCategories: syntacticCategories,
		EntityVariables:     entityVariables,
		Sense:               sense,
	}
}

func (rule GrammarRule) GetAntecedent() string {
	return rule.SyntacticCategories[0]
}

func (rule GrammarRule) GetAntecedentVariables() []string {
	return rule.EntityVariables[0]
}

func (rule GrammarRule) GetConsequents() []string {
	return rule.SyntacticCategories[1:]
}

func (rule GrammarRule) GetConsequent(i int) string {
	return rule.SyntacticCategories[i+1]
}

func (rule GrammarRule) GetAllConsequentVariables() [][]string {
	return rule.EntityVariables[1:]
}

func (rule GrammarRule) GetConsequentVariables(i int) []string {
	return rule.EntityVariables[i+1]
}

func (rule GrammarRule) GetConsequentPositionType(i int) string {
	return rule.PositionTypes[i+1]
}

func (rule GrammarRule) GetConsequentCount() int {
	return len(rule.SyntacticCategories) - 1
}

// returns the index of the consequent specified by, for example: np2
func (rule GrammarRule) FindConsequentIndex(category string, catIndex int) int {

	count := 0

	for index, cat := range rule.GetConsequents() {
		if cat == category {
			count++
			if count == catIndex {
				return index
			}
		}
	}

	return -1
}

func (rule GrammarRule) Equals(otherRule GrammarRule) bool {

	if len(rule.SyntacticCategories) != len(otherRule.SyntacticCategories) {
		return false
	}

	for i, v := range rule.SyntacticCategories {
		if v != otherRule.SyntacticCategories[i] {
			return false
		}
		if rule.PositionTypes[i] != otherRule.PositionTypes[i] {
			return false
		}
		if len(rule.EntityVariables[i]) != len(otherRule.EntityVariables[i]) {
			return false
		}
		for j, w := range rule.EntityVariables[i] {
			if w != otherRule.EntityVariables[i][j] {
				return false
			}
		}
	}

	return true
}

func (rule GrammarRule) Copy() GrammarRule {
	return GrammarRule{
		PositionTypes:		 rule.PositionTypes,
		SyntacticCategories: rule.SyntacticCategories,
		EntityVariables:     rule.EntityVariables,
		Sense:               rule.Sense.Copy(),
	}
}

func (rule GrammarRule) String() string {

	s := ""
	sep2 := ""

	s += rule.SyntacticCategories[0] + "("
	sep2 = ""
	for _, variable := range rule.EntityVariables[0] {
		s += sep2 + variable
		sep2 = ", "
	}
	s += ")"

	s += " -> "

	sep := ""
	for i := 1; i < len(rule.SyntacticCategories); i++ {
		if rule.PositionTypes[i] == PosTypeRelation {
			s += sep + rule.SyntacticCategories[i] + "("
			sep2 = ""
			for _, variable := range rule.EntityVariables[i] {
				s += sep2 + variable
				sep2 = ", "
			}
			s += ")"
		} else if rule.PositionTypes[i] == PosTypeWordForm {
			s += sep + "'" + rule.SyntacticCategories[i] + "'"
		} else {
			s += sep + "/" + rule.SyntacticCategories[i] + "/"
		}
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
