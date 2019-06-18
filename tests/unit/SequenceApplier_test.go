package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"testing"
)

func TestSequenceApplier(t *testing.T) {

	tests := []struct {
		input string
		want  string
	}{
		{
			"[abc(P1) sequence(P1, C, P2) def(P2, X) ghi(X, Y) jkl(S)]",
			"[seq([abc(P1)], C, [def(P2, X) ghi(X, Y)]) jkl(S)]",
		},
// todo: test nested sequences!
	}

	log := common.NewSystemLog(false)
	internalGrammarParser := importer.NewInternalGrammarParser()
	sequenceApplier := mentalese.NewSequenceApplier(log)

	for _, test := range tests {

		input := internalGrammarParser.CreateRelationSet(test.input)
		result := sequenceApplier.ApplySequences(input)

		if result.String() != test.want {
			t.Errorf("got %s, want %s", result.String(), test.want)
		}
	}
}
