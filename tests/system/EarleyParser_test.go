package tests

import (
	"testing"
	"nli-go/lib/importer"
	"fmt"
	"nli-go/lib/parse"
	"nli-go/lib/parse/earley"
)

func TestEarleyParser(test *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()
	grammar := internalGrammarParser.CreateGrammar(`[
		rule: s(P) -> np(E) vp(P),			sense: subject(P, E);
		rule: np(E) -> nbar(E);
		rule: np(E) -> det(E) nbar(E);
		rule: nbar(E) -> noun(E);
		rule: nbar(E) -> adj(E) nbar(E);
		rule: vp(P) -> verb(P);
	]`)

	lexicon := internalGrammarParser.CreateLexicon(`[
		form: 'the',			pos: det;
		form: 'a',  			pos: det;
		form: 'shy',			pos: adj;
		form: 'small',			pos: adj;
		form: 'boy',			pos: noun,			sense: instance_of(this, boy);
		form: 'girl',			pos: noun,			sense: instance_of(this, girl);
		form: 'cries',  		pos: verb,  		sense: predication(this, cry);
		form: 'sings',			pos: verb,			sense: predication(this, sing);
	]`)

	rawInput := "the small shy girl sings"
	tokenizer := parse.NewTokenizer()

	parser := earley.NewParser(grammar, lexicon)

	wordArray := tokenizer.Process(rawInput)

	relations, tree, _ := parser.Parse(wordArray)

	if relations.String() != "[subject(S1, E1) instance_of(E1, girl) predication(S1, sing)]" {
		test.Error(fmt.Sprintf("Relations: %v", relations))
	}
	if tree.String() != "[s [np [det the] [nbar [adj small] [nbar [adj shy] [nbar [noun girl]]]]] [vp [verb sings]]]" {
		test.Error(fmt.Sprintf("tree: %v", tree.String()))
	}
}
