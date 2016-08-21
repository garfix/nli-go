package tests

import (
	"fmt"
	"nli-go/lib/example2"
	"nli-go/lib/example3"
	"testing"
)

func TestSimpleRelationTransformer(test *testing.T) {

	relations := []example2.SimpleRelation{
		{Predicate: "predicate", Arguments: []string{"S1", "name"}},
		{Predicate: "object", Arguments: []string{"S1", "E1"}},
		{Predicate: "instance-of", Arguments: []string{"E1", "customer"}},
		{Predicate: "determiner", Arguments: []string{"E1", "D1"}},
		{Predicate: "instance-of", Arguments: []string{"D1", "all"}},
	}

	transformations := []example3.SimpleRelationTransformation {
		{
			Pattern: []example2.SimpleRelation{},
			Replacement: []example2.SimpleRelation{},
		},
	}

	transformer := example3.NewSimpleRelationTransformer(transformations)

	transformedRelations := transformer.Process(relations)

	if len(transformedRelations) != 4 {
		test.Error(fmt.Sprintf("Wrong number of relations: %d", len(transformedRelations)))
	}

	relationString, sep := "", " "
	for i := 0; i < len(transformedRelations); i++ {
		relationString += sep + RelationToString(transformedRelations[i])
		sep = " "
	}
	if relationString != "" {
		test.Error("Error in relations: " + relationString)
	}
}
