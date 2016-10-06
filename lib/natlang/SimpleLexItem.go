package natlang

import "nli-go/lib/mentalese"

type SimpleLexItem struct {
	Form string
	PartOfSpeech      string
	RelationTemplates []mentalese.SimpleRelation
}

