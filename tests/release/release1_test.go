package tests

import (
	"testing"
	"nli-go/lib/parse"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"nli-go/lib/central"
	"nli-go/lib/knowledge"
	"nli-go/lib/common"
	"nli-go/lib/parse/earley"
	"fmt"
)

func TestRelease1(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()
	internalGrammarParser.SetPanicOnParseFail(true)

	// Data

	grammar := internalGrammarParser.LoadGrammar(common.GetCurrentDir() + "/../../resources/english-1.grammar")

	lexicon := internalGrammarParser.CreateLexicon(`[
		form: 'who',        pos: whWord;
		form: 'married',    pos: verb, 	        sense: isa(this, marry);
		form: 'did',		pos: auxDo;
		form: 'marry',		pos: verb,		    sense: isa(this, marry);
		form: 'de',		    pos: insertion      sense: name(this, 'de', insertion);
		form: '[A-Z].*',	pos: lastName       sense: name(this, form, lastName);
		form: '[A-Z].*',	pos: firstName      sense: name(this, form, firstName);
		form: 'are',		pos: auxBe,		    sense: isa(this, be);
		form: 'and',		pos: conjunction;
		form: 'siblings',	pos: noun,		    sense: isa(this, sibling);
		form: '?'           pos: questionMark;
	]`)

	domainSpecificAnalysis := internalGrammarParser.CreateTransformations(`[
		married_to(A, B) :- predication(P1, marry) subject(P1, A) object(P1, B);
		name(A, N) :- name(A, N);
		question(A) :- info_request(A);
	]`)

	facts := internalGrammarParser.CreateRelationSet(`[
		marriages(11, 14, '1992')
		person(11, 'Courtney Love', 'F', '1964')
		person(14, 'Kurt Cobain', 'M', '1967')
	]`)

	domainSpecificGoalAnalysis := internalGrammarParser.CreateTransformations(`[
		grammatical_subject(B) married_to(A, B) gender(B, G) name(A, N) :- married_to(A, B) question(A);
	]`)

	ds2db := internalGrammarParser.CreateRules(`[
		married_to(A, B) :- marriages(A, B, _);
		name(A, N) :- person(A, N, _, _);
		gender(A, male) :- person(A, _, 'M', _);
		gender(A, female) :- person(A, _, 'F', _);
	]`)

	// Services

	tokenizer := parse.NewTokenizer()
	parser := earley.NewParser(grammar, lexicon)
	transformer := mentalese.NewRelationTransformer()
	factBase1 := knowledge.NewFactBase(facts, ds2db)
	problemSolver := central.NewProblemSolver()
	problemSolver.AddKnowledgeBase(factBase1)

	// Tests

	var tests = []struct {
		question string
		want string
	} {
		{"Who married Jacqueline de Boer?", "Marty"},
		//{"Did Bob marry Sally?", "Yes"},
		//{"Are Jane and Janelle siblings?", "No"},
	}

	for _, test := range tests {

		tokens := tokenizer.Process(test.question)

common.LoggerActive=false
		genericSense, _, _ := parser.Parse(tokens)
common.LoggerActive=false
fmt.Print(genericSense)
		domainSpecificSense := transformer.Extract(domainSpecificAnalysis, genericSense)
		goalSense := transformer.Extract(domainSpecificGoalAnalysis, domainSpecificSense)
		//domainSpecificResponseSenses :=
			problemSolver.Solve(goalSense)

		answer := ""

		if answer != test.want {
//			t.Errorf("release1: got %v, want %v", answer, test.want)
		}
	}
}
