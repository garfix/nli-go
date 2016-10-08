package tests

import (
	"testing"
	"nli-go/lib/importer"
	"fmt"
	"nli-go/lib/parse"
)

func TestSimpleParser(test *testing.T) {

	internalGrammarParser := importer.NewSimpleInternalGrammarParser()
	grammar, _, _ := internalGrammarParser.CreateGrammar(`[
		{
			rule: s(P) :- np(E), vp(P)
			sense: subject(P, E)
		} {
			rule: np(E) :- nbar(E)
		} {
			rule: np(E) :- det(E), nbar(E)
		} {
			rule: nbar(E) :- noun(E)
		} {
			rule: nbar(E) :- adj(E), nbar(E)
		} {
			rule: vp(P) :- verb(P)
		}
	]`)

	lexicon, _, _ := internalGrammarParser.CreateLexicon(`[
		{
			form: 'the'
			pos: det
		} {
			form: 'a'
			pos: det
		} {
			form: 'shy'
			pos: adj
		} {
			form: 'small'
			pos: adj
		} {
			form: 'boy'
			pos: noun
			sense: instance_of('*', boy)
		} {
			form: 'girl'
			pos: noun
			sense: instance_of('*', girl)
		} {
			form: 'cries'
			pos: verb
			sense: predication('*', cry)
		} {
			form: 'sings'
			pos: verb
			sense: predication('*', sing)
		}
	]`)

	rawInput := "the small shy girl sings"
	tokenizer := parse.NewSimpleTokenizer()

	parser := parse.NewSimpleParser(grammar, lexicon)

	wordArray := tokenizer.Process(rawInput)

	length, relations, ok := parser.Process(wordArray)

	if !ok {
		test.Errorf("Parse failed at pos %d", length)
	}
	if relations.String() != "[subject(S1, E1) instance_of(E1, girl) predication(S1, sing)]" {
		test.Error(fmt.Sprintf("Relations: %v", relations))
	}
}
