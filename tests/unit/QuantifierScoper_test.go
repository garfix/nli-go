package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"testing"
)

func TestQuantifierScoper(t *testing.T) {

	tests := []struct {
		input string
		want  string
	}{
		{
			"[have_child(S1, O1) quantification(S1, [ isa(S1, parent) ], D1, [ isa(D1, all) ]) quantification(O1, [ isa(O1, child) ], D2, [ isa(D2, 2) ])]",
			"[quant(S1, [isa(S1, parent)], D1, [isa(D1, all)], [quant(O1, [isa(O1, child)], D2, [isa(D2, 2)], [have_child(S1, O1)])])]",
		},
		{
			"[quantification(O1, [ isa(O1, child) ], D2, [ isa(D2, 2) ]) quantification(S1, [ isa(S1, parent) ], D1, [ isa(D1, all) ]) have_child(S1, O1)]",
			"[quant(S1, [isa(S1, parent)], D1, [isa(D1, all)], [quant(O1, [isa(O1, child)], D2, [isa(D2, 2)], [have_child(S1, O1)])])]",
		},
	}

	log := common.NewSystemLog(false)
	internalGrammarParser := importer.NewInternalGrammarParser()
	quantifierScoper := mentalese.NewQuantifierScoper(log)

	for _, test := range tests {

		input := internalGrammarParser.CreateRelationSet(test.input)
		result := quantifierScoper.Scope(input)

		if result.String() != test.want {
			t.Errorf("got %s, want %s", result.String(), test.want)
		}
	}
}
