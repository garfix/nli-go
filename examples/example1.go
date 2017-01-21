package main

import (
	"fmt"
	"nli-go/lib/importer"
	"nli-go/lib/parse/earley"
	"nli-go/lib/parse"
	"nli-go/lib/common"
)

func main() {
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

	common.LoggerActive = true

	_, tree, _ := parser.Parse(wordArray)

	if tree.String() != "[s [np [det the] [nbar [adj small] [nbar [adj shy] [nbar [noun girl]]]]] [vp [verb sings]]]" {
		fmt.Printf("Tree: %v", tree.String())
	}

	//if relations.String() != "[subject(S1, E1) instance_of(E1, girl) predication(S1, sing)]" {
	//	fmt.Printf("Relations: %v", relations)
	//}

}