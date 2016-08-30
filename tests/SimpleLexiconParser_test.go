package tests

import (
	"testing"
	"nli-go/lib/example3"
)

func TestSimpleLexiconParser(test *testing.T) {

	parser := example3.NewSimpleInternalGrammarParser()
	ok := true

	lexicon, _, ok := parser.CreateLexicon("" +
		"[" +
		"\t{ form: 'boek'\npos: noun }" +
		"]")
	if !ok {
		test.Error("Parse error")
	}

	_, ok = lexicon.GetLexItem("boek", "noun")
	if !ok {
		test.Error("Parse error")
	}
}