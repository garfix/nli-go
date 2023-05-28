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
	solver := central.NewProblemSolver(matcher, variableGenerator, log)
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
	factBase := knowledge.NewInMemoryFactBase("memory", facts, matcher, readMap, writeMap, log)
	solver.AddFactBase(factBase)
	functionBase := knowledge.NewSystemFunctionBase("function", meta, log)
	solver.AddFunctionBase(functionBase)
	processList := central.NewProcessList()
	dialogContext := central.NewDialogContext(variableGenerator)
	nestedBase := function.NewSystemSolverFunctionBase(dialogContext, meta, log, nil)
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

		break_out(Start, End, Result) :-
			go:range_foreach(Start, End, Index,
				if [Index == 5] then
					[Result := Index]
					break
				end
			)
		;

		return_out(Start, End, Result) :-
			go:list_foreach([1, 2, 3, 4, 5, 6, 7], Index,
				if [Index == 5] then
					[Result := Index]
					return
				end
			)
		;

	`)
	ruleBase := knowledge.NewInMemoryRuleBase("mem", rules, []string{}, log)
	solver.AddRuleBase(ruleBase)
	solver.Reindex()
	runner := central.NewProcessRunner(processList, solver, log)

	tests := []struct {
		goal           string
		binding        string
		resultBindings string
	}{
		//{"break_out(1, 10, Result)", "{}", "[{} {Result:5}]"},
		//{"return_out(1, 10, Result)", "{}", "[{} {Result:5}]"},
		{"pow(2, 3, Pow)", "{}", "[{Pow:8}]"},
		//{"first(5, Result)", "{}", "[{Result:16}]"},
	}

	for _, test := range tests {

		log.Clear()

		goal := parser.CreateRelation(test.goal)
		binding := parser.CreateBinding(test.binding)

		resultBindings := runner.RunRelationSetWithBindings(central.NO_RESOURCE, mentalese.RelationSet{goal}, mentalese.InitBindingSet(binding)).String()

		if !log.IsOk() {
			t.Errorf(log.String())
		}

		if resultBindings != test.resultBindings {
			t.Errorf("SolveRuleBase: got %v, want %v", resultBindings, test.resultBindings)
			t.Errorf(log.String())
		}
	}
}
