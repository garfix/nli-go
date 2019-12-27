package generate

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type GenerationLexicon struct {
	allItems []GenerationLexeme
	lexemes  map[string][]GenerationLexeme
	matcher  *mentalese.RelationMatcher
	log      *common.SystemLog
}

func NewGenerationLexicon(log *common.SystemLog, matcher *mentalese.RelationMatcher) *GenerationLexicon {
	return &GenerationLexicon{
		allItems: []GenerationLexeme{},
		lexemes:  map[string][]GenerationLexeme{},
		matcher:  matcher,
		log:      log,
	}
}

func (lexicon *GenerationLexicon) AddLexItem(lexItem GenerationLexeme) {

	lexicon.allItems = append(lexicon.allItems, lexItem)

	pos := lexItem.PartOfSpeech
	_, found := lexicon.lexemes[pos]
	if !found {
		lexicon.lexemes[pos] = []GenerationLexeme{}
	}
	lexicon.lexemes[pos] = append(lexicon.lexemes[pos], lexItem)
}

// Example:
// consequent: noun(E1)
func (lexicon *GenerationLexicon) GetLexemeForGeneration(consequent mentalese.Relation, sentenseSense mentalese.RelationSet) (GenerationLexeme, bool) {

	resultLexeme := GenerationLexeme{}

	lexicon.log.StartDebug("GetLexemeForGeneration", consequent)

	if consequent.Predicate == "number" {
		resultLexeme.Form = consequent.Arguments[0].TermValue
		return resultLexeme, true
	}

	if consequent.Predicate == "text" {
		resultLexeme.Form = consequent.Arguments[0].TermValue
		return resultLexeme, true
	}

	partOfSpeech := consequent.Predicate
	applicableLexemeFound := false

	lexemes, found := lexicon.lexemes[partOfSpeech]
	if found {

		binding := mentalese.Binding{"E": consequent.Arguments[0]}

		for _, lexeme := range lexemes {

			_, match := lexicon.matcher.MatchSequenceToSet(lexeme.Condition, sentenseSense, binding)

			if match {
				resultLexeme = lexeme

				applicableLexemeFound = true

				break
			}
		}
	}

	lexicon.log.EndDebug("GetLexemeForGeneration", resultLexeme, applicableLexemeFound)

	return resultLexeme, applicableLexemeFound
}

func (lexicon *GenerationLexicon) ImportFrom(fromLexicon *GenerationLexicon) {
	for _, lexItem := range fromLexicon.allItems {
		lexicon.AddLexItem(lexItem)
	}
}
