package api

import "nli-go/lib/mentalese"

type MorphologicalAnalyser interface {
	Analyse(word string, lexicalCategory string, variables []string) (mentalese.RelationSet, bool)
}
