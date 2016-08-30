package example3

import "nli-go/lib/example2"

type simpleLexicon struct {
	lexItems map[string][]example2.SimpleLexItem
}

func NewSimpleLexicon() *simpleLexicon{
	return &simpleLexicon{lexItems: map[string][]example2.SimpleLexItem{}}
}

func (lexicon *simpleLexicon) AddLexItem(lexItem example2.SimpleLexItem) {
	form := lexItem.Form
	_, found := lexicon.lexItems[form]
	if !found {
		lexicon.lexItems[form] = []example2.SimpleLexItem{}
	}
	lexicon.lexItems[form] = append(lexicon.lexItems[form], lexItem)
}

func (lexicon *simpleLexicon) GetLexItem(word string, partOfSpeech string) (example2.SimpleLexItem, bool) {
	lexItems, found := lexicon.lexItems[word]

	if found {
		for _, lexItem := range lexItems {
			if lexItem.PartOfSpeech == partOfSpeech {
				return lexItem, true
			}
		}
	}

	return example2.SimpleLexItem{}, false
}
