package natlang

import "nli-go/lib/mentalese"

type SimpleGrammarRule struct {
	SyntacticCategories []string
	EntityVariables     []string
	Sense               []mentalese.SimpleRelation
}
