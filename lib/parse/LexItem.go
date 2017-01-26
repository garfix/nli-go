package parse

import "nli-go/lib/mentalese"

type LexItem struct {
	Form              string
	PartOfSpeech      string
	RelationTemplates mentalese.RelationSet
}
