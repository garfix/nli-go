package parse

import "nli-go/lib/mentalese"

type SimpleLexItem struct {
	Form string
	PartOfSpeech      string
	RelationTemplates []mentalese.SimpleRelation
}

