package tests

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse/earley"
	"testing"
)

func TestFillerStack(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()

	grammarRules := internalGrammarParser.CreateGrammarRules(`[

		{ rule: s(P1) -> 'which' np(E1) dep_vp(P1, E1),						sense: which(E1) quant_check($np, $dep_vp) }
		{ rule: np(E1) -> nbar(E1), 										sense: quant(_, some(_), E1, $nbar) }
		{ rule: nbar(E) -> noun(E) }
		{ rule: dep_vp(P1, E1) -> be(_) np(E2) advp(P1) vp(P1, E1, E2), 	sense: quant_check($np, $advp $vp) }
		{ rule: np(E1) -> qp(Q1) nbar(E1), 									sense: quant(Q1, $qp, E1, $nbar) }
		{ rule: advp(P1) -> adverb(P1) }
		{ rule: vp(P1, E1, E2) -> 'to' 'take' 'from', 						sense: take_from(P1, E2, E1)  }
		{ rule: be(P1) -> 'were' }
		{ rule: qp(Q1) -> 'the', 											sense: the(Q1) }
		{ rule: adverb(A1) -> 'easiest', 									sense: easiest(A1) }
		{ rule: noun(E1) -> 'babies', 										sense: baby(E1) }
		{ rule: noun(E1) -> 'toys',											sense: toy(E1) }
	]`)

	log := common.NewSystemLog(true)

	matcher := mentalese.NewRelationMatcher(log)
	dialogContext := central.NewDialogContext()
	meta := mentalese.NewMeta()
	solver := central.NewProblemSolver(matcher, dialogContext, log)
	nameResolver := central.NewNameResolver(solver, meta, matcher, log, dialogContext)

	parser := earley.NewParser(nameResolver, meta, log)

	relationizer := earley.NewRelationizer(log)

	parseTrees := parser.Parse(grammarRules, []string{"Which", "babies", "were", "the", "toys", "easiest", "to", "take", "from"})

	if len(parseTrees) == 0 {
		t.Error(log.String())
		return
	}

	result, _ := relationizer.Relationize(parseTrees[0], nameResolver)

	want := "which(E5) quant_check(quant(_, some(_), E5, baby(E5)), quant_check(quant(Q5, the(Q5), E6, toy(E6)), easiest(S5) take_from(S5, E6, E5)))"
	if result.String() != want {
		t.Errorf("got %s, want %s", result.String(), want)
	}
}
