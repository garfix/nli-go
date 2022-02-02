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

func TestLocalVariables(t *testing.T) {

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
		sibling(alice, bob)
	`)
	readMap := parser.CreateRules(`
		parent(A, B) :- parent(A, B);
		sibling(A, B) :- sibling(A, B);
	`)
	writeMap := []mentalese.Rule{}
	factBase := knowledge.NewInMemoryFactBase("memory", facts, matcher, readMap, writeMap, nil, log)
	solver.AddFactBase(factBase)
	functionBase := knowledge.NewSystemFunctionBase("function", meta, log)
	solver.AddFunctionBase(functionBase)
	deicticCenter := central.NewDeicticCenter()
	discourseEntities := mentalese.NewBinding()
	processList := central.NewProcessList()
	dialogContext := central.NewDialogContext(nil, deicticCenter, processList, variableGenerator, &discourseEntities)
	nestedBase := function.NewSystemSolverFunctionBase(dialogContext, meta, log)
	solver.AddSolverFunctionBase(nestedBase)
	rules := parser.CreateRules(`
		pow(Base, Number, Pow) :- 
			[:Result := 1]
			go:range_foreach(1, Number, _,
				[:Result := [:Result * Base]]
			)
			[Pow := :Result];	

		first(In, Out) :-
			[:X := In]
			[:Y := 13]
			times_three(:X, :Y)
			[:X := [:Y + 1]]
			[Out := :X]
		;

		times_three(In, Out) :-
			[:X := 3]
			[Out := [:X * In]]
		;
	
	`)
	ruleBase := knowledge.NewInMemoryRuleBase("mem", rules, []string{}, nil, log)
	solver.AddRuleBase(ruleBase)
	solver.Reindex()
	runner := central.NewProcessRunner(solver, log)

	tests := []struct {
		goal           string
		binding        string
		resultBindings string
	}{
		{"pow(2, 3, Pow)", "{}", "[{Pow:8}]"},
		{"first(5, Result)", "{}", "[{Result:16}]"},
	}

	for _, test := range tests {

		log.Clear()

		goal := parser.CreateRelation(test.goal)
		binding := parser.CreateBinding(test.binding)

		resultBindings := runner.RunRelationSetWithBindings(mentalese.RelationSet{goal}, mentalese.InitBindingSet(binding)).String()

		if !log.IsOk() {
			t.Errorf(log.String())
		}

		if resultBindings != test.resultBindings {
			t.Errorf("SolveRuleBase: got %v, want %v", resultBindings, test.resultBindings)
			t.Errorf(log.String())
		}
	}
}
