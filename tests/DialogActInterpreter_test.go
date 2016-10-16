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
	sense, _, _ := internalGrammarParser.CreateRelationSet(`
		predication(S1, marry)
		object(S1, E2)
		subject(S1, who)
		name(E1, 'Kurt Cobain')
	`)

	// interpret dialog act (via transformations)
	analysis, _, _ := internalGrammarParser.CreateTransformations(`
		dialog_act(S1, info_request) :- subject(S1, who)
	`)

	analyser := mentalese.NewRelationTransformer(analysis)
	dialogActs := analyser.Extract(sense)

	infoRequest, _, _ := internalGrammarParser.CreateRelationSet(`dialog_act(S, info_request)`)

	_, binding, _ := relationMatcher.MatchSequenceToSet(infoRequest, dialogActs, mentalese.Binding{})
	if binding.String() != "{S:S1}" {
		test.Error("No match")
	}
}
