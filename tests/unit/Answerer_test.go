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

	solutions := parser.CreateSolutions(`[
		{
			condition: write(Person_name, Book_name) publish(Pub_name, Book_name),
			no_results: {
				answer: none()
			},
			some_results: {
				answer: book(Book_name)
			}
		} {
			condition: write(Person, Book) number_of(N, Book),
			no_results: {
				answer: none()
			},
			some_results: {
				answer: focus(N)
			}
		} {
			condition: write(X, Y),
			no_results: {
				answer: none()
			},
			some_results: {
				answer: write(X, Y)
			}
		} {
			condition: publish(A, B),
			no_results: {
				answer: none()
			},
			some_results: {
				preparation: write(C, B),
				answer: publish_author(A, C)
			}
		}
	]`)

	matcher := mentalese.NewRelationMatcher(log)

	stats := mentalese.DbStats{
		"book": {Size: 100, DistinctValues: []int{100, 100}},
		"author": {Size: 100, DistinctValues: []int{100, 200}},
		"person": {Size: 100, DistinctValues: []int{100, 100}},
	}
	factBase := knowledge.NewInMemoryFactBase(facts, matcher, ds2db, stats, log)
	systemAggregateBase := knowledge.NewSystemAggregateBase(log)

	answerer := central.NewAnswerer(matcher, log)
	answerer.AddMultipleBindingsBase(systemAggregateBase)
	answerer.AddFactBase(factBase)
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
		{"[write('Sally Klein', Book) number_of(N, Book)]", "[focus(2)]"},
	}

	for _, test := range tests {

		input := parser.CreateRelationSet(test.input)

		resultRelationSet := answerer.Answer(input)

		if fmt.Sprintf("%v", resultRelationSet) != test.wantRelationSet {
			t.Errorf("Answerer(%v): got %v, want %s", test.input, resultRelationSet, test.wantRelationSet)
		}
	}
}

func TestUnscope(t *testing.T) {

	parser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog(false)

	matcher := mentalese.NewRelationMatcher(log)
	answerer := central.NewAnswerer(matcher, log)

	tests := []struct {
		input           string
		wantRelationSet string
	}{
		// use all arguments
		{"[abc(A, 1) quant(A, [isa(A, 1)], B, [isa(B, 2)], [make(A, B)])]", "[abc(A, 1) isa(A, 1) isa(B, 2) make(A, B) quant(A, [], B, [], [])]"},
		// recurse
		{"[quant(A, [ quant(A, [ isa(A, 1) ], B, [], []) ], B, [], [])]", "[isa(A, 1) quant(A, [], B, [], []) quant(A, [], B, [], [])]"},
	}

	for _, test := range tests {

		input := parser.CreateRelationSet(test.input)

		resultRelationSet := answerer.Unscope(input)

		if resultRelationSet.String() != test.wantRelationSet {
			t.Errorf("Answerer: got %v, want %s", resultRelationSet, test.wantRelationSet)
		}
	}
}
