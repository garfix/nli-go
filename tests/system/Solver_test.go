package tests

import (
	"testing"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"fmt"
	"nli-go/lib/central"
)

func TestSolver(t *testing.T) {

	parser := importer.NewInternalGrammarParser()

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

	rules := parser.CreateRules(`
		write(PersonName, BookName) :- book(BookId, BookName, _) author(PersonId, BookId) person(PersonId, PersonName);
		publish(PubName, BookName) :- book(BookId, BookName, PubId) publisher(PubId, PubName);
	`)

	factBase := knowledge.NewFactBase(facts, rules)

	solver := central.NewProblemSolver()
	solver.AddKnowledgeBase(factBase)

	tests := []struct {
		input string
		wantRelationSets string
	} {
		{"[write('Sally Klein', B)]", "[[write('Sally Klein', 'The red book')] [write('Sally Klein', 'The green book')]]"},
		{"[write('Sally Klein', B) publish(P, B)]", "[[write('Sally Klein', 'The red book') publish('Orbital', 'The red book')] [write('Sally Klein', 'The green book') publish('Bookworm inc', 'The green book')]]"},
	}

	for _, test := range tests {

		input := parser.CreateRelationSet(test.input)

		resultRelationSets := solver.Solve(input)

		if fmt.Sprintf("%v", resultRelationSets) != test.wantRelationSets {
			t.Errorf("FactBase,Bind(%v): got %v, want %s", test.input, resultRelationSets, test.wantRelationSets)
		}
	}
}