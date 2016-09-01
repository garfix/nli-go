package tests

import (
	"fmt"
	"nli-go/lib/example3"
	"testing"
)

func TestSimpleRelationTransformer(test *testing.T) {

	internalGrammarParser := example3.NewSimpleInternalGrammarParser()

	// "name all customers"
	relations, _, _ := internalGrammarParser.CreateRelationSet(
		"[" +
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

	transformer := example3.NewSimpleRelationTransformer(transformations)
	transformedSet := transformer.Extract(relations)

	if len(transformedSet.Relations) != 1 {
		test.Error(fmt.Sprintf("Wrong number of relations: %d", len(transformedSet.Relations)))
	}

	relationString, sep := "", ""
	for _, relation := range transformedSet.Relations {
		relationString += sep + relation.ToString()
		sep = " "
	}

	if relationString != "task(S1, list_customers)" {
		test.Error("Error in relations: " + relationString)
	}
}
