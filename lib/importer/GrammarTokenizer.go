package importer

// A tokenizer for expressions like
//
// list_customers(P1) :- predicate(P1, name), object(P1, E1), instance_of(E1, customer)

import (
	"regexp"
	"strings"
)

const (
	_ = iota // number these constants 1, 2, ...
	t_predicate
	t_variable
	t_anonymousVariable
	t_stringConstant
	t_number
	t_comma
	t_implication
	t_colon
	t_opening_parenthesis
	t_closing_parenthesis
	t_opening_bracket
	t_closing_bracket
	t_opening_brace
	t_closing_brace
	_newline
	_other
)

type GrammarTokenizer struct {

}

func NewGrammarTokenizer() *GrammarTokenizer {
	return &GrammarTokenizer{}
}

func (tok *GrammarTokenizer) Tokenize(source string) ([]Token, int, bool) {

	// 2x \ = 1x escape for Go; 4x \ = 1x escape for Regexp engine

	tokenExpressions := []struct{
		id int
		pattern string} {
		{t_predicate, "[a-z][a-z_]*"},
		{t_variable, "[A-Z][A-Za-z0-9_]*"},
		{t_anonymousVariable, "_"},
		{t_stringConstant, "'(?:\\\\'|\\\\\\\\|[^'])*'"},
		{t_number, "[0-9]+"},
		{t_comma, ","},
		{t_implication, ":-"},
		{t_colon, ":"},
		{t_opening_parenthesis, "\\("},
		{t_closing_parenthesis, "\\)"},
		{t_opening_bracket, "\\["},
		{t_closing_bracket, "\\]"},
		{t_opening_brace, "\\{"},
		{t_closing_brace, "\\}"},
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

	tokens := []Token{}
	lineNumber := 1
	for _, tokenValues := range expression.FindAllStringSubmatch(source, -1) {

		tokenId := 0
		tokenValue := ""

		for i := 1; i < len(tokenValues); i++  {
			// http://stackoverflow.com/a/18595217
			if len(tokenValues[i]) != 0 {
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
			tokenValue = tokenValue[1:len(tokenValue)-1]
		} else if tokenId == _other || tokenId == 0 {
			success = false
			break;
		}

		tokens = append(tokens, Token{tokenId, lineNumber, tokenValue})
	}

	return tokens, lineNumber, success
}