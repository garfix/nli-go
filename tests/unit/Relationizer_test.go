package tests

import (
"testing"
	"nli-go/lib/parse/earley"
	"nli-go/lib/importer"
	"nli-go/lib/common"
)

func TestRelationizer(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()

	grammar := internalGrammarParser.LoadGrammar(common.GetCurrentDir() + "/../../resources/english-1.grammar")
	lexicon := internalGrammarParser.CreateLexicon(`[
		form: 'the',        pos: determiner,        sense: isa(E, the);
		form: 'book',       pos: noun,              sense: isa(E, book);
		form: 'falls',      pos: verb,              sense: isa(E, fall);
		form: '.',          pos: period;
	]`)
	parser := earley.NewParser(grammar, lexicon)
	relationizer := earley.NewRelationizer(lexicon)

	parseTree, _ := parser.Parse([]string{"the", "book", "falls", "."})
	result := relationizer.Relationize(parseTree)

	want := "[declaration(S1) subject(S1, E1) quantification(E1, [isa(E1, book)], D1, [isa(D1, the)]) isa(S1, fall)]"
	if result.String() != want {
		t.Errorf("got %s, want %s", result.String(), want)
	}

	result = relationizer.Relationize(parseTree)

	want = "[declaration(S2) subject(S2, E2) quantification(E2, [isa(E2, book)], D2, [isa(D2, the)]) isa(S2, fall)]"
	if result.String() != want {
		t.Errorf("got %s, want %s", result.String(), want)
	}
}