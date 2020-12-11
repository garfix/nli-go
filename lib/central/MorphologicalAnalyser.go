package central

import "nli-go/lib/mentalese"

type MorphologicalAnalyser struct {

}

func NewMorphologicalAnalyser() *MorphologicalAnalyser {
	return &MorphologicalAnalyser{}
}

func (morph *MorphologicalAnalyser) Analyse(word string, lexicalCategory string, variables []string) (mentalese.RelationSet, bool) {

	sense := mentalese.RelationSet{}

	// segment
	// parse
	// relationize
	return sense, false
}