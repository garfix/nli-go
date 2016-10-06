package natlang

type SimpleLexicon struct {
	lexItems map[string][]SimpleLexItem
}

func NewSimpleLexicon() *SimpleLexicon {
	return &SimpleLexicon{lexItems: map[string][]SimpleLexItem{}}
}

func (lexicon *SimpleLexicon) AddLexItem(lexItem SimpleLexItem) {
	form := lexItem.Form
	_, found := lexicon.lexItems[form]
	if !found {
		lexicon.lexItems[form] = []SimpleLexItem{}
	}
	lexicon.lexItems[form] = append(lexicon.lexItems[form], lexItem)
}

func (lexicon *SimpleLexicon) GetLexItem(word string, partOfSpeech string) (SimpleLexItem, bool) {
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
