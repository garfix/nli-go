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
	grammar := internalGrammarParser.LoadGrammar(common.GetCurrentDir() + "/../resources/english-1.grammar")

	lexicon := internalGrammarParser.CreateLexicon(`[
		form: 'the',			pos: det,            sense: isa(this, the);
		form: 'a',  			pos: det;
		form: 'shy',			pos: adj;
		form: 'small',			pos: adj;
		form: 'boy',			pos: noun,			sense: isa(this, boy);
		form: 'girl',			pos: noun,			sense: isa(this, girl);
		form: 'cries',  		pos: verb,  		sense: predication(this, cry);
		form: 'sings',			pos: verb,			sense: predication(this, sing);
	]`)

	rawInput := "the small shy girl sings"
	tokenizer := parse.NewTokenizer()

	parser := earley.NewParser(grammar, lexicon)

	wordArray := tokenizer.Process(rawInput)

	relations, tree, _ := parser.Parse(wordArray)

	if tree.String() != "[s [np [det the] [nbar [adj small] [nbar [adj shy] [nbar [noun girl]]]]] [vp [verb sings]]]" {
		fmt.Printf("Tree: %v", tree.String())
	}

	if relations.String() != "[subject(S1, E1) determiner(E1, D1) isa(D1, the) isa(E1, girl) predication(S1, sing)]" {
		fmt.Printf("Relations: %v", relations)
	}

}