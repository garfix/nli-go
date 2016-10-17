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

type InternalGrammarParser struct {
	tokenizer      *GrammarTokenizer
	lastParsedLine int
}

func NewInternalGrammarParser() *InternalGrammarParser {
	return &InternalGrammarParser{tokenizer: new(GrammarTokenizer), lastParsedLine: 0}
}

// Parses source into a lexicon
func (parser *InternalGrammarParser) CreateLexicon(source string) (*parse.Lexicon, int, bool) {

	lexicon := parse.NewLexicon()

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
func (parser *InternalGrammarParser) CreateTransformations(source string) ([]mentalese.RelationTransformation, int, bool) {

	transformations := []mentalese.RelationTransformation{}

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
func (parser *InternalGrammarParser) CreateRules(source string) ([]mentalese.Rule, int, bool) {

	rules := []mentalese.Rule{}

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
func (parser *InternalGrammarParser) CreateGrammar(source string) (*parse.Grammar, int, bool) {

	grammar := parse.NewGrammar()

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
func (parser *InternalGrammarParser) CreateRelation(source string) (mentalese.Relation, bool) {

	relation := mentalese.Relation{}

	// tokenize
	tokens, _, tokensOk := parser.tokenizer.Tokenize(source)
	if !tokensOk {
		return relation, false
	}

	// parse
	parser.lastParsedLine = 0
	relation, _, ok := parser.parseRelation(tokens, 0)

	return relation, ok
}

// Parses source into a relation set
func (parser *InternalGrammarParser) CreateRelationSet(source string) (mentalese.RelationSet, int, bool) {

	relationSet := mentalese.RelationSet{}

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

func (parser *InternalGrammarParser) CreateQAPairs(source string) ([]mentalese.QAPair, int, bool) {

	qaPairs := []mentalese.QAPair{}

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

func (parser *InternalGrammarParser) CreateTerm(source string) (mentalese.Term, bool) {

	// tokenize
	tokens, _, tokensOk := parser.tokenizer.Tokenize(source)
	if !tokensOk {
		return mentalese.Term{}, false
	}

	// parse
	term, _, ok := parser.parseTerm(tokens, 0)

	return term, ok
}

func (parser *InternalGrammarParser) CreateBinding(source string) (mentalese.Binding, bool) {

	// tokenize
	tokens, _, tokensOk := parser.tokenizer.Tokenize(source)
	if !tokensOk {
		return mentalese.Binding{}, false
	}

	// parse
	binding, _, ok := parser.parseBinding(tokens, 0)

	return binding, ok
}
