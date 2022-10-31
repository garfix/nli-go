package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"testing"
)

func TestFillerStack(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()

	grammarRules := internalGrammarParser.CreateGrammarRules(`

		{ rule: s(P1) -> 'which' np(E1) dep_vp(P1, E1),						sense: which(E1) check($np, $dep_vp) }
		{ rule: np(E1) -> nbar(E1), 										sense: quant(_, some(_), E1, $nbar) }
		{ rule: nbar(E) -> noun(E) }
		{ rule: dep_vp(P1, E1) -> be(_) np(E2) advp(P1) vp(P1, E1, E2), 	sense: check($np, $advp $vp) }
		{ rule: np(E1) -> qp(Q1) nbar(E1), 									sense: quant(Q1, $qp, E1, $nbar) }
		{ rule: advp(P1) -> adverb(P1) }
		{ rule: vp(P1, E1, E2) -> 'to' 'take' 'from', 						sense: take_from(P1, E2, E1)  }
		{ rule: be(P1) -> 'were' }
		{ rule: qp(Q1) -> 'the', 											sense: the(Q1) }
		{ rule: adverb(A1) -> 'easiest', 									sense: easiest(A1) }
		{ rule: noun(E1) -> 'babies', 										sense: baby(E1) }
		{ rule: noun(E1) -> 'toys',											sense: toy(E1) }
	`)

	log := common.NewSystemLog()
	parser := parse.NewParser(grammarRules, log)

	variableGenerator := mentalese.NewVariableGenerator()
	dialogizer := parse.NewDialogizer(variableGenerator)
	relationizer := parse.NewRelationizer(variableGenerator, log)

	parseTrees := parser.Parse([]string{"Which", "babies", "were", "the", "toys", "easiest", "to", "take", "from"}, "s", []string{"S"})

	if len(parseTrees) == 0 {
		t.Error(log.String())
		return
	}

	tree := dialogizer.Dialogize(&parseTrees[0])
	result := relationizer.Relationize(*tree, []string{"S"})

	want := "which(E$1) check(quant(_, some(_), E$1, baby(E$1)), check(quant(Q$1, the(Q$1), E$2, toy(E$2)), easiest(Sentence$1) take_from(Sentence$1, E$2, E$1)))"
	if result.String() != want {
		t.Errorf("got %s, want %s", result.String(), want)
	}
}
