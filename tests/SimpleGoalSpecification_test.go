package tests

import (
	"testing"
	"nli-go/lib/importer"
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"nli-go/lib/central"
)

func TestSimpleGoalSpecification(test *testing.T) {

	// relations
	internalGrammarParser := importer.NewSimpleInternalGrammarParser()

	// who did Kurt Cobain marry?
	// non-domain specific
	genericSense, _, _ := internalGrammarParser.CreateRelationSet(`[
		predication(S1, marry)
		tense(S1, past)
		subject(S1, E1)
		object(S1, E2)
		info_request(E1)
		name(E2, 'Kurt Cobain')
	]`)

	// transform the generic sense into a domain specific sense. Leaving out material, but making it more compact
	domainSpecificAnalysis, _, _ := internalGrammarParser.CreateTransformations(`[
		married_to(A, B) :- predication(P1, marry), subject(P1, A), object(P1, B)
		name(A, N) :- name(A, N)
		question(A) :- info_request(A)
	]`)

//	common.LoggerActive=true

	// create domain specific representation
	transformer := mentalese.NewSimpleRelationTransformer(domainSpecificAnalysis)
	domainSpecificSense := transformer.Extract(genericSense)

	common.Logf("DS sense %v\n", domainSpecificSense)

	// goal specification
	// if X was the request, Y is what the user really wants to know
	// turn X into Y
	// name ALL X for whom holds that X was married to Y
	// he / she was married to a, b, and c
	// name(B, N) will have multiple bindings for B and N
	// de operator :- is hier niet handig. Het gaat niet om een implicatie, maar om een interpretatie

	// je loopt er nu al tegenaan dat je makkelijk predicaten uit het verkeerde domein gebruikt. dat moet niet kunnen

	// date er direct een antwoord-template beschikbaar is, is omdat die eerder bedacht is; als er nog geen beschikbaar is,
	// moet die misschien bedacht worden; dat is een meta-probleem

	//domainSpecificGoalAnalysis, _, _ := internalGrammarParser.CreateQAPairs(`[
	//	{
	//		Q: married_to(A, B), info_request(B)
	//		A: married_to(A, B), gender(A, G), name(B, N)
	//	}
	//]`)

	domainSpecificGoalAnalysis, _, _ := internalGrammarParser.CreateTransformations(`[
		married_to(A, B), gender(B, G), name(A, N) :- married_to(A, B), question(A)
	]`)

	// A: married_to(A, B), person(A, _, G, _), person(B, N, _, _)
	// 			RQ: B

	// optimalisatie-fase: doe eerst de meest gerestricteerde doelen (bv name(A, 'Kurt Cobain')

	transformer = mentalese.NewSimpleRelationTransformer(domainSpecificGoalAnalysis)
	goalSense := transformer.Extract(domainSpecificSense)

	common.Logf("Goal sense %v\n", goalSense)

	rules, _, _ := internalGrammarParser.CreateRules(`
	`)
//		married_to(X, Y) :- married_to(Y, X)

	ruleBase1 := knowledge.NewSimpleRuleBase(rules)

	ds2db, _, _ := internalGrammarParser.CreateRules(`[
		married_to(A, B) :- marriages(A, B, _)
		name(A, N) :- person(A, N, _, _)
		gender(A, male) :- person(A, _, 'M', _)
		gender(A, female) :- person(A, _, 'F', _)
	]`)

	// voorbeeld van wanneer dit niet werkt:
	// marriages('Kurt Cobain', 'Courtney Love', '1992')
	// married_to(A, B), name(A, AN), name(B, BN) :- marriages(AN, BN)
	// kan echter wel, door de full names als person ids te beschouwen

	// note! db specific
	facts, _, _ := internalGrammarParser.CreateRelationSet(`[
		marriages(11, 14, '1992')
		person(11, 'Courtney Love', 'F', '1964')
		person(14, 'Kurt Cobain', 'M', '1967')
	]`)
//		marriages(14, 11, '1992')

	factBase1 := knowledge.NewSimpleFactBase(facts, ds2db)

	// produce response
	problemSolver := central.NewSimpleProblemSolver()
	problemSolver.AddKnowledgeBase(factBase1)
	problemSolver.AddKnowledgeBase(ruleBase1)
	domainSpecificResponseSense := problemSolver.Solve(goalSense)

// ERR

// ERR

	//// turn domain specific response into generic response
	//specificResponseSpec, _, _ := internalGrammarParser.CreateTransformations(`[
	//	predication(S1, marry), object(S1, E2), subject(S1, A), object(S1, B) :- marry(A, B)
	//	name(A, N) :- name(A, N)
	//	gender(A, N) :- gender(A, N)
	//]`)
	//transformer = example3.NewSimpleRelationTransformer(specificResponseSpec)
	//genericResponseSense := transformer.Extract(domainSpecificResponseSense)
	//

	if len(domainSpecificResponseSense) == 0 {
		test.Error("Wrong response")
	} else if domainSpecificResponseSense[0].String() != "[married_to(11, 14) gender(14, male) name(11, 'Courtney Love')]" {
		test.Errorf("Wrong response: %s", domainSpecificResponseSense[0].String())
	}

}