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

	solutions := parser.CreateSolutions(`[
		{
			condition: write(Person_name, Book_name) publish(Pub_name, Book_name),
			responses: [
				{
					condition: exists(),
					answer: book(Book_name)
				}
				{
					answer: none()
				}
			]
		} {
			condition: write(Person, Book) number_of(Book, N),
			responses: [
				{
					condition: exists(),
					answer: focus(N)
				}
				{
					answer: none()
				}
			]
		} {
			condition: write(X, Y),
			responses: [
				{
					condition: exists(),
					answer: write(X, Y)
				}
				{
					answer: none()
				}
			]
		} {
			condition: publish(A, B),
			responses: [
				{
					condition: exists(),
					preparation: write(C, B),
					answer: publish_author(A, C)
				}
				{
					answer: none()
				}
			]
		}
	]`)

	matcher := mentalese.NewRelationMatcher(log)

	stats := mentalese.DbStats{
		"book": {Size: 100, DistinctValues: []int{100, 100}},
		"author": {Size: 100, DistinctValues: []int{100, 200}},
		"person": {Size: 100, DistinctValues: []int{100, 100}},
	}

	entities := mentalese.Entities{}

	factBase := knowledge.NewInMemoryFactBase("memory", facts, matcher, ds2db, ds2dbWrite, stats, entities, log)
	systemAggregateBase := knowledge.NewSystemAggregateBase("system-aggregate", log)
	predicates := mentalese.Predicates{}

	dialogContext := central.NewDialogContext()
	solver := central.NewProblemSolver(matcher, predicates, dialogContext, log)
	solver.AddMultipleBindingsBase(systemAggregateBase)
	solver.AddFactBase(factBase)

	answerer := central.NewAnswerer(matcher, solver, log)
	answerer.AddSolutions(solutions)

	tests := []struct {
		input           string
		wantRelationSet string
	}{
		// simple
		{"[write('Sally Klein', B)]", "[write('Sally Klein', 'The red book') write('Sally Klein', 'The green book')]"},
		// preparation
		{"[publish('Bookworm inc', B)]", "[publish_author('Bookworm inc', 'Sally Klein') publish_author('Bookworm inc', 'Onslow Bigbrain')]"},
		//// return each relation only once
		{"[write(Person_name, B) publish('Orbital', B)]", "[book('The red book')]"},
		// number_of
		{"[write('Sally Klein', Book) number_of(Book, N)]", "[focus(2)]"},
	}

	for _, test := range tests {

		input := parser.CreateRelationSet(test.input)

		resultRelationSet := answerer.Answer(input, mentalese.Bindings{{}})

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
		{"[abc(A, 1) quant(B, [isa(B, 2)], A, [isa(A, 1)], [make(A, B)])]", "[abc(A, 1) isa(B, 2) isa(A, 1) make(A, B) quant(B, [], A, [], [])]"},
		// recurse
		{"[quant(A, [ quant(B, [], A, [ isa(A, 1) ], []) ], B, [], [])]", "[isa(A, 1) quant(B, [], A, [], []) quant(A, [], B, [], [])]"},
	}

	for _, test := range tests {

		input := parser.CreateRelationSet(test.input)

		resultRelationSet := input.UnScope()

		if resultRelationSet.String() != test.wantRelationSet {
			t.Errorf("Answerer: got %v, want %s", resultRelationSet, test.wantRelationSet)
		}
	}
}
