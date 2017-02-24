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
	"nli-go/lib/generate"
	"fmt"
)

func TestRelease1(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()

	// Data

	grammar := internalGrammarParser.LoadGrammar(common.GetCurrentDir() + "/../../resources/english-1.grammar")

	lexicon := internalGrammarParser.CreateLexicon(`[
		form: 'who',        pos: whWord,        sense: isa(E, who);
		form: 'which',      pos: whWord,        sense: isa(E, which);
		form: 'married',    pos: verb, 	        sense: isa(E, marry);
		form: 'did',		pos: auxVerb,       sense: isa(E, do);
		form: 'has',		pos: auxVerb,       sense: isa(E, have);
		form: 'marry',		pos: verb,		    sense: isa(E, marry);
		form: /^[A-Z]/,	    pos: firstName,     sense: name(E, Form, firstName);
		form: 'de',		    pos: insertion,     sense: name(E, 'de', insertion);
		form: 'van',		pos: insertion,     sense: name(E, 'van', insertion);
		form: /^[A-Z]/,	    pos: lastName,      sense: name(E, Form, lastName);
		form: /^[A-Z]/,	    pos: fullName,      sense: name(E, Form, fullName);
		form: 'are',		pos: auxVerb,		sense: isa(E, be);
		form: 'and',		pos: conjunction;
		form: 'siblings',	pos: noun,		    sense: isa(E, sibling);
		form: 'children',	pos: noun,		    sense: isa(E, child);
		form: '?',          pos: questionMark;
	]`)

	generic2ds := internalGrammarParser.CreateTransformations(`[
		isa(P1, marry) subject(P1, A) object(P1, B) => married_to(A, B);
		isa(P1, be) subject(P1, A) conjunction(A, A1, A2) object(P1, B) isa(B, sibling) => siblings(A1, A2);
		isa(P1, have) subject(P1, S) object(P1, O) isa(S, child) => child(S, O);
		name(A, F, firstName) name(A, I, insertion) name(A, L, lastName) join(N, ' ', F, I, L) => name(A, N);
		name(A, N, fullName) => name(A, N);
		question(S, whQuestion) subject(S, E) isa(E, who) => act(question, who);
		question(S, whQuestion) subject(S, E) determiner(E, D) isa(D, which) => act(question, who);
		question(S, yesNoQuestion) => act(question, yesNo);
		focus(E1) => focus(E1);
	]`)

	dsSolutions := internalGrammarParser.CreateSolutions(`[
		condition: act(question, who) married_to(A, B) focus(A),
		preparation: gender(B, G) name(A, N),
		answer: focus(A) married_to(A, B) gender(B, G) name(A, N);

		condition: act(question, yesNo) married_to(A, B),
		preparation: exists(G, A),
		answer: result(G);

		condition: act(question, yesNo) siblings(A, B),
		preparation: exists(G, A),
		answer: result(G);

		condition: act(question, who) child(A, B) focus(A),
		preparation: name(A, N),
		answer: name(A, N) make_and(A, R);
	]`)

	dsInferenceRules := internalGrammarParser.CreateRules(`[
		siblings(A, B) :- parent(C, A) parent(C, B);
	]`)

	ds2db := internalGrammarParser.CreateRules(`[
		married_to(A, B) :- marriages(A, B, _);
		name(A, N) :- person(A, N, _, _);
		parent(P, C) :- parent(P, C);
		child(C, P) :- parent(P, C);
		gender(A, male) :- person(A, _, 'M', _);
		gender(A, female) :- person(A, _, 'F', _);
	]`)

	dbFacts := internalGrammarParser.CreateRelationSet(`[
		marriages(2, 1, '1992')
		parent(4, 2)
		parent(4, 3)
		parent(4, 5)
		parent(4, 6)
		person(1, 'Jacqueline de Boer', 'F', '1964')
		person(2, 'Mark van Dongen', 'M', '1967')
		person(3, 'Suzanne van Dongen', 'F', '1967')
		person(4, 'John van Dongen', 'M', '1938')
		person(5, 'Dirk van Dongen', 'M', '1972')
		person(6, 'Durkje van Dongen', 'M', '1982')
	]`)

	systemFacts := internalGrammarParser.CreateRelationSet(`[
		act(question, who)
		act(question, yesNo)
		focus(A)
	]`)

	ds2system := internalGrammarParser.CreateRules(`[
		act(question, X) :- act(question, X);
		focus(A) :- focus(A);
	]`)

	ds2generic := internalGrammarParser.CreateTransformations(`[
		married_to(A, B) => isa(P1, marry) subject(P1, A) object(P1, B);
		siblings(A1, A2) => isa(P1, be) subject(P1, A) conjunction(A, A1, A2) object(P1, B) isa(B, sibling);
		name(A, N) => name(A, N);
		and(R, A, B) => conjunction(R, A, B) isa(R, and);
		gender(A, female) => isa(A, female);
		result(true) => declaration(S) modifier(S, M) isa(M, yes);
		result(false) => declaration(S) modifier(S, M) isa(M, no);
	]`)

	generationGrammar := internalGrammarParser.CreateGenerationGrammar(`[
        rule: s(P) -> np(E) vp(P),                                                  condition: subject(P, E);
        rule: s(C) -> np(P1) comma(C) s(P2),                                        condition: conjunction(C, P1, P2) conjunction(P2, _, _);
        rule: s(C) -> np(P1) conjunction(C) np(P2),                                 condition: conjunction(C, P1, P2);
        rule: s(P) -> adverb(M),                                                    condition: modifier(P, M);
        rule: vp(V) -> verb(V) np(H),                                               condition: object(V, H);
        rule: np(F) -> proper_noun(F),                                              condition: name(F, Name);
        rule: np(G) -> pronoun(G),                                                  condition: isa(G, female);
	]`)

	generationLexicon := internalGrammarParser.CreateGenerationLexicon(`[
		form: 'married',	pos: verb,		    condition: isa(E, marry);
		form: 'yes',	    pos: adverb,	    condition: isa(E, yes);
		form: 'no',	        pos: adverb,	    condition: isa(E, no);
		form: 'she',	    pos: pronoun,	    condition: subject(P, S) isa(S, female);
		form: 'her',	    pos: pronoun,	    condition: object(S, O) isa(O, female);
		form: '?',	        pos: proper_noun,	condition: name(E, Name);
		form: 'and',	    pos: conjunction,	condition: isa(E, and);
		form: ',',	        pos: comma,         condition: isa(E, and);
	]`)

	// Services

	tokenizer := parse.NewTokenizer()
	parser := earley.NewParser(grammar, lexicon)
	systemFunctionBase := knowledge.NewSystemFunctionBase()
	matcher := mentalese.NewRelationMatcher()
	matcher.AddFunctionBase(systemFunctionBase)
	transformer := mentalese.NewRelationTransformer(matcher)
	factBase1 := knowledge.NewFactBase(dbFacts, ds2db)
	factBase2 := knowledge.NewFactBase(systemFacts, ds2system)
	ruleBase1 := knowledge.NewRuleBase(dsInferenceRules)
	systemPredicateBase := knowledge.NewSystemPredicateBase()
	answerer := central.NewAnswerer(matcher)
	answerer.AddSolutions(dsSolutions)
	answerer.AddFactBase(factBase1)
	answerer.AddFactBase(factBase2)
	answerer.AddRuleBase(ruleBase1)
	answerer.AddMultipleBindingsBase(systemPredicateBase)
	generator := generate.NewGenerator(generationGrammar, generationLexicon)
	surfacer := generate.NewSurfaceRepresentation()

	// Tests

	var tests = []struct {
		question string
		answer   string
	} {
		{"Who married Jacqueline de Boer?", "Mark van Dongen married her"},
		{"Did Mark van Dongen marry Jacqueline de Boer?", "Yes"},
		{"Did Jacqueline de Boer marry Gerard van As?", "No"},
		{"Are Mark van Dongen and Suzanne van Dongen siblings?", "Yes"},
		{"Are Mark van Dongen and John van Dongen siblings?", "No"},
		{"Which children has John van Dongen?", "Mark van Dongen, Suzanne van Dongen, Dirk van Dongen and Durkje van Dongen"},
		//{"How many children has John van Dongen?", "He has 2 children"},
	}

	for _, test := range tests {

		tokens := tokenizer.Process(test.question)
		genericSense, _, _ := parser.Parse(tokens)
		domainSpecificSense := transformer.Extract(generic2ds, genericSense)
common.LoggerActive=false
		dsAnswer := answerer.Answer(domainSpecificSense)
common.LoggerActive=false
		genericAnswer := transformer.Extract(ds2generic, dsAnswer)
		answerWords := generator.Generate(genericAnswer)
		answer := surfacer.Create(answerWords)

		fmt.Println()
		//fmt.Println(genericSense)
		//fmt.Println(domainSpecificSense)
		//fmt.Println(dsAnswer)
		//fmt.Println(genericAnswer)

		if answer != test.answer {
			t.Errorf("release1: got %v, want %v", answer, test.answer)
		}
	}
}
