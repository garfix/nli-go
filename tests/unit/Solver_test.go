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

func TestSolver(t *testing.T) {

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
		write(PersonName, BookName) => book(BookId, BookName, _) author(PersonId, BookId) person(PersonId, PersonName);
		publish(PubName, BookName) => book(BookId, BookName, PubId) publisher(PubId, PubName);
	]`)

	matcher := mentalese.NewRelationMatcher(log)

	stats := mentalese.DbStats{}
	factBase := knowledge.NewInMemoryFactBase(facts, matcher, ds2db, stats, log)

	solver := central.NewProblemSolver(matcher, log)
	solver.AddFactBase(factBase)

	tests := []struct {
		input            string
		wantRelationSets string
	}{
		{"[write('Sally Klein', B)]", "[[write('Sally Klein', 'The red book')] [write('Sally Klein', 'The green book')]]"},
		{"[write('Sally Klein', B) publish(P, B)]", "[[write('Sally Klein', 'The red book') publish('Orbital', 'The red book')] [write('Sally Klein', 'The green book') publish('Bookworm inc', 'The green book')]]"},
		// stop processing when a predicate fails
		{"[missingPredicate() write('Sally Klein', B)]", "[]"},
		//// a failing predicate should remove existing bindings
		{"[write('Sally Klein', B) missingPredicate()]", "[]"},
	}

	for _, test := range tests {

		input := parser.CreateRelationSet(test.input)
		resultBindings := solver.SolveRelationSet(input, []mentalese.Binding{})
		resultRelationSets := matcher.BindRelationSetMultipleBindings(input, resultBindings)

		if fmt.Sprintf("%v", resultRelationSets) != test.wantRelationSets {
			t.Errorf("SolverTest: got %v, want %s", resultRelationSets, test.wantRelationSets)
		}
	}

	tests2 := []struct {
		input              string
		binding            string
		wantResultBindings string
	}{
		{"publish('Bookworm inc', B)", "{}", "[{B:'The green book'} {B:'The blue book'}]"},
		//{"publish('Bookworm inc', B)", "{X:B}", "[{B:'The green book', X:B} {B:'The blue book', X:B}]"},
		//{"publish('Bookworm inc', B)", "{B:X}", "[{B:'The green book'} {B:'The blue book'}]"},
		{"publish('Bookworm inc', B)", "{A:1}", "[{A:1, B:'The green book'} {A:1, B:'The blue book'}]"},
		{"publish('Bookworm inc', B)", "{B:'The green book'}", "[{B:'The green book'}]"},
	}

	for _, test := range tests2 {

		input := parser.CreateRelation(test.input)
		binding := parser.CreateBinding(test.binding)

		resultBindings := solver.SolveRelationSet([]mentalese.Relation{input}, []mentalese.Binding{binding})

		if fmt.Sprintf("%v", resultBindings) != test.wantResultBindings {
			t.Errorf("SolverTest: got %v, want %s", resultBindings, test.wantResultBindings)
		}
	}

	rules2 := parser.CreateRules(`[
		indirect_link(A, B) :- link(A, C) link(C, B);
	]`)

	facts2 := parser.CreateRelationSet(`[
		link('red', 'blue')
		link('blue', 'green')
		link('blue', 'yellow')
	]`)

	ds2db2 := parser.CreateTransformations(`[
		link(A, B) => link(A, B);
	]`)

	factBase2 := knowledge.NewInMemoryFactBase(facts2, matcher, ds2db2, stats, log)
	ruleBase2 := knowledge.NewRuleBase(rules2, log)

	solver2 := central.NewProblemSolver(matcher, log)
	solver2.AddFactBase(factBase2)
	solver2.AddRuleBase(ruleBase2)

	tests3 := []struct {
		input              string
		binding            string
		wantResultBindings string
	}{
		{"indirect_link(X, Y)", "{}", "[{X:'red', Y:'green'} {X:'red', Y:'yellow'}]"},
		{"indirect_link(X, Y)", "{ Y:'yellow' }", "[{X:'red', Y:'yellow'}]"},
	}

	for _, test := range tests3 {

		input := parser.CreateRelation(test.input)
		binding := parser.CreateBinding(test.binding)
		resultBindings := solver2.SolveSingleRelationSingleBindingSingleRuleBase(input, binding, ruleBase2)

		if fmt.Sprintf("%v", resultBindings) != test.wantResultBindings {
			t.Errorf("SolverTest: got %v, want %s", resultBindings, test.wantResultBindings)
		}
	}

}
