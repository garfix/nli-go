package example3

const field_form = "form"
const field_pos = "pos"
const field_sense = "sense"

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

func (parser *simpleInternalGrammarParser) parseTransformations(tokens []SimpleToken, startIndex int) ([]SimpleRelationTransformation, int, bool) {

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

func (parser *simpleInternalGrammarParser) parseTransformation(tokens []SimpleToken, startIndex int) (SimpleRelationTransformation, int, bool) {

	transformation := SimpleRelationTransformation{}
	ok := true

	transformation.Replacement, startIndex, ok = parser.parseRelations(tokens, startIndex)
	if ok {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_implication)
		if ok {
			transformation.Pattern, startIndex, ok = parser.parseRelations(tokens, startIndex)
		}
	}

	return transformation, startIndex, ok
}

func (parser *simpleInternalGrammarParser) parseLexicon(tokens []SimpleToken, startIndex int) (*simpleLexicon, int, bool) {

	lexicon := NewSimpleLexicon()
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)
	for ok {
		lexItem, newStartIndex, lexItemFound := parser.parseLexItem(tokens, startIndex)
		if lexItemFound {
			lexicon.AddLexItem(lexItem)
			startIndex = newStartIndex
		} else {
			ok = false
		}
	}

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_bracket)

	return lexicon, startIndex, ok
}

func (parser *simpleInternalGrammarParser) parseLexItem(tokens []SimpleToken, startIndex int) (SimpleLexItem, int, bool) {

	lexItem := SimpleLexItem{}
	ok, formFound, posFound := true, false, false

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_brace)

	for ok {
		field := ""
		field, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_colon)
			if ok {
				switch field {
				case field_form:
					lexItem.Form, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_stringConstant)
					formFound = true
				case field_pos:
					lexItem.PartOfSpeech, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
					posFound = true;
				case field_sense:
					lexItem.RelationTemplates, startIndex, ok = parser.parseRelations(tokens, startIndex)
				default:
					ok = false
				}
			}
		}
	}

	// required fields
	if formFound && posFound {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_brace)
	} else {
		ok = false
	}

	return lexItem, startIndex, ok
}

func (parser *simpleInternalGrammarParser) parseRelations(tokens []SimpleToken, startIndex int) ([]SimpleRelation, int, bool) {

	relations := []SimpleRelation{}
	ok := true
	commaFound := false

	for ok {

		if len(relations) > 0 {
			_, startIndex, commaFound = parser.parseSingleToken(tokens, startIndex, t_comma)
			if !commaFound {
				break;
			}
		}

		relation := SimpleRelation{}
		relation, startIndex, ok = parser.parseRelation(tokens, startIndex)
		if ok {
			relations = append(relations, relation)
		}
	}

	return relations, startIndex, ok
}

func (parser *simpleInternalGrammarParser) parseRelation(tokens []SimpleToken, startIndex int) (SimpleRelation, int, bool) {

	relation := SimpleRelation{}
	ok := true
	commaFound, argumentFound := false, false
	argument := ""
	newStartIndex := 0

	relation.Predicate, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
	if ok {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_parenthesis)
		for ok {
			if len(relation.Arguments) > 0 {

				// second and further arguments
				_, startIndex, commaFound = parser.parseSingleToken(tokens, startIndex, t_comma)
				if !commaFound {
					break;
				} else {
					argument, startIndex, ok = parser.parseArgument(tokens, startIndex)
					if ok {
						relation.Arguments = append(relation.Arguments, argument)
					}
				}

			} else {

				// first argument (there may not be one, zero arguments are allowed)
				argument, newStartIndex, argumentFound = parser.parseArgument(tokens, startIndex)
				if !argumentFound {
					break;
				} else {
					relation.Arguments = append(relation.Arguments, argument)
					startIndex = newStartIndex
				}

			}
		}
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_parenthesis)
		}
	}

	return relation, startIndex, ok
}

func (parser *simpleInternalGrammarParser) parseArgument(tokens []SimpleToken, startIndex int) (string, int, bool) {

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
func (parser *simpleInternalGrammarParser) parseSingleToken(tokens []SimpleToken, startIndex int, tokenId int) (string, int, bool) {

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