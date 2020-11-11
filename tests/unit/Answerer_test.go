package tests

import (
	"fmt"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"testing"
)

func TestAnswerer(t *testing.T) {

	parser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog()

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

	ds2db := parser.CreateRules(`
		write(PersonName, BookName) :- book(BookId, BookName, _) author(PersonId, BookId) person(PersonId, PersonName);
		publish(PubName, BookName) :- book(BookId, BookName, PubId) publisher(PubId, PubName);
	`)

	ds2dbWrite := []mentalese.Rule{}

	solutions := parser.CreateSolutions(`
		{
			condition: write(PersonName, BookName) publish(PubName, BookName),
			result: _,
			responses: 
				{
					condition: go:exists(),
					answer: book(BookName)
				}
				{
					answer: none()
				}
		} 
		{
			condition: write(Person, Book) go:number_of(Book, N),
			result: _,
			responses: 
				{
					condition: go:exists(),
					answer: focus(N)
				}
				{
					answer: none()
				}
		} 
		{
			condition: write(X, Y),
			result: _,
			responses: 
				{
					condition: go:exists(),
					answer: write(X, Y)
				}
				{
					answer: none()
				}
		} 
		{
			condition: publish(A, B),
			result: _,
			responses: 
				{
					condition: go:exists(),
					preparation: write(C, B),
					answer: publish_author(A, C)
				}
				{
					answer: none()
				}
		}
	`)

	matcher := central.NewRelationMatcher(log)

	factBase := knowledge.NewInMemoryFactBase("memory", facts, matcher, ds2db, ds2dbWrite, log)
	systemAggregateBase := knowledge.NewSystemMultiBindingBase("system-aggregate", log)

	dialogContext := central.NewDialogContext()
	solver := central.NewProblemSolver(matcher, dialogContext, log)
	solver.AddMultipleBindingBase(systemAggregateBase)
	solver.AddFactBase(factBase)

	answerer := central.NewAnswerer(matcher, solver, log)
	answerer.AddSolutions(solutions)

	tests := []struct {
		input           string
		wantRelationSet string
	}{
		// simple
		{"write('Sally Klein', B)", "write('Sally Klein', 'The red book') write('Sally Klein', 'The green book')"},
		// preparation
		{"publish('Bookworm inc', B)", "publish_author('Bookworm inc', 'Sally Klein') publish_author('Bookworm inc', 'Onslow Bigbrain')"},
		//// return each relation only once
		{"write(PersonName, B) publish('Orbital', B)", "book('The red book')"},
		// number_of
		{"write('Sally Klein', Book) go:number_of(Book, N)", "focus(2)"},
	}

	for _, test := range tests {

		input := parser.CreateRelationSet(test.input)

		resultRelationSet := answerer.Answer(input, mentalese.InitBindingSet(mentalese.NewBinding()))

		if fmt.Sprintf("%v", resultRelationSet) != test.wantRelationSet {
			t.Errorf("Answerer(%v): got %v, want %s", test.input, resultRelationSet, test.wantRelationSet)
		}
	}
}

func TestUnScope(t *testing.T) {

	parser := importer.NewInternalGrammarParser()
	tests := []struct {
		input           string
		wantRelationSet string
	}{
		// use all arguments
		{"abc(A, 1) quant(B, isa(B, 2), A, isa(A, 1), make(A, B))", "abc(A, 1) isa(B, 2) isa(A, 1) make(A, B) quant(B, none, A, none, none)"},
		// recurse
		{"quant(A, quant(B, none, A, isa(A, 1), none), B, none, none)", "isa(A, 1) quant(B, none, A, none, none) quant(A, none, B, none, none)"},
	}

	for _, test := range tests {

		input := parser.CreateRelationSet(test.input)

		resultRelationSet := input.UnScope()

		if resultRelationSet.String() != test.wantRelationSet {
			t.Errorf("Answerer: got %v, want %s", resultRelationSet, test.wantRelationSet)
		}
	}
}
