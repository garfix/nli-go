package example3

type simpleLexicon struct {
	lexItems map[string][]SimpleLexItem
}

func NewSimpleLexicon() *simpleLexicon{
	return &simpleLexicon{lexItems: map[string][]SimpleLexItem{}}
}

func (lexicon *simpleLexicon) AddLexItem(lexItem SimpleLexItem) {
	form := lexItem.Form
	_, found := lexicon.lexItems[form]
	if !found {
		lexicon.lexItems[form] = []SimpleLexItem{}
	}
	lexicon.lexItems[form] = append(lexicon.lexItems[form], lexItem)
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
