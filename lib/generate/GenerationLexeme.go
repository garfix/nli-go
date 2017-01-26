package generate

import "nli-go/lib/mentalese"

type GenerationLexeme struct {
	Form         string
	IsRegExp     bool
	PartOfSpeech string
	Condition    mentalese.RelationSet
}

