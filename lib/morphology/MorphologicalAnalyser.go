package morphology

import (
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
)

type MorphologicalAnalyser struct {
	segmenter *Segmenter
	parsingRules *parse.GrammarRules
	parser api.Parser
}

func NewMorphologicalAnalyser(parser api.Parser, segmenter *Segmenter, parsingRules *parse.GrammarRules, log *common.SystemLog) *MorphologicalAnalyser {
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