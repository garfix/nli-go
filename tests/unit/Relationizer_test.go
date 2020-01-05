package tests

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse/earley"
	"testing"
)

func TestRelationizer(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()

	grammar := internalGrammarParser.CreateGrammar(`[

		{ rule: dp(D1) -> determiner(D1) }
	
		{ rule: np(E1) -> dp(D1) nbar(E1),                                           sense: determiner(E1, D1) }	
		{ rule: nbar(E1) -> noun(E1) }
	
		{ rule: vgp(V1) -> verb(V1) }
		{ rule: vbar(V1) -> vgp(V1) pp(P1),                                          sense: mod(V1, P1) }
		{ rule: vbar(V1) -> vgp(V1) }	
		{ rule: vp(V1) -> vbar(V1) }

		{ rule: pp(E1) -> preposition(P1) np(E1),                                    sense: case(E1, P1) }
	
		{ rule: s_declarative(P1) -> np(E1) vp(P1),									 sense: subject(P1, E1) }
		{ rule: s_declarative(P1) -> s_declarative(P1) period() }
	
		{ rule: s(S1) -> s_declarative(S1),											 sense: declaration(S1) }
	
	]`)

	lexicon := internalGrammarParser.CreateLexicon(`[
		{ form: 'the',        pos: determiner,        sense: isa(E, the) }
		{ form: 'book',       pos: noun,              sense: isa(E, book) }
		{ form: 'falls',      pos: verb,              sense: isa(E, fall) }
		{ form: 'on',         pos: preposition,       sense: isa(E, on) }
	    { form: 'ground',     pos: noun,       		  sense: isa(E, ground) }
		{ form: '.',          pos: period }
	]`)
	log := common.NewSystemLog(false)

	matcher := mentalese.NewRelationMatcher(log)
	dialogContext := central.NewDialogContext(matcher, log)
	solver := central.NewProblemSolver(matcher, dialogContext, log)
	predicates := mentalese.Predicates{}
	nameResolver := central.NewNameResolver(solver, matcher, predicates, log, dialogContext)

	parser := earley.NewParser(grammar, lexicon, nameResolver, predicates, log)

	relationizer := earley.NewRelationizer(lexicon, log)

	parseTree := parser.Parse([]string{"the", "book", "falls", "."})
	result := relationizer.Relationize(parseTree, mentalese.NewKeyCabinet(), nameResolver)

	want := "[declaration(S5) subject(S5, E5) determiner(E5, D5) isa(D5, the) isa(E5, book) isa(S5, fall)]"
	if result.String() != want {
		t.Errorf("got %s, want %s", result.String(), want)
	}

	result = relationizer.Relationize(parseTree, mentalese.NewKeyCabinet(), nameResolver)

	want = "[declaration(S6) subject(S6, E6) determiner(E6, D6) isa(D6, the) isa(E6, book) isa(S6, fall)]"
	if result.String() != want {
		t.Errorf("got %s, want %s", result.String(), want)
	}

	parseTree2 := parser.Parse([]string{"the", "book", "falls", "on", "the", "ground", "."})
	result2 := relationizer.Relationize(parseTree2, mentalese.NewKeyCabinet(), nameResolver)

	want2 := "[declaration(S7) subject(S7, E7) determiner(E7, D7) isa(D7, the) isa(E7, book) mod(S7, P5) isa(S7, fall) case(P5, P6) isa(P6, on) determiner(P5, D8) isa(D8, the) isa(P5, ground)]"
	if result2.String() != want2 {
		t.Errorf("got %s, want %s", result2.String(), want2)
	}
}
