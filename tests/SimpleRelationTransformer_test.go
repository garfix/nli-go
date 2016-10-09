package tests

import (
	"testing"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
)

func TestSimpleRelationTransformer(test *testing.T) {

	internalGrammarParser := importer.NewSimpleInternalGrammarParser()

	// "name all customers"
	relationSet, _, _ := internalGrammarParser.CreateRelationSet(
		"[" +
			"instance_of(E2, name)" +
			"predicate(S1, name)" +
			"object(S1, E1)" +
			"instance_of(E1, customer)" +
			"determiner(E1, D1)" +
			"instance_of(D1, all)" +
		"]")

	transformations, _, _ := internalGrammarParser.CreateTransformations(
		"[" +
			"task(P1, list_customers) :- predicate(P1, name), object(P1, O1), instance_of(O1, customer)" +
		"]")

	transformations2, _, _ := internalGrammarParser.CreateTransformations(
		"[" +
			"task(A, B), subject(Y) :- predicate(A, X), object(A, Y), determiner(Y, Z), instance_of(Z, B)" +
			"done() :- predicate(A, X), object(A, Y), determiner(Y, Z), instance_of(Z, B)" +
			"magic(A, X) :- predicate(A, X), predicate(X, A)" +
		"]")

	transformer := mentalese.NewSimpleRelationTransformer(transformations)
	transformer2 := mentalese.NewSimpleRelationTransformer(transformations2)

	// extract

	transformedSet, _ := transformer.Extract(relationSet)

	if transformedSet[0].String() != "[task(S1, list_customers)]" {
		test.Errorf("Error in result: %s", transformedSet[0].String())
	}

	transformedSet, _ = transformer2.Extract(relationSet)

	if transformedSet[0].String() != "[task(S1, name) subject(E1)]" {
		test.Errorf("Error in result: %s", transformedSet[0].String())
	}
	if transformedSet[1].String() != "[done()]" {
		test.Errorf("Error in result: %s", transformedSet[1].String())
	}
	if len(transformedSet) != 2 {
//		test.Errorf("Error in length: %d", len(transformedSet))
	}

	// replace

	//transformedSet2 := transformer2.Replace(relationSet)
	//
	//if transformedSet2.String() != "[instance_of(E2, name) instance_of(E1, customer) task(S1, all) subject(E1) done()]" {
	//	test.Errorf("Error in result: %s", transformedSet2.String())
	//}
	//
	//// append
	//
	//transformedSet2 = transformer2.Append(relationSet)
	//
	//if transformedSet2.String() != "[instance_of(E2, name) predicate(S1, name) object(S1, E1) instance_of(E1, customer) determiner(E1, D1) " +
	//	"instance_of(D1, all) task(S1, name) subject(E1) done()]" {
	//	test.Errorf("Error in result: %s", transformedSet2.String())
	//}
}
