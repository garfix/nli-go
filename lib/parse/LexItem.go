package parse

import "nli-go/lib/mentalese"

type LexItem struct {
	Form              string
	IsRegExp          bool
	PartOfSpeech      string
	RelationTemplates mentalese.RelationSet
}
