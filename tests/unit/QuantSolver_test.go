package tests

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/knowledge/nested"
	"nli-go/lib/mentalese"
	"testing"
)

func TestQuantSolver(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog(false)

	dbFacts := internalGrammarParser.CreateRelationSet(`
		person(1, 'Jacqueline de Boer', 'F', '1964')
		person(2, 'Mark van Dongen', 'M', '1967')
		person(3, 'Suzanne van Dongen', 'F', '1967')
		person(4, 'John van Dongen', 'M', '1938')
		person(5, 'Dirk van Dongen', 'M', '1972')
		person(6, 'Durkje van Dongen', 'M', '1982')
		person(7, 'Huub de Boer', 'M', '1998')
		person(8, 'Babs de Boer', 'F', '1999')
		person(9, 'Johanneke de Boer', 'M', '1998')
		person(10, 'Baukje de Boer', 'F', '1999')
		have_child(4, 2)
		have_child(4, 3)
		have_child(1, 7)
		have_child(1, 8)
		have_child(8, 9)
		have_child(8, 10)
	`)

	ds2db := internalGrammarParser.CreateRules(`[
		is_person(Id) :- person(Id, _, _, _);
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
			`
				quant_check(
					quant(quantifier(ResultCount, RangeCount, equals(ResultCount, RangeCount)), S1, isa(S1, parent)), 
					have_child(S1, O1) number_of(O1, 2))`,
			"{}",
			"{O1:2, S1:4}{O1:3, S1:4}{O1:7, S1:1}{O1:8, S1:1}{O1:9, S1:8}{O1:10, S1:8}",
		},
		{
			// does every parent have 3 children?
			`
				quant_check(
					quant(quantifier(ResultCount1, RangeCount1, equals(ResultCount1, RangeCount1)), S1, isa(S1, parent)),
					quant_check(
						quant(quantifier(ResultCount2, RangeCount2, equals(ResultCount1, 3)), O1, isa(O1, child))
				, 
				have_child(S1, O1)))`,
			"{}",
			"",
		},
		{
			// keep extra bindings?
			`
				quant_check(
					quant(quantifier(ResultCount, RangeCount, equals(ResultCount, RangeCount)), S1, isa(S1, parent)), 
					have_child(S1, O1) number_of(O1, 2)
				)`,
			"{X: 3}",
			"{O1:2, S1:4, X:3}{O1:3, S1:4, X:3}{O1:7, S1:1, X:3}{O1:8, S1:1, X:3}{O1:9, S1:8, X:3}{O1:10, S1:8, X:3}",
		},
		{
			// xor
			// the first quant in the xor has a range, but only the second quant has a range and scope bindings
			`
				quant_check(
					or(_,	
						and(_,
							quant(some, S1, is_person(S1) equals(S1, 8)),
							quant(some, S1, is_person(S1) have_child(S1, 9))
						),
						xor(_,
							quant(some, S1, is_person(S1) equals(S1, 4)),
							quant(some, S1, is_person(S1) equals(S1, 1))
						)
					),
					or(_, have_child(S1, 7), have_child(S1, 10))
				)`,
			"{}",
			"{S1:8}{S1:1}",
		},
	}

	matcher := mentalese.NewRelationMatcher(log)

	entities := mentalese.Entities{}
	factBase1 := knowledge.NewInMemoryFactBase("memory", dbFacts, matcher, ds2db, ds2dbWrite, entities, log)
	dialogContext := central.NewDialogContext()
	predicates := mentalese.Predicates{}
	solver := central.NewProblemSolver(mentalese.NewRelationMatcher(log), predicates, dialogContext, log)
	solver.AddFactBase(factBase1)

	systemFunctionBase := knowledge.NewSystemFunctionBase("system-function", log)
	solver.AddFunctionBase(systemFunctionBase)

	nestedStructureBase := nested.NewSystemNestedStructureBase(solver, dialogContext, predicates, log)
	solver.AddNestedStructureBase(nestedStructureBase)

	aggregateBase := knowledge.NewSystemAggregateBase("system-aggregate", log)
	solver.AddMultipleBindingsBase(aggregateBase)

	for _, test := range tests {

		log.Clear()

		quant := internalGrammarParser.CreateRelation(test.quant)
		binding := internalGrammarParser.CreateBinding(test.binding)

		result := solver.SolveRelationSet(mentalese.RelationSet{ quant }, mentalese.Bindings{ binding })
		result = result.UniqueBindings()

		resultString := ""
		for _, result := range result {
			resultString += result.String()
		}

		if !log.IsOk() {
			t.Error(log.String())
		}

		if resultString != test.result {
			t.Errorf("got %s, want %s", resultString, test.result)
			t.Error(log.String())
		}
	}
}
