package example3

import (
	"nli-go/lib/example2"
)

// parses a string of transformations like
//
// list-customers(P1) :- predicate(P1, name), object(P1, E1), instance_of(E1, customer)


type simpleRelationTransformationParser struct {
	lastParsedLine int
}

func NewSimpleRelationTransformationParser() *simpleRelationTransformationParser{
	return &simpleRelationTransformationParser{}
}

// Parses source and returns its transformations
func (parser *simpleRelationTransformationParser) ParseString(source string) ([]SimpleRelationTransformation, int, bool) {

	transformations := []SimpleRelationTransformation{}
	tok := NewSimpleGrammarTokenizer()

	// tokenize
	tokens, lineNumber, tokensOk := tok.Tokenize(source)
	if !tokensOk {
		return transformations, lineNumber, false
	}

	// parse
	transformations, _, ok := parser.parseTransformations(tokens, 0)

	return transformations, parser.lastParsedLine, ok
}

func (parser *simpleRelationTransformationParser) parseTransformations(tokens []SimpleToken, startIndex int) ([]SimpleRelationTransformation, int, bool) {

	transformations := []SimpleRelationTransformation{}
	ok := true

	for startIndex < len(tokens) {
		transformation := SimpleRelationTransformation{}
		transformation, startIndex, ok = parser.parseTransformation(tokens, startIndex)
		if ok {
			transformations = append(transformations, transformation)
		} else {
			break;
		}
	}

	return transformations, startIndex, ok
}

func (parser *simpleRelationTransformationParser) parseTransformation(tokens []SimpleToken, startIndex int) (SimpleRelationTransformation, int, bool) {

	transformation := SimpleRelationTransformation{}
	ok := true

	transformation.Pattern, startIndex, ok = parser.parseRelations(tokens, startIndex)
	if ok {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_implication)
		if ok {
			transformation.Replacement, startIndex, ok = parser.parseRelations(tokens, startIndex)
		}
	}

	return transformation, startIndex, ok
}

func (parser *simpleRelationTransformationParser) parseRelations(tokens []SimpleToken, startIndex int) ([]example2.SimpleRelation, int, bool) {

	relations := []example2.SimpleRelation{}
	ok := true
	commaFound := false

	for ok {
		relation := example2.SimpleRelation{}
		relation, startIndex, ok = parser.parseRelation(tokens, startIndex)
		if ok {
			relations = append(relations, relation)
		} else {
			break;
		}

		_, startIndex, commaFound = parser.parseSingleToken(tokens, startIndex, t_comma)
		if !commaFound {
			break;
		}
	}

	return relations, startIndex, ok
}

func (parser *simpleRelationTransformationParser) parseRelation(tokens []SimpleToken, startIndex int) (example2.SimpleRelation, int, bool) {

	relation := example2.SimpleRelation{}
	ok := true
	commaFound := false

	relation.Predicate, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
	if ok {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_parenthesis)
		for ok {
			if len(relation.Arguments) > 0 {
				_, startIndex, commaFound = parser.parseSingleToken(tokens, startIndex, t_comma)
				if (!commaFound) {
					break;
				}
			}
			argument := ""
			argument, startIndex, ok = parser.parseArgument(tokens, startIndex)
			if ok {
				relation.Arguments = append(relation.Arguments, argument)
			}
		}
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_parenthesis)
		}
	}

	return relation, startIndex, ok
}

func (parser *simpleRelationTransformationParser) parseArgument(tokens []SimpleToken, startIndex int) (string, int, bool) {

	ok := false
	tokenValue := ""

	tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
	if !ok {
		tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_variable)
		if !ok {
			tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_number)
			if !ok {
				tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_stringConstant)
			}
		}
	}

	return tokenValue, startIndex, ok
}

// (!) startIndex increases only if the specified token could be matched
func (parser *simpleRelationTransformationParser) parseSingleToken(tokens []SimpleToken, startIndex int, tokenId int) (string, int, bool) {

	ok := false
	tokenValue := ""

	if startIndex < len(tokens) {
		token := tokens[startIndex]
		ok = (token.TokenId == tokenId)
		if ok {
			tokenValue = token.TokenValue
			if tokens[startIndex].LineNumber > parser.lastParsedLine {
				parser.lastParsedLine = tokens[startIndex].LineNumber
			}
			startIndex++
		}
	}

	return tokenValue, startIndex, ok
}