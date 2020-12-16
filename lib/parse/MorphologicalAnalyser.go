package parse

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse/morphology"
)

type MorphologicalAnalyser struct {
	segmenter *morphology.Segmenter
	parsingRules *GrammarRules
	parser *EarleyParser
}

func NewMorphologicalAnalyser(parser *EarleyParser, segmenter *morphology.Segmenter, parsingRules *GrammarRules, log *common.SystemLog) *MorphologicalAnalyser {
	return &MorphologicalAnalyser{
		segmenter: segmenter,
		parsingRules: parsingRules,
		parser: parser,
	}
}

func (morph *MorphologicalAnalyser) Analyse(word string, lexicalCategory string, variables []string) (mentalese.RelationSet, bool) {

	sense := mentalese.RelationSet{}

	// segment
	// parse
	// relationize
	return sense, false
}