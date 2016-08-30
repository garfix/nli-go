package tests

import (
	"fmt"
	"nli-go/lib/example3"
	"testing"
)

func TestSimpleRelationTransformer(test *testing.T) {

	// "name all customers"
	relations := []example3.SimpleRelation{
		{Predicate: "predicate", Arguments: []string{"S1", "name"}},
		{Predicate: "object", Arguments: []string{"S1", "E1"}},
		{Predicate: "instance_of", Arguments: []string{"E1", "customer"}},
		{Predicate: "determiner", Arguments: []string{"E1", "D1"}},
		{Predicate: "instance_of", Arguments: []string{"D1", "all"}},
	}

	transformations := []example3.SimpleRelationTransformation {
		// list-customers(P1) :- predicate(P1, name), object(P1, E1), instance_of(E1, customer)
		{
			Pattern: []example3.SimpleRelation{
				{Predicate: "predicate", Arguments: []string{"P1", "name"}},
				{Predicate: "object", Arguments: []string{"P1", "O1"}},
				{Predicate: "instance_of", Arguments: []string{"O1", "customer"}},
			},
			Replacement: []example3.SimpleRelation{
				{Predicate: "task", Arguments: []string{"P1", "list_customers"}},
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
