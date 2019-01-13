package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/parse/earley"
	"testing"
)

func TestRelationizer(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()

	grammar := internalGrammarParser.CreateGrammar(internalGrammarParser.LoadText("../../resources/relationships/english-1.grammar"))
	lexicon := internalGrammarParser.CreateLexicon(`[
		{ form: 'the',        pos: determiner,        sense: isa(E, the) }
		{ form: 'book',       pos: noun,              sense: isa(E, book) }
		{ form: 'falls',      pos: verb,              sense: isa(E, fall) }
		{ form: 'on',         pos: preposition,       sense: isa(E, on) }
	    { form: 'ground',     pos: noun,       		  sense: isa(E, ground) }
		{ form: '.',          pos: period }
	]`)
	log := common.NewSystemLog(false)
	parser := earley.NewParser(grammar, lexicon, log)
	relationizer := earley.NewRelationizer(lexicon, log)

	parseTree := parser.Parse([]string{"the", "book", "falls", "."})
	result := relationizer.Relationize(parseTree)

	want := "[declaration(S5) subject(S5, E5) determiner(E5, D5) isa(D5, the) isa(E5, book) isa(S5, fall)]"
	if result.String() != want {
		t.Errorf("got %s, want %s", result.String(), want)
	}

	result = relationizer.Relationize(parseTree)

	want = "[declaration(S6) subject(S6, E6) determiner(E6, D6) isa(D6, the) isa(E6, book) isa(S6, fall)]"
	if result.String() != want {
		t.Errorf("got %s, want %s", result.String(), want)
	}

	parseTree2 := parser.Parse([]string{"the", "book", "falls", "on", "the", "ground", "."})
	result2 := relationizer.Relationize(parseTree2)

	want2 := "[declaration(S7) subject(S7, E7) determiner(E7, D7) isa(D7, the) isa(E7, book) isa(S7, fall) case(P5, E8) isa(P5, on) determiner(E8, D8) isa(D8, the) isa(E8, ground)]"
	if result2.String() != want2 {
		t.Errorf("got %s, want %s", result2.String(), want2)
	}
}
