package tests

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/knowledge/function"
	"nli-go/lib/mentalese"
	"testing"
)

func TestInMemoryRuleBase(t *testing.T) {

	parser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog()
	matcher := central.NewRelationMatcher(log)
	meta := mentalese.NewMeta()
	variableGenerator := mentalese.NewVariableGenerator()
	solver := central.NewProblemSolver(matcher, variableGenerator, log)
	facts := parser.CreateRelationSet(`
		parent(john, jack)
		parent(james, jack)
		parent(vince, bob)
		parent(pat, bob)
		parent(sue, cyrill)
		age(pat, 51)
		age(sue, 49)
		-sibling(alice, bob)
	`)
	readMap := parser.CreateRules(`
		parent(A, B) :- parent(A, B);
		-sibling(A, B) :- -sibling(A, B);
		age(X, Y) :- age(X, Y);
	`)
	writeMap := []mentalese.Rule{}
	factBase := knowledge.NewInMemoryFactBase("memory", facts, matcher, readMap, writeMap, log)
	solver.AddFactBase(factBase)
	functionBase := knowledge.NewSystemFunctionBase("function", meta, log)
	solver.AddFunctionBase(functionBase)
	processList := central.NewProcessList()
	dialogContext := central.NewDialogContext(variableGenerator)
	nestedBase := function.NewSystemSolverFunctionBase(dialogContext, meta, log, nil)
	solver.AddSolverFunctionBase(nestedBase)
	runner := central.NewProcessRunner(processList, solver, log)
	rules := parser.CreateRules(`
		sibling(A, B) :- parent(A, C) parent(B, C) go:not( -sibling(A, B) );
		-sibling(A, B) :- [A == B];
		older(A, B) :- [age(A, rv) > age(B, rv)];
		isa(man, animal);
	`)
	ruleBase := knowledge.NewInMemoryRuleBase("mem", rules, []string{}, log)
	solver.AddRuleBase(ruleBase)
	solver.Reindex()

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
		// facts as functions
		{"older(A, B)", "{A:pat, B:sue}", "[{A:pat, B:sue}]"},
		{"older(A, B)", "{B:pat, A: sue}", "[]"},
		// bind by match
		{"isa(man, Type)", "{}", "[{Type:animal}]"},
	}

	for _, test := range tests {

		goal := parser.CreateRelation(test.goal)
		binding := parser.CreateBinding(test.binding)

		resultBindings := runner.RunRelationSetWithBindings(central.NO_RESOURCE, mentalese.RelationSet{goal}, mentalese.InitBindingSet(binding)).String()

		if !log.IsOk() {
			t.Errorf(log.String())
		}

		if resultBindings != test.resultBindings {
			t.Errorf("SolveRuleBase: got %v, want %v", resultBindings, test.resultBindings)
		}
	}
}
