package parse

import "strings"

type Lexicon struct {
	lexItems map[string][]LexItem
}

func NewLexicon() *Lexicon {
	return &Lexicon{lexItems: map[string][]LexItem{}}
}

func (lexicon *Lexicon) AddLexItem(lexItem LexItem) {
	form := lexItem.Form
	_, found := lexicon.lexItems[form]
	if !found {
		lexicon.lexItems[form] = []LexItem{}
	}
	lexicon.lexItems[form] = append(lexicon.lexItems[form], lexItem)
}

func (lexicon *Lexicon) GetLexItem(word string, partOfSpeech string) (LexItem, bool) {

	// try the word as is
	lexItems, found := lexicon.lexItems[word]

	// try the lowercase version
	if !found {
		lexItems, found = lexicon.lexItems[strings.ToLower(word)]
	}

	if found {
		for _, lexItem := range lexItems {
			if lexItem.PartOfSpeech == partOfSpeech {
				return lexItem, true
			}
		}
	}

	return LexItem{}, false
}
