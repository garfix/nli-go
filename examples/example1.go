package main

import (
	"fmt"
	"nli-go/lib/importer"
)

func main() {
	parser := importer.NewInternalGrammarParser()

	grammar := parser.CreateGrammar("[" +
		"rule: s(P) -> np(E) vp(P),         sense: subject(P, E)" +
		"]")

	rules := grammar.FindRules("s")
	if len(rules) == 0 {
		fmt.Print("No rules found")
	}

	if rules[0].SyntacticCategories[0] != "s" {
		fmt.Printf("Error in rule: %s", rules[0].SyntacticCategories[0])
	}
}