package generate

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
)

type GenerationLexicon struct {
	lexemes map[string][]GenerationLexeme
	matcher mentalese.RelationMatcher
}

func NewGenerationLexicon() *GenerationLexicon {
	return &GenerationLexicon{lexemes: map[string][]GenerationLexeme{}}
}

func (lexicon *GenerationLexicon) AddLexItem(lexItem GenerationLexeme) {
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

	common.LogTree("GetLexemeForGeneration", consequent)

	if consequent.Predicate == "number" {
		resultLexeme.Form = consequent.Arguments[0].TermValue
		return resultLexeme, true
	}


	partOfSpeech := consequent.Predicate

	lexemes, found := lexicon.lexemes[partOfSpeech]
	if found {

		binding := mentalese.Binding{"E": consequent.Arguments[0]}

		for _, lexeme := range lexemes {

			bindings, _, match := lexicon.matcher.MatchSequenceToSet(lexeme.Condition, sentenseSense, binding)

			if match {
				resultLexeme = lexeme

				if partOfSpeech == "proper_noun" {
					resultLexeme.Form = bindings[0]["Name"].TermValue
				}

				break
			}
		}
	}

	common.LogTree("GetLexemeForGeneration", resultLexeme, found)

	return resultLexeme, found
}
