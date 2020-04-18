package tests

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"testing"
)

func TestQuantSolver(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog(false)

	dbFacts := internalGrammarParser.CreateRelationSet(`[
		person(1, 'Jacqueline de Boer', 'F', '1964')
		person(2, 'Mark van Dongen', 'M', '1967')
		person(3, 'Suzanne van Dongen', 'F', '1967')
		person(4, 'John van Dongen', 'M', '1938')
		person(5, 'Dirk van Dongen', 'M', '1972')
		person(6, 'Durkje van Dongen', 'M', '1982')
		person(7, 'Huub de Boer', 'M', '1998')
		person(8, 'Babs de Boer', 'F', '1999')
		person(7, 'Johanneke de Boer', 'M', '1998')
		person(8, 'Baukje de Boer', 'F', '1999')
		have_child(4, 2)
		have_child(4, 3)
		have_child(1, 7)
		have_child(1, 8)
	]`)

	ds2db := internalGrammarParser.CreateRules(`[
		have_child(A, B) :- have_child(A, B);
		isa(A, parent) :- have_child(A, _);
		isa(A, child) :- have_child(_, A);
	]`)

	ds2dbWrite := internalGrammarParser.CreateRules(`[]`)

	tests := []struct {
		quant   string
		binding string
		result  string
	}{
		{
			// does every parent have 2 children?
			"quant(D1, [ isa(D1, all) ], S1, [ isa(S1, parent) ], [ have_child(S1, O1) number_of(O1, 2) ])",
			"{}",
			"{O1:2, S1:4}{O1:3, S1:4}{O1:7, S1:1}{O1:8, S1:1}",
		},
		{
			// does every parent have 3 children?
			"quant(D1, [ isa(D1, all) ], S1, [ isa(S1, parent) ], [quant(O1, [ isa(O1, child) ], D2, [ isa(D2, 3) ], [ have_child(S1, O1) ]) ])",
			"{}",
			"",
		},
		{
			// keep extra bindings?
			"quant(D1, [ isa(D1, all) ], S1, [ isa(S1, parent) ], [have_child(S1, O1) number_of(O1, 2) ])",
			"{X: 3}",
			"{O1:2, S1:4, X:3}{O1:3, S1:4, X:3}{O1:7, S1:1, X:3}{O1:8, S1:1, X:3}",
		},

// do 2 parents each have 2 children?
	}

	matcher := mentalese.NewRelationMatcher(log)

	entities := mentalese.Entities{}
	factBase1 := knowledge.NewInMemoryFactBase("memory", dbFacts, matcher, ds2db, ds2dbWrite, entities, log)
	dialogContext := central.NewDialogContext()
	predicates := mentalese.Predicates{}
	solver := central.NewProblemSolver(mentalese.NewRelationMatcher(log), predicates, dialogContext, log)
	solver.AddFactBase(factBase1)

	nestedStructureBase := knowledge.NewSystemNestedStructureBase(log)
	solver.AddNestedStructureBase(nestedStructureBase)

	aggregateBase := knowledge.NewSystemAggregateBase("system-aggregate", log)
	solver.AddMultipleBindingsBase(aggregateBase)

	for _, test := range tests {

		quant := internalGrammarParser.CreateRelation(test.quant)
		binding := internalGrammarParser.CreateBinding(test.binding)

		result := solver.SolveFind(quant, binding)
		result = mentalese.UniqueBindings(result)

		resultString := ""
		for _, result := range result {
			resultString += result.String()
		}

		if resultString != test.result {
			t.Errorf("got %s, want %s", resultString, test.result)
		}
	}
}
