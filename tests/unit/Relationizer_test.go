package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
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

	variableGenerator := mentalese.NewVariableGenerator()
	relationizer := parse.NewRelationizer(variableGenerator, log)

	parseTrees := parser.Parse([]string{"the", "book", "falls", "."}, "s", []string{"S"})
	result, _ := relationizer.Relationize(parseTrees[0], []string{ "S"})

	want := "isa(S, fall) subject(S, E1$1) isa(D1$1, the) isa(E1$1, book) determiner(E1$1, D1$1) declaration(S)"
	if result.String() != want {
		t.Errorf("got %s, want %s", result.String(), want)
	}

	result, _ = relationizer.Relationize(parseTrees[0], []string{ "S"})

	want = "isa(S, fall) subject(S, E1$2) isa(D1$2, the) isa(E1$2, book) determiner(E1$2, D1$2) declaration(S)"
	if result.String() != want {
		t.Errorf("got %s, want %s", result.String(), want)
	}

	parseTrees2 := parser.Parse([]string{"the", "book", "falls", "on", "the", "ground", "."}, "s", []string{"S"})
	result2, _ := relationizer.Relationize(parseTrees2[0], []string{ "S"})

	want2 := "isa(S, fall) isa(P1$2, on) isa(D1$4, the) isa(P1$1, ground) determiner(P1$1, D1$4) case(P1$1, P1$2) mod(S, P1$1) subject(S, E1$3) isa(D1$3, the) isa(E1$3, book) determiner(E1$3, D1$3) declaration(S)"
	if result2.String() != want2 {
		t.Errorf("got %s, want %s", result2.String(), want2)
	}
}
