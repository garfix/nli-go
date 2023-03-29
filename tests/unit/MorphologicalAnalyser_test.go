package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse/morphology"
	"strings"
	"testing"
)

func TestSegmenter(t *testing.T) {

	log := common.NewSystemLog()
	parser := importer.NewInternalGrammarParser()

	segmentationRules := parser.CreateSegmentationRules(`
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
	segmenter := morphology.NewSegmenter(segmentationRules, mentalese.NewGrammarRules())

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

		result := strings.Join(segmenter.Segment(test.input, test.cat, 0), " ")

		if result != test.want {
			t.Errorf("call %v: got %v, want %v", test.input, result, test.want)
		}
	}

	if len(log.GetErrors()) > 0 {
		t.Errorf("errors: %v", log.String())
	}
}
