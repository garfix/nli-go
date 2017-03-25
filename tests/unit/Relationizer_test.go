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

	want := "[declaration(S5) subject(S5, E5) quantification(E5, [isa(E5, book)], D5, [isa(D5, the)]) isa(S5, fall)]"
	if result.String() != want {
		t.Errorf("got %s, want %s", result.String(), want)
	}

	result = relationizer.Relationize(parseTree)

	want = "[declaration(S6) subject(S6, E6) quantification(E6, [isa(E6, book)], D6, [isa(D6, the)]) isa(S6, fall)]"
	if result.String() != want {
		t.Errorf("got %s, want %s", result.String(), want)
	}
}