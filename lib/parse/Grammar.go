package parse

type Grammar struct {
	locale 			      string
	readRules             *GrammarRules
	writeRules            *GrammarRules
	tokenizer             *Tokenizer
	morphologicalAnalyzer *MorphologicalAnalyzer
	texts				  map[string]string
}

func NewGrammar(locale string) Grammar {
	return Grammar{
		locale: 			   locale,
		readRules:             NewGrammarRules(),
		writeRules:            NewGrammarRules(),
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

func (grammar *Grammar) GetReadRules() *GrammarRules {
	return grammar.readRules
}

func (grammar *Grammar) GetWriteRules() *GrammarRules {
	return grammar.writeRules
}