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

	facts := parser.CreateRelationSet(`
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
	`)

	ds2db := parser.CreateRules(`[
		write(PersonName, BookName) :- book(BookId, BookName, _) author(PersonId, BookId) person(PersonId, PersonName);
		publish(PubName, BookName) :- book(BookId, BookName, PubId) publisher(PubId, PubName);
	]`)

	ds2dbWrite := parser.CreateRules(`[]`)

	entities := mentalese.Entities{}
	matcher := mentalese.NewRelationMatcher(log)
	factBase := knowledge.NewInMemoryFactBase("memory", facts, matcher, ds2db, ds2dbWrite, entities, log)
	dialogContext := central.NewDialogContext()
	predicates := &mentalese.Predicates{}
	solver := central.NewProblemSolver(matcher, predicates, dialogContext, log)

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

		resultBindings := solver.FindFacts(factBase, input, mentalese.Binding{})

		if fmt.Sprintf("%v", resultBindings) != test.wantBindings {
			t.Errorf("FactBase,BindSingle(%v): got %v, want %s", test.input, resultBindings, test.wantBindings)
		}
	}
}
