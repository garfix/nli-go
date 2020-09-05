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

	matcher := mentalese.NewRelationMatcher(log)

	entities := mentalese.Entities{}
	factBase := knowledge.NewInMemoryFactBase("memory", facts, matcher, ds2db, ds2dbWrite, entities, log)

	dialogContext := central.NewDialogContext()
	solver := central.NewProblemSolver(matcher, dialogContext, log)
	solver.AddFactBase(factBase)

	tests := []struct {
		input            string
		wantRelationSets string
	}{
		{"write('Sally Klein', B)", "[write('Sally Klein', 'The red book') write('Sally Klein', 'The green book')]"},
		{"write('Sally Klein', B) publish(P, B)", "[write('Sally Klein', 'The red book') publish('Orbital', 'The red book') write('Sally Klein', 'The green book') publish('Bookworm inc', 'The green book')]"},
		// stop processing when a predicate fails
		{"missing_predicate() write('Sally Klein', B)", "[]"},
		//// a failing predicate should remove existing bindings
		{"write('Sally Klein', B) missing_predicate()", "[]"},
	}

	for _, test := range tests {

		input := parser.CreateRelationSet(test.input)
		resultBindings := solver.SolveRelationSet(input, mentalese.Bindings{{}})
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

		resultBindings := solver.SolveRelationSet([]mentalese.Relation{input}, mentalese.Bindings{binding})

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
		resultBindings := solver.SolveRelationSet(mentalese.RelationSet{ input }, mentalese.Bindings{ binding })

		if fmt.Sprintf("%v", resultBindings) != test.wantResultBindings {
			t.Errorf("SolverTest: got %v, want %s", resultBindings, test.wantResultBindings)
		}
	}

	rules2 := parser.CreateRules(`[
		indirect_link(A, B) :- link(A, C) link(C, B);
	]`)

	facts2 := parser.CreateRelationSet(`
		link('red', 'blue')
		link('blue', 'green')
		link('blue', 'yellow')
	`)

	ds2db2 := parser.CreateRules(`[
		link(A, B) :- link(A, B);
	]`)

	factBase2 := knowledge.NewInMemoryFactBase("memory-1", facts2, matcher, ds2db2, ds2dbWrite, entities, log)
	ruleBase2 := knowledge.NewInMemoryRuleBase("memory-2", rules2, log)

	solver2 := central.NewProblemSolver(matcher, dialogContext, log)
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
		resultBindings := solver2.SolveRelationSet(mentalese.RelationSet{ input }, mentalese.Bindings{ binding})

		if fmt.Sprintf("%v", resultBindings) != test.wantResultBindings {
			t.Errorf("SolverTest: got %v, want %s", resultBindings, test.wantResultBindings)
		}
	}
}

func TestMissingHandlerError(t *testing.T) {

	parser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog(false)

	facts := mentalese.RelationSet{}
	ds2db := parser.CreateRules(`[]`)
	ds2dbWrite := parser.CreateRules(`[]`)
	matcher := mentalese.NewRelationMatcher(log)
	entities := mentalese.Entities{}
	factBase := knowledge.NewInMemoryFactBase("memory", facts, matcher, ds2db, ds2dbWrite, entities, log)

	dialogContext := central.NewDialogContext()

	solver := central.NewProblemSolver(matcher, dialogContext, log)
	solver.AddFactBase(factBase)

	input := parser.CreateRelationSet("not_a_relation()")
	bindings := parser.CreateBindings("[{}]")
	resultBindings := solver.SolveRelationSet(input, bindings)

	if fmt.Sprintf("%v", resultBindings) != "[]" {
		t.Errorf("SolverTest: got %v, want []", resultBindings)
	}

	errors := log.GetErrors()
	if len(errors) == 0 {
		t.Errorf("Expected error message")
	} else if errors[0] != "Predicate not supported by any knowledge base: not_a_relation" {
		t.Errorf("Unexpected message: " + errors[0])
	}
}