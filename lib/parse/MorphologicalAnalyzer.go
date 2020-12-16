package parse

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse/morphology"
)

type MorphologicalAnalyzer struct {
	segmenter *morphology.Segmenter
	parsingRules *GrammarRules
	parser *EarleyParser
	relationizer *Relationizer
	log *common.SystemLog
}

func NewMorphologicalAnalyzer(parsingRules *GrammarRules, segmenter *morphology.Segmenter, parser *EarleyParser, relationizer *Relationizer, log *common.SystemLog) *MorphologicalAnalyzer {
	return &MorphologicalAnalyzer{
		segmenter: segmenter,
		parsingRules: parsingRules,
		parser: parser,
		relationizer: relationizer,
		log: log,
	}
}

func (morph *MorphologicalAnalyzer) Analyse(word string, lexicalCategory string, variables []string) (mentalese.RelationSet, bool) {

	sense := mentalese.RelationSet{}

//fmt.Println(word, lexicalCategory)

	segments := morph.segmenter.Segment(word, lexicalCategory)
	if len(segments) == 0 {
		return sense, false
	}

	trees := morph.parser.Parse(segments)
	if trees == nil {
		return sense, false
	}

	// keep just the first tree, for now
	sense, _ = morph.relationizer.Relationize(trees[0])

	return sense, len(sense) > 0
}