package tests

import (
	"testing"
	"nli-go/lib/example3"
)

func TestSimpleLexiconParser(test *testing.T) {

	parser := example3.NewSimpleInternalGrammarParser()
	lexicon := example3.NewSimpleLexicon()
	ok := true

	lexicon, _, ok = parser.CreateLexicon("" +
		"[" +
		"{ form: 'boek'\npos: noun }" +
		"]")
	if !ok {
		test.Error("Parse error")
	}

	_, ok = lexicon.GetLexItem("boek", "noun")
	if !ok {
		test.Error("Parse error")
	}

	lexicon, _, ok = parser.CreateLexicon("" +
		"[" +
		"{ form: 'boek' pos: noun }" +
		"{ form: 'lees' pos: verb }" +
		"]")
	if !ok {
		test.Error("Parse error")
	}
	_, ok = lexicon.GetLexItem("boek", "noun")
	if !ok {
		test.Error("Parse error")
	}
	_, ok = lexicon.GetLexItem("lees", "verb")
	if !ok {
		test.Error("Parse error")
	}
}