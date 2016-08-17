package example2

type simpleLexicon struct {
    lexItems map[string][]SimpleLexItem
}

func NewSimpleLexicon(lexItems map[string][]SimpleLexItem) *simpleLexicon {
    return &simpleLexicon{lexItems:lexItems}
}

func (lexicon *simpleLexicon) CheckPartOfSpeech(word string, parOfSpeech string) bool {
    lexItems, found := lexicon.lexItems[word]

    if found {
        for i := 0; i < len(lexItems); i++ {
            lexItem := lexItems[i]
            if lexItem.PartOfSpeech == parOfSpeech {
                return true
            }
        }
    }

    return false
}

func (lexicon *simpleLexicon) GetLexItem(word string, partOfSpeech string) (SimpleLexItem, bool) {
    lexItems, found := lexicon.lexItems[word]

    if found {
        for i := 0; i < len(lexItems); i++ {
            lexItem := lexItems[i]
            if lexItem.PartOfSpeech == partOfSpeech {
                return lexItem, true
            }
        }
    }

    return SimpleLexItem{}, false
}
