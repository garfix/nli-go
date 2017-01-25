package parse

import (
	"strings"
	"nli-go/lib/mentalese"
)

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

	// proper noun?
	if !found {
		if partOfSpeech == "fullName" || partOfSpeech == "firstName" || partOfSpeech == "lastName" {

		if strings.ToUpper(string(word[0])) == string(word[0]) {
			return LexItem{
				Form: word,
				PartOfSpeech: "name",
				RelationTemplates: []mentalese.Relation{
					{ Predicate: "name", Arguments: []mentalese.Term{
						{ TermType: mentalese.Term_predicateAtom, TermValue: "this" },
						{ TermType: mentalese.Term_stringConstant, TermValue: word },
						{ TermType: mentalese.Term_predicateAtom, TermValue: partOfSpeech }}}}}, true
			}
		}
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
