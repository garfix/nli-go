package tests

import (
	"nli-go/lib/importer"
	"testing"
)

func TestLexiconParser(test *testing.T) {

	parser := importer.NewInternalGrammarParser()
	ok := true

	lexicon := parser.CreateLexicon(`[
		form: 'boek', pos: noun;
	]`)

	_, ok = lexicon.GetLexItem("boek", "noun")
	if !ok {
		test.Error("Parse error")
	}

	lexicon = parser.CreateLexicon(`[
		form: 'boek',   pos: noun;
		form: 'lees',   pos: verb;
	]`)

	_, ok = lexicon.GetLexItem("boek", "noun")
	if !ok {
		test.Error("Parse error")
	}
	_, ok = lexicon.GetLexItem("lees", "verb")
	if !ok {
		test.Error("Parse error")
	}
}
