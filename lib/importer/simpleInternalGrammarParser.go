package importer

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
)

const (
	field_form = "form"
	field_pos = "pos"
	field_sense = "sense"
	field_rule = "rule"
	field_question = "Q"
	field_answer = "A"
)

type simpleInternalGrammarParser struct {
	tokenizer      *simpleGrammarTokenizer
	lastParsedLine int
}

func NewSimpleInternalGrammarParser() *simpleInternalGrammarParser{
	return &simpleInternalGrammarParser{tokenizer: new(simpleGrammarTokenizer), lastParsedLine: 0}
}

// Parses source into a lexicon
func (parser *simpleInternalGrammarParser) CreateLexicon(source string) (*parse.SimpleLexicon, int, bool) {

	lexicon := parse.NewSimpleLexicon()

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
func (parser *simpleInternalGrammarParser) CreateTransformations(source string) ([]mentalese.SimpleRelationTransformation, int, bool) {

	transformations := []mentalese.SimpleRelationTransformation{}

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

// Parses source into rules
func (parser *simpleInternalGrammarParser) CreateRules(source string) ([]mentalese.SimpleRule, int, bool) {

	rules := []mentalese.SimpleRule{}

	// tokenize
	tokens, lineNumber, tokensOk := parser.tokenizer.Tokenize(source)
	if !tokensOk {
		return rules, lineNumber, false
	}

	// parse
	parser.lastParsedLine = 0
	rules, _, ok := parser.parseRules(tokens, 0)

	return rules, parser.lastParsedLine, ok
}

// Parses source into a grammar
func (parser *simpleInternalGrammarParser) CreateGrammar(source string) (*parse.SimpleGrammar, int, bool) {

	grammar := parse.NewSimpleGrammar()

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
func (parser *simpleInternalGrammarParser) CreateRelationSet(source string) (*mentalese.SimpleRelationSet, int, bool) {

	relationSet := mentalese.NewSimpleRelationSet()

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

func (parser *simpleInternalGrammarParser) CreateQAPairs(source string) ([]mentalese.SimpleQAPair, int, bool) {

	qaPairs := []mentalese.SimpleQAPair{}

	// tokenize
	tokens, lineNumber, tokensOk := parser.tokenizer.Tokenize(source)
	if !tokensOk {
		return qaPairs, lineNumber, false
	}

	// parse
	parser.lastParsedLine = 0
	qaPairs, _, ok := parser.parseQAPairs(tokens, 0)

	return qaPairs, parser.lastParsedLine, ok
}
