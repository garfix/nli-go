package morphology

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"nli-go/lib/parse/earley"
)

type MorphologicalAnalyser struct {
	segmenter *Segmenter
	parsingRules *parse.GrammarRules
	parser *earley.Parser
}

func NewMorphologicalAnalyser(segmenter *Segmenter, parsingRules *parse.GrammarRules, log *common.SystemLog) *MorphologicalAnalyser {
	return &MorphologicalAnalyser{
		segmenter: segmenter,
		parsingRules: parsingRules,
		parser: earley.NewParser(log),
	}
}

func (morph *MorphologicalAnalyser) Analyse(word string, lexicalCategory string, variables []string) (mentalese.RelationSet, bool) {

	sense := mentalese.RelationSet{}

	// segment
	// parse
	// relationize
	return sense, false
}