package parse

type Grammar struct {
	readRules             *GrammarRules
	writeRules            *GrammarRules
	tokenizer             *Tokenizer
	morphologicalAnalyzer *MorphologicalAnalyzer
}

func NewGrammar() Grammar {
	return Grammar{
		readRules:             NewGrammarRules(),
		writeRules:            NewGrammarRules(),
		tokenizer:             NewTokenizer(DefaultTokenizerExpression),
		morphologicalAnalyzer: nil,
	}
}

func (grammar *Grammar) SetTokenizer(tokenizer *Tokenizer) {
	grammar.tokenizer = tokenizer
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

func (grammar *Grammar) GetReadRules() *GrammarRules {
	return grammar.readRules
}

func (grammar *Grammar) GetWriteRules() *GrammarRules {
	return grammar.writeRules
}