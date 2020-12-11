package tests

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/morphology"
	"strings"
	"testing"
)

func TestSegmenter(t *testing.T) {

	log := common.NewSystemLog()
	parser := importer.NewInternalGrammarParser()
	segmenter := morphology.NewSegmenter()

	tests := []struct {
		input      string
		cat     string
		variables string
		want string
		ok bool
	}{
		{"biggest", "super", "E1", "height(E1, A1) order(A1, desc) first()", true},
	}

	for _, test := range tests {

		want := parser.CreateRelationSet(test.want)
		variables := strings.Split(test.variables, ",")

		result, ok := segmenter.Analyse(test.input, test.cat, variables)

		if result.String() != want.String() || ok != test.ok {
			t.Errorf("call %v: got %v, want %v", test.input, result.String(), want.String())
		}
	}

	if len(log.GetErrors()) > 0 {
		t.Errorf("errors: %v", log.String())
	}
}

func TestMorphologicalAnalyser(t *testing.T) {

	log := common.NewSystemLog()
	parser := importer.NewInternalGrammarParser()
	analyser := central.NewMorphologicalAnalyser()

	tests := []struct {
		input      string
		cat     string
		variables string
		want string
		ok bool
	}{
		{"biggest", "super", "E1", "height(E1, A1) order(A1, desc) first()", true},
	}

	for _, test := range tests {

		want := parser.CreateRelationSet(test.want)
		variables := strings.Split(test.variables, ",")

		result, ok := analyser.Analyse(test.input, test.cat, variables)

		if result.String() != want.String() || ok != test.ok {
			t.Errorf("call %v: got %v, want %v", test.input, result.String(), want.String())
		}
	}

	if len(log.GetErrors()) > 0 {
		t.Errorf("errors: %v", log.String())
	}
}
