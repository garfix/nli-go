package tests

import (
	"testing"
	"fmt"
	"nli-go/lib/importer"
)

func TestSimpleGrammar(test *testing.T) {

	parser := importer.NewSimpleInternalGrammarParser()
	ok := true

	grammar, _, ok := parser.CreateGrammar("[" +
		"{" +
		"rule: s(P) :- np(E), vp(P)" +
		"sense: subject(P, E)" +
		"}" +
		"]")

	if !ok {
		test.Error("Parse error")
	}

	rules := grammar.FindRules("s")
	if len(rules) == 0 {
		test.Error("No rules found")
	}

	if rules[0].SyntacticCategories[0] != "s" {
		test.Error(fmt.Printf("Error in rule: %s", rules[0].SyntacticCategories[0]))
	}
	if rules[0].SyntacticCategories[1] != "np" {
		test.Error(fmt.Printf("Error in rule: %s", rules[0].SyntacticCategories[1]))
	}
	if rules[0].EntityVariables[0] != "P" {
		test.Error(fmt.Printf("Error in rule: %s", rules[0].EntityVariables[0]))
	}
	if rules[0].EntityVariables[1] != "E" {
		test.Error(fmt.Printf("Error in rule: %s", rules[0].EntityVariables[1]))
	}
	if len(rules[0].Sense) != 1 {
		test.Error(fmt.Printf("Error in number of sense relations: %s", len(rules[0].Sense)))
	}

	grammar, _, ok = parser.CreateGrammar("[" +
		"{" +
		"rule: s(P) :- np(E), vp(P)" +
		"sense: subject(P, E)" +
		"}" +
		"{" +
		"rule: np(P) :- nbar(E)" +
		"}" +
		"]")

	rules = grammar.FindRules("s")
	if len(rules) != 1 {
		test.Error("No rules found")
	}
	rules = grammar.FindRules("np")
	if len(rules) != 1 {
		test.Error("No rules found")
	}

	grammar, _, ok = parser.CreateGrammar("[]")
	if !ok {
		test.Error("Parse error")
	}
}