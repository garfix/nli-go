package example3

// A tokenizer for expressions like
//
// list_customers(P1) :- predicate(P1, name), object(P1, E1), instance_of(E1, customer)

import (
	"regexp"
	"strings"
)

const t_predicate = 1
const t_variable = 2
const t_stringConstant = 3
const t_number = 4
const t_comma = 5
const t_implication = 6
const t_opening_parenthesis = 7
const t_closing_parenthesis = 8
const _newline = 9
const _other = 10

type simpleGrammarTokenizer struct {

}

func NewSimpleGrammarTokenizer() *simpleGrammarTokenizer {
	return &simpleGrammarTokenizer{}
}

func (tok *simpleGrammarTokenizer) Tokenize(source string) ([]SimpleToken, int, bool) {

	// 2x \ = 1x escape for Go; 4x \ = 1x escape for Regexp engine

	tokenExpressions := []struct{
		id int
		pattern string} {
		{t_predicate, "[a-z][a-z_]*"},
		{t_variable, "[A-Z][A-Za-z0-9_]*"},
		{t_stringConstant, "'(?:\\\\'|\\\\\\\\|[^'])*'"},
		{t_number, "[0-9]+"},
		{t_comma, ","},
		{t_implication, ":="},
		{t_opening_parenthesis, "\\("},
		{t_closing_parenthesis, "\\)"},
		{_newline, "(?:\r\n|\n|\r)"},
		{_other, "."},
	}

	success := true

	pattern, sep := "", ""
	for _, tokenExpression := range tokenExpressions {
		pattern += sep + "[ \t]*(" + tokenExpression.pattern + ")"
		sep = "|"
	}

	expression, _ := regexp.Compile("(?:" + pattern + ")[ \t]*")

	tokens := []SimpleToken{}
	lineNumber := 1
	for _, tokenValues := range expression.FindAllStringSubmatch(source, -1) {

		tokenId := 0
		tokenValue := ""

		for i := 1; i < len(tokenValues); i++  {
			if tokenValues[i] != "" {
				tokenId = i
				tokenValue = tokenValues[i]
			}
		}

		if tokenId == _newline {
			lineNumber++
			continue;
		} else if tokenId == t_stringConstant {
			tokenValue = strings.Replace(tokenValue, "\\'", "'", -1)
			tokenValue = strings.Replace(tokenValue, "\\\\", "\\", -1)
		} else if tokenId == _other || tokenId == 0 {
			success = false
			break;
		}

		tokens = append(tokens, SimpleToken{tokenId, lineNumber, tokenValue})
	}

	return tokens, lineNumber, success
}