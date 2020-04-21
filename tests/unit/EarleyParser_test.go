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
	grammar := internalGrammarParser.CreateGrammar(`[
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
	]`)

	log := common.NewSystemLog(false)

	rawInput := "the small shy girl speaks up"
	tokenizer := parse.NewTokenizer(log)

	matcher := mentalese.NewRelationMatcher(log)
	dialogContext := central.NewDialogContext()
	predicates := mentalese.Predicates{}
	solver := central.NewProblemSolver(matcher, predicates, dialogContext, log)
	nameResolver := central.NewNameResolver(solver, matcher, predicates, log, dialogContext)

	parser := earley.NewParser(grammar, nameResolver, predicates, log)
	relationizer := earley.NewRelationizer(log)

	wordArray := tokenizer.Process(rawInput)

	trees := parser.Parse(wordArray)
	relations, _ := relationizer.Relationize(trees[0], nameResolver)

	if relations.String() != "[subject(S5, E5) determiner(E5, D5) isa(D5, the) isa(E5, girl) predication(S5, speak_up)]" {
		test.Error(fmt.Sprintf("Relations: %v", relations))
	}
	if trees[0].String() != "[s [np [det [the the]] [nbar [adj [small small]] [nbar [adj [shy shy]] [nbar [noun [girl girl]]]]]] [vp [verb [speaks speaks] [up up]]]]" {
		test.Error(fmt.Sprintf("tree: %v", trees[0].String()))
	}
}

