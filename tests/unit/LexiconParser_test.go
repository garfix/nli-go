package tests

import (
	"nli-go/lib/importer"
	"testing"
)

func TestLexiconParser(test *testing.T) {

	parser := importer.NewInternalGrammarParser()
	ok := true

	lexicon := parser.CreateLexicon(`[
		{ form: 'boek', pos: noun }
	]`)

	_, ok, _ = lexicon.GetLexItem("boek", "noun")
	if !ok {
		test.Error("Parse error")
	}

	lexicon = parser.CreateLexicon(`[
		{ form: 'boek',   pos: noun }
		{ form: 'lees',   pos: verb }
	]`)

	_, ok, _ = lexicon.GetLexItem("boek", "noun")
	if !ok {
		test.Error("Parse error")
	}
	_, ok, _ = lexicon.GetLexItem("lees", "verb")
	if !ok {
		test.Error("Parse error")
	}
}
