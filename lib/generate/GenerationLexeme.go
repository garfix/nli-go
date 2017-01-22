package generate

import "nli-go/lib/mentalese"

type GenerationLexeme struct {
	Form         string
	PartOfSpeech string
	Condition    mentalese.RelationSet
}

