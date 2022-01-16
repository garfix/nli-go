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
	tComment
	tPredicate
	tVariable
	tEquals
	tNotEquals
	tAssign
	tPlaceholder
	tAnonymousVariable
	tId
	tStringConstant
	tRegExp
	tNumber
	tComma
	tRewrite
	tImplication
	tColon
	tAmpersand
	tSemicolon
	tGtEq
	tGt
	tLtEq
	tLt
	tOpeningParenthesis
	tClosingParenthesis
	tOpeningBracket
	tClosingBracket
	tDoubleOpeningBrace
	tDoubleClosingBrace
	tOpeningBrace
	tClosingBrace
	tNegative
	tPositive
	tSlash
	tMultiply
	tUp
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

	tokenExpressions := []struct {
		id      int
		pattern string
	}{
		{tComment, "/\\*.*?\\*/"},
		{tPredicate, "[a-z][a-z0-9_]*"},
		{tVariable, "[A-Z][a-zA-Z0-9]*"},
		{tEquals, "=="},
		{tNotEquals, "!="},
		{tAssign, ":="},
		{tPlaceholder, "\\$"},
		{tAnonymousVariable, "_"},
		{tId, "`[^`]+`"},
		{tStringConstant, "'(?:\\\\'|\\\\\\\\|[^'])*'"},
		{tRegExp, "~(?:\\\\~|\\\\\\\\|[^~])*~"},
		{tNumber, "-?[0-9]+"},
		{tComma, ","},
		{tRewrite, "->"},
		{tImplication, ":-"},
		{tColon, ":"},
		{tAmpersand, "&"},
		{tSemicolon, ";"},
		{tGtEq, ">="},
		{tGt, ">"},
		{tLtEq, "<="},
		{tLt, "<"},
		{tOpeningParenthesis, "\\("},
		{tClosingParenthesis, "\\)"},
		{tOpeningBracket, "\\["},
		{tClosingBracket, "\\]"},
		{tDoubleOpeningBrace, "\\{\\{"},
		{tDoubleClosingBrace, "\\}\\}"},
		{tOpeningBrace, "\\{"},
		{tClosingBrace, "\\}"},
		{tNegative, "-"},
		{tPositive, "\\+"},
		{tSlash, "/"},
		{tMultiply, "\\*"},
		{tUp, "\\.\\."},
		{_newline, "(?:\r\n|\n|\r)"},
		{_other, "."},
	}

	success := true

	pattern, sep := "", ""
	for _, tokenExpression := range tokenExpressions {
		pattern += sep + "[ \t]*(" + tokenExpression.pattern + ")"
		sep = "|"
	}

	expression, _ := regexp.Compile("(?s)(?:" + pattern + ")[ \t]*")

	tokens := []Token{}
	lineNumber := 1
	for _, tokenValues := range expression.FindAllStringSubmatch(source, -1) {

		tokenId := 0
		tokenValue := ""

		for i := 1; i < len(tokenValues); i++ {
			// http://stackoverflow.com/a/18595217
			if len(tokenValues[i]) != 0 {
				tokenId = i
				tokenValue = tokenValues[i]
			}
		}

		if tokenId == _newline {
			lineNumber++
			continue
		} else if tokenId == tComment {
			lineNumber += strings.Count(tokenValue, "\n")
			continue
		} else if tokenId == tStringConstant {
			tokenValue = strings.Replace(tokenValue, "\\'", "'", -1)
			tokenValue = strings.Replace(tokenValue, "\\\\", "\\", -1)
			tokenValue = tokenValue[1 : len(tokenValue)-1]
		} else if tokenId == tId {
			tokenValue = strings.Replace(tokenValue, "\\`", "`", -1)
			tokenValue = strings.Replace(tokenValue, "\\\\", "\\", -1)
			tokenValue = tokenValue[1 : len(tokenValue)-1]
		} else if tokenId == tRegExp {
			tokenValue = strings.Replace(tokenValue, "\\/", "/", -1)
			tokenValue = strings.Replace(tokenValue, "\\\\", "\\", -1)
			tokenValue = tokenValue[1 : len(tokenValue)-1]
		} else if tokenId == _other || tokenId == 0 {
			success = false
			break
		}

		tokens = append(tokens, Token{tokenId, lineNumber, tokenValue})
	}

	return tokens, lineNumber, success
}
