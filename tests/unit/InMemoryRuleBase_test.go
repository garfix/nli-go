package tests

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"testing"
)

func TestInMemoryRuleBase(t *testing.T) {

	parser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog(false)
	matcher := mentalese.NewRelationMatcher(log)
	dialogContext := central.NewDialogContext(matcher, log)
	solver := central.NewProblemSolver(matcher, dialogContext, log)
	facts := parser.CreateRelationSet(`[
		parent(john, jack)
		parent(james, jack)
		parent(vince, bob)
		parent(pat, bob)
		parent(sue, cyrill)
	]`)
	stats := mentalese.DbStats{}
	entities := mentalese.Entities{}
	ds2db := parser.CreateTransformations(`[
		parent(A, B) => parent(A, B);
	]`)
	ds2dbWrite := parser.CreateTransformations(`[
	]`)
	factBase := knowledge.NewInMemoryFactBase("memory", facts, matcher, ds2db, ds2dbWrite, stats, entities, log)
	solver.AddFactBase(factBase)
	rules := parser.CreateRules(`[
		sibling(A, B) :- parent(A, C) parent(B, C);
	]`)
	ruleBase := knowledge.NewInMemoryRuleBase("mem", rules, log)
	keyCabinet := &mentalese.KeyCabinet{}

	tests := []struct {
		goal           string
		binding        string
		resultBindings string
	}{
		//{"sibling(X, Y)", "{}", "[{X:john, Y:john} {X:john, Y:james} {X:james, Y:john} {X:james, Y:james} {X:vince, Y:vince} {X:vince, Y:pat} {X:pat, Y:vince} {X:pat, Y:pat} {X:sue, Y:sue}], want [{X:john, Y:john} {X:john, Y:james} {X:james, Y:john} {X:james, Y:james} {X:vince, Y:vince} {X:vince, Y:pat} {X:pat, Y:vince} {X:pat, Y:pat} {X:sue, Y:sue}], want [{X:john, Y:james} {X:john, Y:john} {X:james, Y:james}]"},
		{"sibling(X, Y)", "{X:john, Z:sue}", "[{X:john, Y:john, Z:sue} {X:john, Y:james, Z:sue}]"},
		{"sibling(X, Y)", "{X:bob}", "[]"},
	}

	for _, test := range tests {

		goal := parser.CreateRelation(test.goal)
		binding := parser.CreateBinding(test.binding)

		resultBindings := solver.SolveSingleRelationSingleBindingSingleRuleBase(goal, keyCabinet, binding, ruleBase).String()

		if resultBindings != test.resultBindings {
			t.Errorf("SolveRuleBase: got %v, want %v", resultBindings, test.resultBindings)
		}
	}
}
