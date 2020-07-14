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

func TestInMemoryRuleBase(t *testing.T) {

	parser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog(false)
	matcher := mentalese.NewRelationMatcher(log)
	dialogContext := central.NewDialogContext()
	predicates := mentalese.Predicates{}
	solver := central.NewProblemSolver(matcher, predicates, dialogContext, log)
	facts := parser.CreateRelationSet(`
		parent(john, jack)
		parent(james, jack)
		parent(vince, bob)
		parent(pat, bob)
		parent(sue, cyrill)
		-sibling(alice, bob)
	`)
	entities := mentalese.Entities{}
	ds2db := parser.CreateRules(`[
		parent(A, B) :- parent(A, B);
		-sibling(A, B) :- -sibling(A, B);
	]`)
	ds2dbWrite := parser.CreateRules(`[]`)
	factBase := knowledge.NewInMemoryFactBase("memory", facts, matcher, ds2db, ds2dbWrite, entities, log)
	solver.AddFactBase(factBase)
	functionBase := knowledge.NewSystemFunctionBase("function", log)
	solver.AddFunctionBase(functionBase)
	nestedBase := nested.NewSystemNestedStructureBase(solver, dialogContext, predicates, log)
	solver.AddNestedStructureBase(nestedBase)
	rules := parser.CreateRules(`[
		sibling(A, B) :- parent(A, C) parent(B, C) not( -sibling(A, B) );
		-sibling(A, B) :- equals(A, B);
	]`)
	ruleBase := knowledge.NewInMemoryRuleBase("mem", rules, log)
	solver.AddRuleBase(ruleBase)

	tests := []struct {
		goal           string
		binding        string
		resultBindings string
	}{
		// do not promote temporary variable C
		// exception
		{"sibling(X, Y)", "{X:john, Z:sue}", "[{X:john, Y:james, Z:sue}]"},
		// do not enter unnecessary variables, because they may conflict with the temporary variables
		{"sibling(X, Y)", "{X:john, C:sue}", "[{C:sue, X:john, Y:james}]"},
		{"sibling(X, Y)", "{X:bob}", "[]"},
		// negative succeed
		{"-sibling(X, Y)", "{X:john, Y:john}", "[{X:john, Y:john}]"},
		{"-sibling(X, Y)", "{X:alice, Y:bob}", "[{X:alice, Y:bob}]"},
		// negative fail
		{"-sibling(X, Y)", "{X:john, Y:sue}", "[]"},
	}

	for _, test := range tests {

		goal := parser.CreateRelation(test.goal)
		binding := parser.CreateBinding(test.binding)

		resultBindings := solver.SolveRelationSet(mentalese.RelationSet{ goal }, mentalese.Bindings{ binding }).String()

		if !log.IsOk() {
			t.Errorf(log.String())
		}

		if resultBindings != test.resultBindings {
			t.Errorf("SolveRuleBase: got %v, want %v", resultBindings, test.resultBindings)
		}
	}
}
