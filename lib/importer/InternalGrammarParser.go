package importer

import (
	"fmt"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
)

const (
	field_form            = "form"
	field_pos             = "pos"
	field_sense           = "sense"
	field_condition       = "condition"
	field_transformations = "transformations"
	field_rule            = "rule"
	field_preparation     = "preparation"
	field_answer          = "answer"
	field_no_results      = "no_results"
	field_some_results    = "some_results"
)

type InternalGrammarParser struct {
	tokenizer        *GrammarTokenizer
	lastParsedResult ParseResult
	panicOnParseFail bool
}

func NewInternalGrammarParser() *InternalGrammarParser {
	return &InternalGrammarParser{
		tokenizer:        new(GrammarTokenizer),
		lastParsedResult: ParseResult{},
		panicOnParseFail: true,
	}
}

// automatically panic with meaningful error message on tokenization / parse fail
func (parser *InternalGrammarParser) SetPanicOnParseFail(doPanic bool) {
	parser.panicOnParseFail = doPanic
}

func (parser *InternalGrammarParser) GetLastParseResult() ParseResult {
	return parser.lastParsedResult
}

// Parses source into a lexicon
func (parser *InternalGrammarParser) CreateLexicon(source string) *parse.Lexicon {

	lexicon := parse.NewLexicon()

	// tokenize
	parser.lastParsedResult.LineNumber = 0
	tokens, lineNumber, tokensOk := parser.tokenizer.Tokenize(source)
	parser.processResult(service_tokenizer, tokensOk, source, lineNumber)
	if !tokensOk {
		return lexicon
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	lexicon, _, parseOk := parser.parseLexicon(tokens, 0)
	parser.processResult(service_parser, parseOk, source, parser.lastParsedResult.LineNumber)

	return lexicon
}

// Parses source into a lexicon
func (parser *InternalGrammarParser) CreateGenerationLexicon(source string, log *common.SystemLog) *generate.GenerationLexicon {

	lexicon := generate.NewGenerationLexicon(log)

	// tokenize
	parser.lastParsedResult.LineNumber = 0
	tokens, lineNumber, tokensOk := parser.tokenizer.Tokenize(source)
	parser.processResult(service_tokenizer, tokensOk, source, lineNumber)
	if !tokensOk {
		return lexicon
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	lexicon, _, parseOk := parser.parseGenerationLexicon(tokens, 0, log)
	parser.processResult(service_parser, parseOk, source, parser.lastParsedResult.LineNumber)

	return lexicon
}

// Parses source into transformations
func (parser *InternalGrammarParser) CreateTransformations(source string) []mentalese.RelationTransformation {

	transformations := []mentalese.RelationTransformation{}

	// tokenize
	parser.lastParsedResult.LineNumber = 0
	tokens, lineNumber, tokensOk := parser.tokenizer.Tokenize(source)
	parser.processResult(service_tokenizer, tokensOk, source, lineNumber)
	if !tokensOk {
		return transformations
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	transformations, _, parseOk := parser.parseTransformations(tokens, 0)
	parser.processResult(service_parser, parseOk, source, parser.lastParsedResult.LineNumber)

	return transformations
}

// Parses source into rules
func (parser *InternalGrammarParser) CreateRules(source string) []mentalese.Rule {

	rules := []mentalese.Rule{}

	// tokenize
	parser.lastParsedResult.LineNumber = 0
	tokens, lineNumber, tokensOk := parser.tokenizer.Tokenize(source)
	parser.processResult(service_tokenizer, tokensOk, source, lineNumber)
	if !tokensOk {
		return rules
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	rules, _, parseOk := parser.parseRules(tokens, 0)
	parser.processResult(service_parser, parseOk, source, parser.lastParsedResult.LineNumber)

	return rules
}

// Parses source into a grammar
func (parser *InternalGrammarParser) CreateGrammar(source string) *parse.Grammar {

	grammar := parse.NewGrammar()

	// tokenize
	parser.lastParsedResult.LineNumber = 0
	tokens, lineNumber, tokensOk := parser.tokenizer.Tokenize(source)
	parser.processResult(service_tokenizer, tokensOk, source, lineNumber)
	if !tokensOk {
		return grammar
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	grammar, _, parseOk := parser.parseGrammar(tokens, 0)
	parser.processResult(service_parser, parseOk, source, parser.lastParsedResult.LineNumber)

	return grammar
}

// Parses source into a grammar
func (parser *InternalGrammarParser) CreateGenerationGrammar(source string) *generate.GenerationGrammar {

	grammar := generate.NewGenerationGrammar()

	// tokenize
	parser.lastParsedResult.LineNumber = 0
	tokens, lineNumber, tokensOk := parser.tokenizer.Tokenize(source)
	parser.processResult(service_tokenizer, tokensOk, source, lineNumber)
	if !tokensOk {
		return grammar
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	grammar, _, parseOk := parser.parseGenerationGrammar(tokens, 0)
	parser.processResult(service_parser, parseOk, source, parser.lastParsedResult.LineNumber)

	return grammar
}

func (parser *InternalGrammarParser) LoadText(path string) string {

	source, err := common.ReadFile(path)

	if err != nil {
		parser.processResult(file_read, false, fmt.Sprint(err), 0)
	}
	return source
}

// Parses source into a relation set
func (parser *InternalGrammarParser) CreateRelation(source string) mentalese.Relation {

	relation := mentalese.Relation{}

	// tokenize
	parser.lastParsedResult.LineNumber = 0
	tokens, lineNumber, tokensOk := parser.tokenizer.Tokenize(source)
	parser.processResult(service_tokenizer, tokensOk, source, lineNumber)
	if !tokensOk {
		return relation
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	relation, _, parseOk := parser.parseRelation(tokens, 0)
	parser.processResult(service_parser, parseOk, source, parser.lastParsedResult.LineNumber)

	return relation
}

// Parses source into a relation set
func (parser *InternalGrammarParser) CreateRelationSet(source string) mentalese.RelationSet {

	relationSet := mentalese.RelationSet{}

	// tokenize
	parser.lastParsedResult.LineNumber = 0
	tokens, lineNumber, tokensOk := parser.tokenizer.Tokenize(source)
	parser.processResult(service_tokenizer, tokensOk, source, lineNumber)
	if !tokensOk {
		return relationSet
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	relationSet, _, parseOk := parser.parseRelationSet(tokens, 0)
	parser.processResult(service_parser, parseOk, source, parser.lastParsedResult.LineNumber)

	return relationSet
}

func (parser *InternalGrammarParser) CreateSolutions(source string) []mentalese.Solution {

	solutions := []mentalese.Solution{}

	// tokenize
	parser.lastParsedResult.LineNumber = 0
	tokens, lineNumber, tokensOk := parser.tokenizer.Tokenize(source)
	parser.processResult(service_tokenizer, tokensOk, source, lineNumber)
	if !tokensOk {
		return solutions
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	solutions, _, parseOk := parser.parseSolutions(tokens, 0)
	parser.processResult(service_parser, parseOk, source, parser.lastParsedResult.LineNumber)

	return solutions
}

func (parser *InternalGrammarParser) CreateTerm(source string) mentalese.Term {

	// tokenize
	parser.lastParsedResult.LineNumber = 0
	tokens, _, tokensOk := parser.tokenizer.Tokenize(source)
	parser.processResult(service_tokenizer, tokensOk, source, parser.lastParsedResult.LineNumber)
	if !tokensOk {
		return mentalese.Term{}
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	term, _, tokensOk := parser.parseTerm(tokens, 0)
	parser.processResult(service_parser, tokensOk, source, parser.lastParsedResult.LineNumber)

	return term
}

func (parser *InternalGrammarParser) CreateBinding(source string) mentalese.Binding {

	// tokenize
	parser.lastParsedResult.LineNumber = 0
	tokens, _, tokensOk := parser.tokenizer.Tokenize(source)
	parser.processResult(service_tokenizer, tokensOk, source, parser.lastParsedResult.LineNumber)
	if !tokensOk {
		return mentalese.Binding{}
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	binding, _, tokensOk := parser.parseBinding(tokens, 0)
	parser.processResult(service_parser, tokensOk, source, parser.lastParsedResult.LineNumber)

	return binding
}

func (parser *InternalGrammarParser) CreateBindings(source string) []mentalese.Binding {

	// tokenize
	parser.lastParsedResult.LineNumber = 0
	tokens, _, tokensOk := parser.tokenizer.Tokenize(source)
	parser.processResult(service_tokenizer, tokensOk, source, parser.lastParsedResult.LineNumber)
	if !tokensOk {
		return []mentalese.Binding{}
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	bindings, _, tokensOk := parser.parseBindings(tokens, 0)
	parser.processResult(service_parser, tokensOk, source, parser.lastParsedResult.LineNumber)

	return bindings
}

func (parser *InternalGrammarParser) processResult(service string, ok bool, source string, lineNumber int) {

	parser.lastParsedResult.Service = service
	parser.lastParsedResult.Ok = ok
	parser.lastParsedResult.LineNumber = lineNumber
	parser.lastParsedResult.Source = source

	if !ok && parser.panicOnParseFail {
		panic(parser.lastParsedResult.String())
	}
}
