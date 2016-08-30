package tests

import (
	"fmt"
	"nli-go/lib/example3"
	"testing"
)

func TestSimpleRelationTransformer(test *testing.T) {

	// "name all customers"
	relations := []example3.SimpleRelation{
		{Predicate: "predicate", Arguments: []example3.SimpleTerm{{example3.Term_variable, "S1"}, {example3.Term_predicateAtom, "name"}}},
		{Predicate: "object", Arguments: []example3.SimpleTerm{{example3.Term_variable, "S1"}, {example3.Term_variable, "E1"}}},
		{Predicate: "instance_of", Arguments: []example3.SimpleTerm{{example3.Term_variable, "E1"}, {example3.Term_predicateAtom, "customer"}}},
		{Predicate: "determiner", Arguments: []example3.SimpleTerm{{example3.Term_variable, "E1"}, {example3.Term_variable, "D1"}}},
		{Predicate: "instance_of", Arguments: []example3.SimpleTerm{{example3.Term_variable, "D1"}, {example3.Term_predicateAtom, "all"}}},
	}

	transformations := []example3.SimpleRelationTransformation {
		// list-customers(P1) :- predicate(P1, name), object(P1, E1), instance_of(E1, customer)
		{
			Pattern: []example3.SimpleRelation{
				{Predicate: "predicate", Arguments: []example3.SimpleTerm{{example3.Term_variable, "P1"}, {example3.Term_predicateAtom, "name"}}},
				{Predicate: "object", Arguments: []example3.SimpleTerm{{example3.Term_variable, "P1"}, {example3.Term_variable, "O1"}}},
				{Predicate: "instance_of", Arguments: []example3.SimpleTerm{{example3.Term_variable, "O1"}, {example3.Term_predicateAtom, "customer"}}},
			},
			Replacement: []example3.SimpleRelation{
				{Predicate: "task", Arguments: []example3.SimpleTerm{{example3.Term_variable, "P1"}, {example3.Term_predicateAtom, "list_customers"}}},
			},
		},
	}

	transformer := example3.NewSimpleRelationTransformer(transformations)
	transformedRelations := transformer.Extract(relations)

	if len(transformedRelations) != 1 {
		test.Error(fmt.Sprintf("Wrong number of relations: %d", len(transformedRelations)))
	}

	relationString, sep := "", ""
	for _, relation := range transformedRelations {
		relationString += sep + relation.ToString()
		sep = " "
	}

	if relationString != "task(S1, list_customers)" {
		test.Error("Error in relations: " + relationString)
	}
}
