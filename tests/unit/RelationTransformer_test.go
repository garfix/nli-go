package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"testing"
)

func TestRelationTransformer(t *testing.T) {

	log := common.NewSystemLog(false)
	parser := importer.NewInternalGrammarParser()
	matcher := mentalese.NewRelationMatcher(log)
	transformer := mentalese.NewRelationTransformer(matcher, log)

	// "name all customers"
	relationSet := parser.CreateRelationSet(`
		instance_of(E2, name)
		predicate(S1, name)
		object(S1, E1)
		instance_of(E1, customer)
		instance_of(D1, all)
	`)

	tests := []struct {
		transformations string
		wantReplaced    string
	}{
		{
			`[
				instance_of(Z, B) :- task(A, B);
				predicate(A, X) :- magic(A, X);
				predicate(A, X) :- label(A, O);
			]`,
			"task(A, name) magic(S1, name) label(S1, O) object(S1, E1) task(A, customer) task(A, all)",
		},
		{
			`[
				instance_of(Z, B) :- isa(Z, B);
			]`,
			"isa(E2, name) predicate(S1, name) object(S1, E1) isa(E1, customer) isa(D1, all)",
		},
		{
			`[
				instance_of(Z, B) :- first(Z, C) second(C, B);
			]`,
			"first(E2, C) second(C, name) predicate(S1, name) object(S1, E1) first(E1, C) second(C, customer) first(D1, C) second(C, all)",
		},
	}
	for _, test := range tests {

		transformations := parser.CreateRules(test.transformations)

		wantReplaced := parser.CreateRelationSet(test.wantReplaced)

		replacedResult := transformer.Replace(transformations, relationSet)

		if replacedResult.String() != wantReplaced.String() {
			t.Errorf("RelationTransformer: got\n%v,\nwant\n%v", replacedResult, wantReplaced)
		}
	}
}

func TestRelationTransformerWithRelationSetArguments(t *testing.T) {

	log := common.NewSystemLog(false)
	parser := importer.NewInternalGrammarParser()
	matcher := mentalese.NewRelationMatcher(log)
	transformer := mentalese.NewRelationTransformer(matcher, log)

	relationSet := parser.CreateRelationSet(`
		quant(E1, isa(E1, ball), D1, isa(D1, every), none)
	`)

	transformations := parser.CreateRules(`[
		quant(E2, isa(E2, ball), D2, isa(D2, every), none) :- quant(E2, isa(E2, ball), D2, isa(D2, every), none) ok(E2, D2);
	]`)

	replacedResult := transformer.Replace(transformations, relationSet)
	wantReplaced := "quant(E1, isa(E1, ball), D1, isa(D1, every), none) ok(E1, D1)"

	if replacedResult.String() != wantReplaced {
		t.Errorf("RelationTransformer:\ngot\n%v,\nwant\n%v", replacedResult, wantReplaced)
	}
}
