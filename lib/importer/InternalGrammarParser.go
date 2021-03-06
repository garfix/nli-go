package importer

import (
	"fmt"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"nli-go/lib/parse/morphology"
)

const (
	field_sense           = "sense"
	field_condition       = "condition"
	field_result		  = "result"
	field_transformations = "transformations"
	field_rule            = "rule"
	field_preparation     = "preparation"
	field_answer          = "answer"
	field_responses       = "responses"
)

type InternalGrammarParser struct {
	tokenizer        *GrammarTokenizer
	lastParsedResult ParseResult
	panicOnParseFail bool
	// map predicate alias to system-wide module index
	aliasMap map[string]string
}

func NewInternalGrammarParser() *InternalGrammarParser {
	return &InternalGrammarParser{
		tokenizer:        new(GrammarTokenizer),
		lastParsedResult: ParseResult{},
		panicOnParseFail: true,
		aliasMap: map[string]string{"": "", "go": "go"},
	}
}

func (parser *InternalGrammarParser) SetAliasMap(aliasMap map[string]string) {
	parser.aliasMap = aliasMap
}

// automatically panic with meaningful error message on tokenization / parse fail
func (parser *InternalGrammarParser) SetPanicOnParseFail(doPanic bool) {
	parser.panicOnParseFail = doPanic
}

func (parser *InternalGrammarParser) GetLastParseResult() ParseResult {
	return parser.lastParsedResult
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
func (parser *InternalGrammarParser) CreateGrammarRules(source string) *parse.GrammarRules {

	grammar := parse.NewGrammarRules()

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
func (parser *InternalGrammarParser) CreateGenerationGrammar(source string) *parse.GrammarRules {

	grammar := parse.NewGrammarRules()

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
	relation, _, parseOk := parser.parseRelation(tokens, 0, true)
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
	relationSet, _, parseOk := parser.parseRelations(tokens, 0, true)
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

func (parser *InternalGrammarParser) CreateSortRelations(source string) []mentalese.SortRelation {

	sortRelations := []mentalese.SortRelation{}

	// tokenize
	parser.lastParsedResult.LineNumber = 0
	tokens, lineNumber, tokensOk := parser.tokenizer.Tokenize(source)
	parser.processResult(service_tokenizer, tokensOk, source, lineNumber)
	if !tokensOk {
		return sortRelations
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	sortRelations, _, parseOk := parser.parseSortRelations(tokens, 0)
	parser.processResult(service_parser, parseOk, source, parser.lastParsedResult.LineNumber)

	return sortRelations
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
		return mentalese.NewBinding()
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	binding, _, tokensOk := parser.parseBinding(tokens, 0)
	parser.processResult(service_parser, tokensOk, source, parser.lastParsedResult.LineNumber)

	return binding
}

func (parser *InternalGrammarParser) CreateBindings(source string) mentalese.BindingSet {

	// tokenize
	parser.lastParsedResult.LineNumber = 0
	tokens, _, tokensOk := parser.tokenizer.Tokenize(source)
	parser.processResult(service_tokenizer, tokensOk, source, parser.lastParsedResult.LineNumber)
	if !tokensOk {
		return mentalese.NewBindingSet()
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	bindings, _, tokensOk := parser.parseBindings(tokens, 0)
	parser.processResult(service_parser, tokensOk, source, parser.lastParsedResult.LineNumber)

	return bindings
}

func (parser *InternalGrammarParser) CreateSegmentationRules(source string) *morphology.SegmentationRules {
	segmentationRules := morphology.NewSegmentationRules()

	// tokenize
	parser.lastParsedResult.LineNumber = 0
	tokens, _, tokensOk := parser.tokenizer.Tokenize(source)
	parser.processResult(service_tokenizer, tokensOk, source, parser.lastParsedResult.LineNumber)
	if !tokensOk {
		return segmentationRules
	}

	// parse
	parser.lastParsedResult.LineNumber = 0
	segmentationRules, _, tokensOk = parser.parseSegmentationRulesAndCharacterClasses(tokens, 0)
	parser.processResult(service_parser, tokensOk, source, parser.lastParsedResult.LineNumber)

	return segmentationRules
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
