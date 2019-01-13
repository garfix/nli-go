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
	relationSet := parser.CreateRelationSet(`[
		instance_of(E2, name)
		predicate(S1, name)
		object(S1, E1)
		instance_of(E1, customer)
		determiner(E1, D1)
		instance_of(D1, all)
	]`)

	tests := []struct {
		transformations string
		wantReplaced    string
	}{
		{
			`[
				predicate(A, X) object(A, Y) determiner(Y, Z) instance_of(Z, B) => task(A, B) subject(Y);
				predicate(A, X) object(A, Y) determiner(Y, Z) instance_of(Z, B) => done();
				predicate(A, X) predicate(X, A) => magic(A, X);
				IF object(A, O) THEN predicate(A, X) => label(A, O);
				IF object(A, A) THEN predicate(A, X) => signal(A, X);
			]`,
			"[instance_of(E2, name) instance_of(E1, customer) task(S1, all) subject(E1) done() label(S1, E1)]",
		},
		{
			`[
				instance_of(Z, B) => isa(Z, B);
			]`,
			"[predicate(S1, name) object(S1, E1) determiner(E1, D1) isa(E2, name) isa(E1, customer) isa(D1, all)]",
		},
	}
	for _, test := range tests {

		transformations := parser.CreateTransformations(test.transformations)

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

	relationSet := parser.CreateRelationSet(`[
		quant(E1, [ isa(E1, ball) ], D1, [ isa(D1, every) ], [])
	]`)

	transformations := parser.CreateTransformations(`[
		quant(E2, [ isa(E2, ball) ], D2, [ isa(D2, every) ], []) => quant(E2, [ isa(E2, ball) ], D2, [ isa(D2, every) ], []) ok(E2, D2);
	]`)

	replacedResult := transformer.Replace(transformations, relationSet)
	wantReplaced := "[quant(E1, [isa(E1, ball)], D1, [isa(D1, every)], []) ok(E1, D1)]"

	if replacedResult.String() != wantReplaced {
		t.Errorf("RelationTransformer:\ngot\n%v,\nwant\n%v", replacedResult, wantReplaced)
	}
}

func TestRelationTransformerWithQuant(t *testing.T) {

	log := common.NewSystemLog(false)
	parser := importer.NewInternalGrammarParser()
	matcher := mentalese.NewRelationMatcher(log)
	transformer := mentalese.NewRelationTransformer(matcher, log)

	input := parser.CreateRelationSet(`[
		quant(E1, [ isa(E1, how) isa(E1, many)], D1, [ isa(D1, how) isa(D1, many) ], [])
	]`)

	transformations := parser.CreateTransformations(`[
		isa(A, how) isa(A, many) => how_many(A);
	]`)

	result := transformer.Replace(transformations, input)
	want := "[quant(E1, [how_many(E1)], D1, [how_many(D1)], [])]"

	if result.String() != want {
		t.Errorf("RelationTransformer:\ngot\n%v,\nwant\n%v", result, want)
	}
}
