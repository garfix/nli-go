package example3

const (
	field_form = "form"
	field_pos = "pos"
	field_sense = "sense"
	field_rule = "rule"
)

type simpleInternalGrammarParser struct {
	tokenizer      *simpleGrammarTokenizer
	lastParsedLine int
}

func NewSimpleInternalGrammarParser() *simpleInternalGrammarParser{
	return &simpleInternalGrammarParser{tokenizer: new(simpleGrammarTokenizer), lastParsedLine: 0}
}

// Parses source into a lexicon
func (parser *simpleInternalGrammarParser) CreateLexicon(source string) (*simpleLexicon, int, bool) {

	lexicon := NewSimpleLexicon()

	// tokenize
	tokens, lineNumber, ok := parser.tokenizer.Tokenize(source)
	if !ok {
		return lexicon, lineNumber, false
	}

	// parse
	parser.lastParsedLine = 0
	lexicon, _, ok = parser.parseLexicon(tokens, 0)

	return lexicon, parser.lastParsedLine, ok
}

// Parses source into transformations
func (parser *simpleInternalGrammarParser) CreateTransformations(source string) ([]SimpleRelationTransformation, int, bool) {

	transformations := []SimpleRelationTransformation{}

	// tokenize
	tokens, lineNumber, tokensOk := parser.tokenizer.Tokenize(source)
	if !tokensOk {
		return transformations, lineNumber, false
	}

	// parse
	parser.lastParsedLine = 0
	transformations, _, ok := parser.parseTransformations(tokens, 0)

	return transformations, parser.lastParsedLine, ok
}

// Parses source into a grammar
func (parser *simpleInternalGrammarParser) CreateGrammar(source string) (*SimpleGrammar, int, bool) {

	grammar := NewSimpleGrammar()

	// tokenize
	tokens, lineNumber, tokensOk := parser.tokenizer.Tokenize(source)
	if !tokensOk {
		return grammar, lineNumber, false
	}

	// parse
	parser.lastParsedLine = 0
	grammar, _, ok := parser.parseGrammar(tokens, 0)

	return grammar, parser.lastParsedLine, ok
}

// Parses source into a relation set
func (parser *simpleInternalGrammarParser) CreateRelationSet(source string) (*SimpleRelationSet, int, bool) {

	relationSet := NewSimpleRelationSet()

	// tokenize
	tokens, lineNumber, tokensOk := parser.tokenizer.Tokenize(source)
	if !tokensOk {
		return relationSet, lineNumber, false
	}

	// parse
	parser.lastParsedLine = 0
	relationSet, _, ok := parser.parseRelationSet(tokens, 0)

	return relationSet, parser.lastParsedLine, ok
}
