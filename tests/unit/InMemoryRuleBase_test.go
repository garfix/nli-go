package tests

import (
	"nli-go/lib/central"
	"nli-go/lib/central/goal"
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
	solver := central.NewProblemSolverAsync(matcher, variableGenerator, log)
	facts := parser.CreateRelationSet(`
		parent(john, jack)
		parent(james, jack)
		parent(vince, bob)
		parent(pat, bob)
		parent(sue, cyrill)
		-sibling(alice, bob)
	`)
	readMap := parser.CreateRules(`
		parent(A, B) :- parent(A, B);
		-sibling(A, B) :- -sibling(A, B);
	`)
	writeMap := []mentalese.Rule{}
	factBase := knowledge.NewInMemoryFactBase("memory", facts, matcher, readMap, writeMap, nil, log)
	solver.AddFactBase(factBase)
	functionBase := knowledge.NewSystemFunctionBase("function", meta, log)
	solver.AddFunctionBase(functionBase)
	anaphoraQueue := central.NewAnaphoraQueue()
	deicticCenter := central.NewDeicticCenter()
	discourseEntities := mentalese.NewBinding()
	processList := goal.NewProcessList()
	dialogContext := central.NewDialogContext(nil, anaphoraQueue, deicticCenter, processList, variableGenerator, &discourseEntities)
	nestedBase := function.NewSystemSolverFunctionBase(dialogContext, meta, log)
	solver.AddSolverFunctionBase(nestedBase)
	runner := central.NewProcessRunner(solver, log)
	rules := parser.CreateRules(`
		sibling(A, B) :- parent(A, C) parent(B, C) go:not( -sibling(A, B) );
		-sibling(A, B) :- go:equals(A, B);
	`)
	ruleBase := knowledge.NewInMemoryRuleBase("mem", rules, []string{}, nil, log)
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
	}

	for _, test := range tests {

		goal := parser.CreateRelation(test.goal)
		binding := parser.CreateBinding(test.binding)

		resultBindings := runner.RunRelationSetWithBindings(mentalese.RelationSet{goal}, mentalese.InitBindingSet(binding)).String()

		if !log.IsOk() {
			t.Errorf(log.String())
		}

		if resultBindings != test.resultBindings {
			t.Errorf("SolveRuleBase: got %v, want %v", resultBindings, test.resultBindings)
		}
	}
}
