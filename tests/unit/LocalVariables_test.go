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

func TestLocalVariables(t *testing.T) {

	parser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog(false)
	matcher := mentalese.NewRelationMatcher(log)
	dialogContext := central.NewDialogContext()
	meta := mentalese.NewMeta()
	solver := central.NewProblemSolver(matcher, dialogContext, log)
	facts := parser.CreateRelationSet(`
		parent(john, jack)
		parent(james, jack)
		parent(vince, bob)
		parent(pat, bob)
		parent(sue, cyrill)
		sibling(alice, bob)
	`)
	readMap := parser.CreateRules(`[
		parent(A, B) :- parent(A, B);
		sibling(A, B) :- sibling(A, B);
	]`)
	writeMap := parser.CreateRules(`[]`)
	factBase := knowledge.NewInMemoryFactBase("memory", facts, matcher, readMap, writeMap, log)
	solver.AddFactBase(factBase)
	functionBase := knowledge.NewSystemFunctionBase("function", log)
	solver.AddFunctionBase(functionBase)
	nestedBase := nested.NewSystemNestedStructureBase(solver, dialogContext, meta, log)
	solver.AddNestedStructureBase(nestedBase)
	rules := parser.CreateRules(`[
		pow(Base, Number, Pow) :- 
			go:let(Result, 1)
			go:range_foreach(1, Number, _,
				go:multiply(Result, Base, Result)
			)
			go:unify(Pow, Result);	

		first(In, Out) :-
			go:let(X, In)
			go:let(Y, 13)
			times_three(X, Y)
			go:add(Y, 1, X)
			go:unify(Out, X)
		;

		times_three(In, Out) :-
			go:let(X, 3)
			go:multiply(X, In, Out)
		;
	
	]`)
	ruleBase := knowledge.NewInMemoryRuleBase("mem", rules, log)
	solver.AddRuleBase(ruleBase)

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

		resultBindings := solver.SolveRelationSet(mentalese.RelationSet{ goal }, mentalese.Bindings{ binding }).String()

		if !log.IsOk() {
			t.Errorf(log.String())
		}

		if resultBindings != test.resultBindings {
			t.Errorf("SolveRuleBase: got %v, want %v", resultBindings, test.resultBindings)
			t.Errorf(log.String())
		}
	}
}
