package tests

import (
	"testing"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/common"
)

func TestFactBase(t *testing.T) {

	parser := importer.NewInternalGrammarParser()

	facts, _, _ := parser.CreateRelationSet(`
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

	rules, _, _ := parser.CreateRules(`
		write(PersonName, BookName) :- book(BookId, BookName, _), author(PersonId, BookId), person(PersonId, PersonName)
	`)

	factBase := knowledge.NewFactBase(facts, rules)

	tests := []struct {
		input string
		wantRelations string
		wantBindings string
	} {
		{"write('Sally Klein', B)", "", "[{B:'The red book'}{B:'The green book'}]"},
	}

	for _, test := range tests {

		input, _ := parser.CreateRelation(test.input)
common.LoggerActive=false
		_, resultBindings := factBase.Bind(input)
common.LoggerActive=false
		bind := "["
		for _, resultBinding := range resultBindings  {
			bind += resultBinding.String()
		}
		bind += "]"

		if bind != test.wantBindings {
			t.Errorf("FactBase,Bind(%v): got %s, want %s", test.input, bind, test.wantBindings)
		}
	}
}