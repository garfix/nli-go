package tests

import (
	"fmt"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"nli-go/lib/parse/earley"
	"testing"
)

func TestEarleyParser(test *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()
	grammarRules := internalGrammarParser.CreateGrammarRules(`
		{ rule: s(P) -> np(E) vp(P),			sense: subject(P, E) }
		{ rule: np(E) -> nbar(E) }
		{ rule: np(E) -> det(D) nbar(E),      sense: determiner(E, D) }
		{ rule: nbar(E) -> noun(E) }
		{ rule: nbar(E) -> adj(E) nbar(E) }
		{ rule: vp(P) -> verb(P) }
		{ rule: det(E1) -> 'the', sense: isa(E1, the) }
		{ rule: det(E1) -> 'a' }
		{ rule: adj(E1) -> 'shy' }
		{ rule: adj(E1) -> 'small' }
		{ rule: noun(E1) -> 'boy', sense: isa(E1, boy) }
		{ rule: noun(E1) -> 'girl', sense: isa(E1, girl) }
		{ rule: verb(P1) -> 'cries', sense: predication(P1, cry) }
		{ rule: verb(P1) -> 'speaks' 'up', sense: predication(P1, speak_up) }

		{ rule: s(P) -> first(P) second(P) }
		{ rule: first(P) -> early(P) middle(P) middle(P) }
		{ rule: first(P) -> early(P) middle(P) }
		{ rule: first(P) -> early(P) }
		{ rule: second(P) -> middle(P) middle(P) last(P) }
		{ rule: second(P) -> middle(P) last(P) }
		{ rule: second(P) -> last(P) }
		{ rule: early(P) -> 'a' }
		{ rule: middle(P) -> 'b' }
		{ rule: last(P) -> 'c' }
	`)

	log := common.NewSystemLog()

	rawInput := "the small shy girl speaks up"
	tokenizer := parse.NewTokenizer(parse.DefaultTokenizerExpression)

	matcher := central.NewRelationMatcher(log)
	dialogContext := central.NewDialogContext()
	meta := mentalese.NewMeta()
	solver := central.NewProblemSolver(matcher, dialogContext, log)
	nameResolver := central.NewNameResolver(solver, meta, log, dialogContext)

	parser := earley.NewParser(nameResolver, meta, log)
	relationizer := earley.NewRelationizer(log)

	{
		wordArray := tokenizer.Process(rawInput)

		trees := parser.Parse(grammarRules, wordArray)

		if len(trees) != 1 {
			test.Error(fmt.Sprintf("expected : 1 tree, found %d", len(trees)))
			return
		}

		relations, _ := relationizer.Relationize(trees[0])

		if relations.String() != "isa(D5, the) isa(E5, girl) determiner(E5, D5) predication(S5, speak_up) subject(S5, E5)" {
			test.Error(fmt.Sprintf("Relations: %v", relations))
		}
		if trees[0].String() != "[s [np [det [the the]] [nbar [adj [small small]] [nbar [adj [shy shy]] [nbar [noun [girl girl]]]]]] [vp [verb [speaks speaks] [up up]]]]" {
			test.Error(fmt.Sprintf("tree: %v", trees[0].String()))
		}
	}

	{
		wordArray := tokenizer.Process("a b b c")

		trees := parser.Parse(grammarRules, wordArray)

		if len(trees) != 3 {
			test.Error(fmt.Sprintf("expected : 3 trees, found %d", len(trees)))
			return
		}

		expected := []string{
			"[s [first [early [a a]]] [second [middle [b b]] [middle [b b]] [last [c c]]]]",
			"[s [first [early [a a]] [middle [b b]]] [second [middle [b b]] [last [c c]]]]",
			"[s [first [early [a a]] [middle [b b]] [middle [b b]]] [second [last [c c]]]]",
		}

		for i, exp := range expected {
			if trees[i].String() != exp {
				test.Error(fmt.Sprintf("ERR tree %d: %v", i, trees[i].String()))
			}
		}
	}
}

