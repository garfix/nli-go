package parse

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse/morphology"
)

type MorphologicalAnalyzer struct {
	segmenter    *morphology.Segmenter
	parsingRules *mentalese.GrammarRules
	parser       *EarleyParser
	relationizer *Relationizer
	log          *common.SystemLog
}

func NewMorphologicalAnalyzer(parsingRules *mentalese.GrammarRules, segmenter *morphology.Segmenter, parser *EarleyParser, relationizer *Relationizer, log *common.SystemLog) *MorphologicalAnalyzer {
	return &MorphologicalAnalyzer{
		segmenter:    segmenter,
		parsingRules: parsingRules,
		parser:       parser,
		relationizer: relationizer,
		log:          log,
	}
}

func (morph *MorphologicalAnalyzer) Analyse(word string, lexicalCategory string, variables []string) (mentalese.RelationSet, bool) {

	sense := mentalese.RelationSet{}

	segments := morph.segmenter.Segment(word, lexicalCategory, 0)
	if len(segments) == 0 {
		return sense, false
	}

	trees, _ := morph.parser.Parse(segments, lexicalCategory, variables)
	if len(trees) == 0 {
		return sense, false
	}

	// keep just the first tree, for now
	sense = morph.relationizer.Relationize(trees[0], variables)

	return sense, true
}
