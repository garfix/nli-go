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
			"[have_child(S1, S2) isa(S2, man) quantification(Q1, R1, S1) isa(R1, parent) isa(Q1, all) quantification(Q2, R2, S2) isa(R2, child) isa(Q2, 2)]",
			"[quant(Q1, [isa(Q1, all)], R1, [isa(R1, parent)], [quant(Q2, [isa(Q2, 2)], R2, [isa(R2, child)], [have_child(R1, R2) isa(R2, man)])])]",
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
