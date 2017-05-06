package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/parse/earley"
	"testing"
)

func TestRelationizer(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()

	grammar := internalGrammarParser.CreateGrammar(internalGrammarParser.LoadText("../../resources/english-1.grammar"))
	lexicon := internalGrammarParser.CreateLexicon(`[
		form: 'the',        pos: determiner,        sense: isa(E, the);
		form: 'book',       pos: noun,              sense: isa(E, book);
		form: 'falls',      pos: verb,              sense: isa(E, fall);
		form: '.',          pos: period;
	]`)
	log := common.NewSystemLog(false)
	parser := earley.NewParser(grammar, lexicon, log)
	relationizer := earley.NewRelationizer(lexicon, log)

	parseTree := parser.Parse([]string{"the", "book", "falls", "."})
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
