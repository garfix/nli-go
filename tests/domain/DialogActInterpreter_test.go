package tests

import (
	"testing"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
)

func TestDialogActInterpreter(test *testing.T) {

	// relations
	internalGrammarParser := importer.NewInternalGrammarParser()
	relationMatcher := mentalese.NewRelationMatcher()

	// who did Kurt Cobain marry?
	// Note: this representation is rubbish :) I will get to that later
	sense := internalGrammarParser.CreateRelationSet(`[
		predication(S1, marry)
		object(S1, E2)
		subject(S1, who)
		name(E1, 'Kurt Cobain')
	]`)

	// interpret dialog act (via transformations)
	analysis := internalGrammarParser.CreateTransformations(`[
		subject(S1, who) => dialog_act(S1, info_request);
	]`)

	matcher := mentalese.NewRelationMatcher()
	analyser := mentalese.NewRelationTransformer(matcher)
	dialogActs := analyser.Extract(analysis, sense)

	infoRequest := internalGrammarParser.CreateRelationSet(`[ dialog_act(S, info_request) ]`)

	binding, _, _ := relationMatcher.MatchSequenceToSet(infoRequest, dialogActs, mentalese.Binding{})
	if binding[0].String() != "{S:S1}" {
		test.Error("No match")
	}
}
