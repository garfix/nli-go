package main

import (
)
import (
	"nli-go/lib/importer"
	"fmt"
	"nli-go/lib/parse"
	"nli-go/lib/common"
)

func main() {

	//internalGrammarParser := importer.NewSimpleInternalGrammarParser()
	//grammar, _, _ := internalGrammarParser.CreateGrammar(`[
	//	{
	//		rule: s(P) :- np(E), vp(P)
	//		sense: subject(P, E)
	//	} {
	//		rule: np(E) :- nbar(E)
	//	} {
	//		rule: np(E) :- det(E), nbar(E)
	//	} {
	//		rule: nbar(E) :- noun(E)
	//	} {
	//		rule: nbar(E) :- adj(E), noun(E)
	//	} {
	//		rule: vp(P) :- verb(E)
	//		sense: predicate(P, E)
	//	}
	//]`)
	//
	//lexicon, _, _ := internalGrammarParser.CreateLexicon(`[
	//	{ form: 'the' pos: det }
	//	{ form: 'a' pos: det }
	//	{ form: 'shy' pos: adj }
	//	{ form: 'small' pos: adj }
	//	{ form: 'boy' pos: noun sense: instance_of('*', boy) }
	//	{ form: 'girl' pos: noun sense: instance_of('*', girl) }
	//	{ form: 'cries' pos: verb sense: predicate('*', cry) }
	//	{ form: 'sings' pos: verb sense: predicate('*', sing) }
	//]`)
	//
	//common.LoggerActive = true
	//
	//rawInput := "the small shy girl sings"
	//tokenizer := parse.NewSimpleTokenizer()
	//
	//parser := parse.NewSimpleParser(grammar, lexicon)
	//
	//wordArray := tokenizer.Process(rawInput)
	//
	//_, relations, _ := parser.Process(wordArray)
	//
	//fmt.Printf("%v", relations)

}