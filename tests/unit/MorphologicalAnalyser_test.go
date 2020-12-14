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

	segmentationRules := parser.CreateSegmentationRulesAndCharacterClasses(`
		vowel: ['a', 'e', 'i', 'o', 'u', 'y']
		consonant: ['b', 'c', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'p', 'q', 'r', 's', 't', 'v', 'w', 'x', 'z']
		comp: '*{consonant1}{consonant1}er' -> adj: '*{consonant1}', suffix: 'er'
		super: '*{consonant1}{consonant1}est' -> adj: '*{consonant1}', suffix: 'est'
		super: '*est' -> adj: '*e', suffix: 'est'
		super: '*est' -> adj: '*', suffix: 'est'
		adj: 'high'
		adj: 'big'
		adj: 'little'
		suffix: 'est'
		suffix: 'er'
	`)

	tests := []struct {
		input string
		cat   string
		want  string
	}{
		{"highest", "super", "high est"},
		{"biggest", "super", "big est"},
		{"bigger", "comp", "big er"},
		{"littlest", "super", "little est"},
	}

	for _, test := range tests {

		result := strings.Join(segmenter.Segment(segmentationRules, test.input, test.cat), " ")

		if result != test.want {
			t.Errorf("call %v: got %v, want %v", test.input, result, test.want)
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
