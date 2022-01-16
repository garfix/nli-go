package tests

import (
	"fmt"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"testing"
)

func TestInternalGrammarParser(t *testing.T) {

	parser := importer.NewInternalGrammarParser()

	tests := []string{
		"determiner(E, some, D, some)",
	}

	for _, test := range tests {
		result := parser.CreateRelation(test)
		if result.String() != test {
			t.Errorf("got %s, want %s", result.String(), test)
		}
	}

	// =====================================================

	grammar := parser.CreateGrammarRules("{ rule: s(P) -> np(E) vp(P),         sense: subject(P, E) }")

	rules := grammar.FindRules("s", 1)
	if len(rules) == 0 {
		t.Error("No rules found")
	}

	if rules[0].GetAntecedent() != "s" {
		t.Error(fmt.Printf("Error in rule: %s", rules[0].GetAntecedent()))
	}
	if rules[0].GetConsequent(0) != "np" {
		t.Error(fmt.Printf("Error in rule: %s", rules[0].GetConsequent(0)))
	}
	if rules[0].GetAntecedentVariables()[0] != "P" {
		t.Error(fmt.Printf("Error in rule: %s", rules[0].GetAntecedentVariables()))
	}
	if rules[0].GetConsequentVariables(0)[0] != "E" {
		t.Error(fmt.Printf("Error in rule: %s", rules[0].GetConsequentVariables(0)))
	}
	if len(rules[0].Sense) != 1 {
		t.Error(fmt.Printf("Error in number of sense relations: %d", len(rules[0].Sense)))
	}

	grammar = parser.CreateGrammarRules(
		"{ rule: s(P) -> np(E) vp(P),    sense: subject(P, E) }" +
			"{ rule: np(P) -> nbar(E) }")

	rules = grammar.FindRules("s", 1)
	if len(rules) != 1 {
		t.Error("No rules found")
	}
	rules = grammar.FindRules("np", 1)
	if len(rules) != 1 {
		t.Error("No rules found")
	}

	grammar = &mentalese.GrammarRules{}

	parser.CreateRelationSet("assert(at(5, 3))")
	parser.CreateRelationSet("learn(own(X, Y) :- fish(Y))")
	parser.CreateRelationSet("sort([])")
	parser.CreateRelationSet("sort([5])")
	parser.CreateRelationSet("sort([5,2,3,1])")
	parser.CreateRelationSet("[A = 1]")
	parser.CreateRelationSet("[A = :B]")
	parser.CreateRelationSet("[A == B]")
	parser.CreateRelationSet("[A != B]")
	parser.CreateGrammarRules("{ rule: a(P) -> b(P), ellipsis: [prev_sentence] }")
	parser.CreateGrammarRules("{ rule: a(P) -> b(P), ellipsis: np(E) }")
	parser.CreateGrammarRules("{ rule: a(P) -> b(P), ellipsis: .. }")
	parser.CreateGrammarRules("{ rule: a(P) -> b(P), ellipsis: ..s(P) }")
	parser.CreateGrammarRules("{ rule: a(P) -> b(P), ellipsis: - }")
	parser.CreateGrammarRules("{ rule: a(P) -> b(P), ellipsis: -np(E) }")
	parser.CreateGrammarRules("{ rule: a(P) -> b(P), ellipsis: + }")
	parser.CreateGrammarRules("{ rule: a(P) -> b(P), ellipsis: +np(E) }")
	parser.CreateGrammarRules("{ rule: a(P) -> b(P), ellipsis: +- }")
	parser.CreateGrammarRules("{ rule: a(P) -> b(P), ellipsis: +-np(E) }")
	parser.CreateGrammarRules("{ rule: a(P) -> b(P), ellipsis: ..np(E)/vp(P)/+noun(E)/-noun(E)/+-noun(E) }")
	parser.CreateGrammarRules("{ rule: a(P) -> b(P), ellipsis: vp(P)//noun(E) }")
	parser.CreateGrammarRules("{ rule: a(P) -> b(P), tag: function(P, subject) }")
	parser.CreateGrammarRules("{ rule: a(P) -> b(P), tag: number(P, singular) person(P, 3) }")
	parser.CreateRules("describe_event(P1, DescSet) :- pick_up(P1, Subject, Object);")
	parser.CreateRules("describe_event(P1, DescSet) :- if a(1) a(2) then b(1) b(2) end;")
	parser.CreateRules("describe_event(P1, DescSet) :- if a(1) a(2) then b(1) b(2) else c(1) c(2) end;")
	parser.CreateRules("a(X) :- return;")
	parser.CreateRules("a(X) :- break;")
	parser.CreateRules("a(X) :- cancel;")
	parser.CreateRules("a(X) :- fail;")

	set := parser.CreateRelationSet("quant_foreach($np, quant_foreach($np2, none))")
	if set.String() != "quant_foreach(go_sem(np, 1), quant_foreach(go_sem(np, 2), none))" {
		t.Error(set.String())
	}

	set2 := parser.CreateRelationSet("unify(A, parent(a, b)) {{ A }}")
	if set2.String() != "unify(A, parent(a, b)) $go$_include_relations(A)" {
		t.Error(set2.String())
	}
}
