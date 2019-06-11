package tests

import (
	"fmt"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"testing"
	"nli-go/lib/mentalese"
	"nli-go/lib/central"
)

func TestFactBase(t *testing.T) {

	parser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog(false)

	facts := parser.CreateRelationSet(`[
		book(1, 'The red book', 5)
		book(2, 'The green book', 6)
		book(3, 'The blue book', 6)
		publisher(5, 'Orbital')
		publisher(6, 'Bookworm inc')
		author(8, 1)
		author(9, 1)
		author(9, 2)
		author(10, 1)
		author(11, 3)
		person(8, 'John Graydon')
		person(9, 'Sally Klein')
		person(10, 'Keith Partridge')
		person(11, 'Onslow Bigbrain')
	]`)

	ds2db := parser.CreateTransformations(`[
		write(Person_name, Book_name) => book(Book_id, Book_name, _) author(Person_id, Book_id) person(Person_id, Person_name);
		publish(Pub_name, Book_name) => book(Book_id, Book_name, Pub_id) publisher(Pub_id, Pub_name);
	]`)

	ds2dbWrite := parser.CreateTransformations(`[]`)

	stats := mentalese.DbStats{}
	entities := mentalese.Entities{}
	matcher := mentalese.NewRelationMatcher(log)
	factBase := knowledge.NewInMemoryFactBase("memory", facts, matcher, ds2db, ds2dbWrite, stats, entities, log)
	solver := central.NewProblemSolver(matcher, log)

	tests := []struct {
		input         string
		wantRelations string
		wantBindings  string
	}{
		{"write('Sally Klein', B)", "", "[{B:'The red book'} {B:'The green book'}]"},
		{"publish(X, Y)", "", "[{X:'Orbital', Y:'The red book'} {X:'Bookworm inc', Y:'The green book'} {X:'Bookworm inc', Y:'The blue book'}]"},
		{"write('Keith Partridge', 'The red book')", "", "[{}]"},
	}

	for _, test := range tests {

		input := parser.CreateRelation(test.input)

		resultBindings := solver.FindFacts(factBase, mentalese.RelationSet{input})

		if fmt.Sprintf("%v", resultBindings) != test.wantBindings {
			t.Errorf("FactBase,Bind(%v): got %v, want %s", test.input, resultBindings, test.wantBindings)
		}
	}
}
