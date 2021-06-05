package parse

import "nli-go/lib/mentalese"

type Grammar struct {
	locale 			      string
	readRules             *mentalese.GrammarRules
	writeRules            *mentalese.GrammarRules
	tokenizer             *Tokenizer
	morphologicalAnalyzer *MorphologicalAnalyzer
	texts				  map[string]string
}

func NewGrammar(locale string) Grammar {
	return Grammar{
		locale: 			   locale,
		readRules:             mentalese.NewGrammarRules(),
		writeRules:            mentalese.NewGrammarRules(),
		tokenizer:             NewTokenizer(DefaultTokenizerExpression),
		morphologicalAnalyzer: nil,
		texts:				   map[string]string{},
	}
}

func (grammar *Grammar) SetTokenizer(tokenizer *Tokenizer) {
	grammar.tokenizer = tokenizer
}

func (grammar *Grammar) SetTexts(texts map[string]string) {
	grammar.texts = texts
}

func (grammar *Grammar) GetText(text string) string {
	translation, found := grammar.texts[text]
	if found {
		return translation
	} else {
		return text
	}
}

func (grammar *Grammar) GetTokenizer() *Tokenizer {
	return grammar.tokenizer
}

func (grammar *Grammar) SetMorphologicalAnalyzer(morphologicalAnalyzer *MorphologicalAnalyzer) {
	grammar.morphologicalAnalyzer = morphologicalAnalyzer
}

func (grammar *Grammar) GetMorphologicalAnalyzer() *MorphologicalAnalyzer {
	return grammar.morphologicalAnalyzer
}

func (grammar *Grammar) GetLocale() string {
	return grammar.locale
}

func (grammar *Grammar) GetReadRules() *mentalese.GrammarRules {
	return grammar.readRules
}

func (grammar *Grammar) GetWriteRules() *mentalese.GrammarRules {
	return grammar.writeRules
}