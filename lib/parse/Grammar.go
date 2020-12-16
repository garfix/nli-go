package parse

type Grammar struct {
	readRules       *GrammarRules
	writeRules      *GrammarRules
	tokenizer 		*Tokenizer
	morphologicalAnalyser *MorphologicalAnalyser
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

func (grammar *Grammar) GetTokenizer() *Tokenizer {
	return grammar.tokenizer
}

func (grammar *Grammar) SetMorphologicalAnalyser(morphologicalAnalyzer *MorphologicalAnalyser) {
	grammar.morphologicalAnalyser = morphologicalAnalyzer
}

func (grammar *Grammar) GetMorphologicalAnalyser() *MorphologicalAnalyser {
	return grammar.morphologicalAnalyser
}

func (grammar *Grammar) GetReadRules() *GrammarRules {
	return grammar.readRules
}

func (grammar *Grammar) GetWriteRules() *GrammarRules {
	return grammar.writeRules
}