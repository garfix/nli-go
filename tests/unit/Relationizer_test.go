package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/parse"
	"testing"
)

func TestRelationizer(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()

	grammarRules := internalGrammarParser.CreateGrammarRules(`

		{ rule: dp(D1) -> determiner(D1) }
	
		{ rule: np(E1) -> dp(D1) nbar(E1),                                           sense: determiner(E1, D1) }	
		{ rule: nbar(E1) -> noun(E1) }
	
		{ rule: vgp(V1) -> verb(V1) }
		{ rule: vbar(V1) -> vgp(V1) pp(P1),                                          sense: mod(V1, P1) }
		{ rule: vbar(V1) -> vgp(V1) }	
		{ rule: vp(V1) -> vbar(V1) }

		{ rule: pp(E1) -> preposition(P1) np(E1),                                    sense: case(E1, P1) }
	
		{ rule: s_declarative(P1) -> np(E1) vp(P1),									 sense: subject(P1, E1) $np } /* test explicit child semantics append */
		{ rule: s_declarative(P1) -> s_declarative(P1) '.' }
	
		{ rule: s(S1) -> s_declarative(S1),											 sense: declaration(S1) }

		{ rule: determiner(Q1) -> 'the', sense: isa(Q1, the) }
		{ rule: noun(E1) -> 'book', sense: isa(E1, book) }
		{ rule: verb(P1) -> 'falls', sense: isa(P1, fall) }
		{ rule: preposition(P1) -> 'on', sense: isa(P1, on) }
	    { rule: noun(E1) -> 'ground', sense: isa(E1, ground) }
	
	`)

	log := common.NewSystemLog()
	parser := parse.NewParser(grammarRules, log)

	relationizer := parse.NewRelationizer(log)

	parseTrees := parser.Parse([]string{"the", "book", "falls", "."})
	result, _ := relationizer.Relationize(parseTrees[0])

	want := "isa(S5, fall) subject(S5, E5) isa(D5, the) isa(E5, book) determiner(E5, D5) declaration(S5)"
	if result.String() != want {
		t.Errorf("got %s, want %s", result.String(), want)
	}

	result, _ = relationizer.Relationize(parseTrees[0])

	want = "isa(S6, fall) subject(S6, E6) isa(D6, the) isa(E6, book) determiner(E6, D6) declaration(S6)"
	if result.String() != want {
		t.Errorf("got %s, want %s", result.String(), want)
	}

	parseTrees2 := parser.Parse([]string{"the", "book", "falls", "on", "the", "ground", "."})
	result2, _ := relationizer.Relationize(parseTrees2[0])

	want2 := "isa(S7, fall) isa(P6, on) isa(D8, the) isa(P5, ground) determiner(P5, D8) case(P5, P6) mod(S7, P5) subject(S7, E7) isa(D7, the) isa(E7, book) determiner(E7, D7) declaration(S7)"
	if result2.String() != want2 {
		t.Errorf("got %s, want %s", result2.String(), want2)
	}
}
