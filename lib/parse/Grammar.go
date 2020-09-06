package parse

type Grammar struct {
	readRules       *GrammarRules
	writeRules      *GrammarRules
	tokenizer 		*Tokenizer
}

func NewGrammar() Grammar {
	return Grammar{
		readRules: NewGrammarRules(),
		writeRules: NewGrammarRules(),
		tokenizer: NewTokenizer(DefaultTokenizerExpression),
	}
}

func (grammar *Grammar) SetTokenizer(tokenizer *Tokenizer) {
	grammar.tokenizer = tokenizer
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