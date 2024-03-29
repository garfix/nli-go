package mentalese

import (
	"nli-go/lib/common"
)

type GrammarRule struct {
	// relation -> relation word-form
	PositionTypes []string
	// vp -> np vbar
	SyntacticCategories []string
	// (P1, E1) -> (E1) (P1, E1)
	EntityVariables [][]string
	Sense           RelationSet
	Ellipsis        CategoryPathList
	Tag             RelationSet
	Intent          RelationSet
}

const PosTypeRelation = "relation"
const PosTypeWordForm = "word-form"
const PosTypeRegExp = "reg-exp"

func NewGrammarRule(positionTypes []string, syntacticCategories []string, entityVariables [][]string, sense RelationSet) GrammarRule {
	return GrammarRule{
		PositionTypes:       positionTypes,
		SyntacticCategories: syntacticCategories,
		EntityVariables:     entityVariables,
		Sense:               sense,
		Ellipsis:            CategoryPathList{},
		Tag:                 RelationSet{},
		Intent:              RelationSet{},
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

func (rule GrammarRule) BindSimple(binding Binding) GrammarRule {
	bound := rule.Copy()

	for i, variables := range rule.EntityVariables {
		for j, variable := range variables {
			val, found := binding.Get(variable)
			if found {
				bound.EntityVariables[i][j] = val.String()
			}
		}
	}

	bound.Sense = bound.Sense.BindSingle(binding)
	bound.Ellipsis = bound.Ellipsis.BindSingle(binding)
	bound.Tag = bound.Tag.BindSingle(binding)
	bound.Intent = bound.Intent.BindSingle(binding)

	return bound
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
		PositionTypes:       common.StringArrayCopy(rule.PositionTypes),
		SyntacticCategories: common.StringArrayCopy(rule.SyntacticCategories),
		EntityVariables:     common.StringMatrixCopy(rule.EntityVariables),
		Sense:               rule.Sense.Copy(),
		Ellipsis:            rule.Ellipsis.Copy(),
		Tag:                 rule.Tag.Copy(),
		Intent:              rule.Intent.Copy(),
	}
}

func (rule GrammarRule) ReplaceVariable(fromVar string, toVar string) GrammarRule {
	newRule := rule.Copy()
	for i, entityVariableArray := range rule.EntityVariables {
		for j, entityVariable := range entityVariableArray {
			if entityVariable == fromVar {
				newRule.EntityVariables[i][j] = toVar
			}
		}
	}
	newRule.Sense = newRule.Sense.ReplaceTerm(NewTermVariable(fromVar), NewTermVariable(toVar))
	newRule.Tag = newRule.Tag.ReplaceTerm(NewTermVariable(fromVar), NewTermVariable(toVar))
	newRule.Intent = newRule.Intent.ReplaceTerm(NewTermVariable(fromVar), NewTermVariable(toVar))
	return newRule
}

func (rule GrammarRule) BasicForm() string {

	s := rule.SyntacticCategories[0] + "("
	sep2 := ""
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

	return s
}

func (rule GrammarRule) String() string {

	s := rule.BasicForm()

	s += " { "

	sep := ""
	for _, senseRelation := range rule.Sense {
		s += sep + senseRelation.String()
		sep = ", "
	}

	s += " }"

	return s
}
