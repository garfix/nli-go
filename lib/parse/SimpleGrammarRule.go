package parse

import "nli-go/lib/mentalese"

type SimpleGrammarRule struct {
	SyntacticCategories []string
	EntityVariables     []string
	Sense               []mentalese.SimpleRelation
}

func (rule SimpleGrammarRule) String() string {

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