package tests

import (
	"testing"
	"nli-go/lib/parse"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"nli-go/lib/central"
	"nli-go/lib/knowledge"
)

func TestRelease1(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()
	internalGrammarParser.SetPanicOnParseFail(true)

	// Data

	grammar := internalGrammarParser.CreateGrammar(`[
		{
			rule: s(P) :- np(E), vp(P)
			sense: subject(P, E)
		} {
			rule: np(E) :- nbar(E)
		} {
			rule: np(E) :- det(E), nbar(E)
		} {
			rule: nbar(E) :- noun(E)
		} {
			rule: nbar(E) :- adj(E), nbar(E)
		} {
			rule: vp(P) :- verb(P)
		}
	]`)

	lexicon := internalGrammarParser.CreateLexicon(`[
		{
			form: 'the'
			pos: det
		} {
			form: 'a'
			pos: det
		} {
			form: 'shy'
			pos: adj
		} {
			form: 'small'
			pos: adj
		} {
			form: 'boy'
			pos: noun
			sense: instance_of(this, boy)
		} {
			form: 'girl'
			pos: noun
			sense: instance_of(this, girl)
		} {
			form: 'cries'
			pos: verb
			sense: predication(this, cry)
		} {
			form: 'sings'
			pos: verb
			sense: predication(this, sing)
		}
	]`)

	domainSpecificAnalysis := internalGrammarParser.CreateTransformations(`[
		married_to(A, B) :- predication(P1, marry), subject(P1, A), object(P1, B)
		name(A, N) :- name(A, N)
		question(A) :- info_request(A)
	]`)

	facts := internalGrammarParser.CreateRelationSet(`[
		marriages(11, 14, '1992')
		person(11, 'Courtney Love', 'F', '1964')
		person(14, 'Kurt Cobain', 'M', '1967')
	]`)

	domainSpecificGoalAnalysis := internalGrammarParser.CreateTransformations(`[
		grammatical_subject(B), married_to(A, B), gender(B, G), name(A, N) :- married_to(A, B), question(A)
	]`)

	ds2db := internalGrammarParser.CreateRules(`[
		married_to(A, B) :- marriages(A, B, _)
		name(A, N) :- person(A, N, _, _)
		gender(A, male) :- person(A, _, 'M', _)
		gender(A, female) :- person(A, _, 'F', _)
	]`)

	// Services

	tokenizer := parse.NewTokenizer()
	parser := parse.NewParser(grammar, lexicon)
	transformer := mentalese.NewRelationTransformer()
	factBase1 := knowledge.NewFactBase(facts, ds2db)
	problemSolver := central.NewProblemSolver()
	problemSolver.AddKnowledgeBase(factBase1)

	// Tests

	var tests = []struct {
		question string
		want string
	} {
		{"Who married Jacqueline?", "Marty"},
		{"Did Bob marry Sally?", "Yes"},
		{"Are Jane and Janelle siblings?", "No"},
	}

	for _, test := range tests {

		tokens := tokenizer.Process(test.question)
		genericSense, _, _ := parser.Process(tokens)
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
