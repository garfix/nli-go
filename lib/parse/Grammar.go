package parse

import (
	"nli-go/lib/morphology"
)

type Grammar struct {
	readRules       *GrammarRules
	writeRules      *GrammarRules
	tokenizer 		*Tokenizer
	morphologicalAnalyser *morphology.MorphologicalAnalyser
}

func NewGrammar() Grammar {
	return Grammar{
		readRules: NewGrammarRules(),
		writeRules: NewGrammarRules(),
		tokenizer: NewTokenizer(DefaultTokenizerExpression),
		morphologicalAnalyser: nil,
	}
}

func (grammar *Grammar) SetTokenizer(tokenizer *Tokenizer) {
	grammar.tokenizer = tokenizer
}

func (grammar *Grammar) SetMorphologicalAnalyser(morphologicalAnalyzer *morphology.MorphologicalAnalyser) {
	grammar.morphologicalAnalyser = morphologicalAnalyzer
}

func (grammar *Grammar) GetTokenizer() *Tokenizer {
	return grammar.tokenizer
}

func (grammar *Grammar) GetReadRules() *GrammarRules {
	return grammar.readRules
}

func (grammar *Grammar) GetWriteRules() *GrammarRules {
	return grammar.writeRules
}