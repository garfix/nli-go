package example1

type simpleLexicon struct {
    lexItems map[string][]string
}

func NewSimpleLexicon(lexItems map[string][]string) *simpleLexicon {
    return &simpleLexicon{lexItems:lexItems}
}

func (lexicon *simpleLexicon) CheckPartOfSpeech(word string, parOfSpeech string) bool {
    partsOfSpeech, found := lexicon.lexItems[word]

    if found {
        for i := 0; i < len(partsOfSpeech); i++ {
            pos := partsOfSpeech[i]
            if pos == parOfSpeech {
                return true
            }
        }
    }

    return false
}