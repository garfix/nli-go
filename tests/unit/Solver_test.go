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

	readMap := parser.CreateRules(`
		write(PersonName, BookName) :- book(BookId, BookName, _) author(PersonId, BookId) person(PersonId, PersonName);
		publish(PubName, BookName) :- book(BookId, BookName, PubId) publisher(PubId, PubName);
	`)

	writeMap := []mentalese.Rule{}

	matcher := central.NewRelationMatcher(log)

	factBase := knowledge.NewInMemoryFactBase("memory", facts, matcher, readMap, writeMap, log)

	variableGenerator := mentalese.NewVariableGenerator()
	solver := central.NewProblemSolver(matcher, variableGenerator, log)
	solver.AddFactBase(factBase)
	solver.Reindex()
	processList := central.NewProcessList()
	runner := central.NewProcessRunner(processList, solver, log)

	tests := []struct {
		input            string
		wantRelationSets string
	}{
		{"write('Sally Klein', B)", "[write('Sally Klein', 'The red book') write('Sally Klein', 'The green book')]"},
		{"write('Sally Klein', B) publish(P, B)", "[write('Sally Klein', 'The red book') publish('Orbital', 'The red book') write('Sally Klein', 'The green book') publish('Bookworm inc', 'The green book')]"},
		// stop processing when a predicate fails
		//{"missing_predicate() write('Sally Klein', B)", "[]"},
		// a failing predicate should remove existing bindings
		//{"write('Sally Klein', B) missing_predicate()", "[]"},
	}

	for _, test := range tests {

		input := parser.CreateRelationSet(test.input)
		resultBindings := runner.RunRelationSetWithBindings(central.NO_RESOURCE, input, mentalese.InitBindingSet(mentalese.NewBinding()))
		resultRelationSets := input.BindRelationSetMultipleBindings(resultBindings)

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

		resultBindings := runner.RunRelationSetWithBindings(central.NO_RESOURCE, []mentalese.Relation{input}, mentalese.InitBindingSet(binding))

		if fmt.Sprintf("%v", resultBindings) != test.wantResultBindings {
			t.Errorf("SolverTest: got %v, want %s", resultBindings, test.wantResultBindings)
		}
	}

	tests4 := []struct {
		input              string
		binding            string
		wantResultBindings string
	}{
		{"write(PersonName, BookName)", "{BookId:100, PersonName: 'John Graydon', BookName: 'The red book'}", "[{BookId:100, BookName:'The red book', PersonName:'John Graydon'}]"},
	}

	for _, test := range tests4 {

		input := parser.CreateRelation(test.input)
		binding := parser.CreateBinding(test.binding)
		resultBindings := runner.RunRelationSetWithBindings(central.NO_RESOURCE, mentalese.RelationSet{input}, mentalese.InitBindingSet(binding))

		if fmt.Sprintf("%v", resultBindings) != test.wantResultBindings {
			t.Errorf("SolverTest: got %v, want %s", resultBindings, test.wantResultBindings)
		}
	}

	rules2 := parser.CreateRules(`
		indirect_link(A, B) :- link(A, C) link(C, B);
	`)

	facts2 := parser.CreateRelationSet(`
		link('red', 'blue')
		link('blue', 'green')
		link('blue', 'yellow')
	`)

	readMap2 := parser.CreateRules(`
		link(A, B) :- link(A, B);
	`)

	factBase2 := knowledge.NewInMemoryFactBase("memory-1", facts2, matcher, readMap2, writeMap, log)
	ruleBase2 := knowledge.NewInMemoryRuleBase("memory-2", rules2, []string{}, log)

	solver2 := central.NewProblemSolver(matcher, variableGenerator, log)
	solver2.AddFactBase(factBase2)
	solver2.AddRuleBase(ruleBase2)
	solver2.Reindex()

	runner2 := central.NewProcessRunner(processList, solver2, log)

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
		resultBindings := runner2.RunRelationSetWithBindings(central.NO_RESOURCE, mentalese.RelationSet{input}, mentalese.InitBindingSet(binding))

		if fmt.Sprintf("%v", resultBindings) != test.wantResultBindings {
			t.Errorf("SolverTest: got %v, want %s", resultBindings, test.wantResultBindings)
		}
	}
}
