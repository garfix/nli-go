package example2

type simpleLexicon struct {
	lexItems map[string][]SimpleLexItem
}

func NewSimpleLexicon(lexItems map[string][]SimpleLexItem) *simpleLexicon {
	return &simpleLexicon{lexItems: lexItems}
}

func (lexicon *simpleLexicon) GetLexItem(word string, partOfSpeech string) (SimpleLexItem, bool) {
	lexItems, found := lexicon.lexItems[word]

	if found {
		for _, lexItem := range lexItems {
			if lexItem.PartOfSpeech == partOfSpeech {
				return lexItem, true
			}
		}
	}

	return SimpleLexItem{}, false
}
