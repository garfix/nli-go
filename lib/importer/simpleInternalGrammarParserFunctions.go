package importer

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
)

func (parser *simpleInternalGrammarParser) parseRelationSet(tokens []SimpleToken, startIndex int) (mentalese.SimpleRelationSet, int, bool) {

	relationSet := mentalese.SimpleRelationSet{}
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)

	for startIndex < len(tokens) {
		relation := mentalese.SimpleRelation{}
		relation, startIndex, ok = parser.parseRelation(tokens, startIndex)
		if ok {
			relationSet = append(relationSet, relation)
		} else {
			break;
		}
	}

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_bracket)

	return relationSet, startIndex, ok
}

func (parser *simpleInternalGrammarParser) parseTransformations(tokens []SimpleToken, startIndex int) ([]mentalese.SimpleRelationTransformation, int, bool) {

	transformations := []mentalese.SimpleRelationTransformation{}
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)

	for startIndex < len(tokens) {
		transformation := mentalese.SimpleRelationTransformation{}
		transformation, startIndex, ok = parser.parseTransformation(tokens, startIndex)
		if ok {
			transformations = append(transformations, transformation)
		} else {
			break;
		}
	}

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_bracket)

	return transformations, startIndex, ok
}

func (parser *simpleInternalGrammarParser) parseTransformation(tokens []SimpleToken, startIndex int) (mentalese.SimpleRelationTransformation, int, bool) {

	transformation := mentalese.SimpleRelationTransformation{}
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

func (parser *simpleInternalGrammarParser) parseRules(tokens []SimpleToken, startIndex int) ([]mentalese.SimpleRule, int, bool) {

	rules := []mentalese.SimpleRule{}
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)

	for startIndex < len(tokens) {
		rule := mentalese.SimpleRule{}
		rule, startIndex, ok = parser.parseRule(tokens, startIndex)
		if ok {
			rules = append(rules, rule)
		} else {
			break;
		}
	}

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_bracket)

	return rules, startIndex, ok
}

func (parser *simpleInternalGrammarParser) parseRule(tokens []SimpleToken, startIndex int) (mentalese.SimpleRule, int, bool) {

	rule := mentalese.SimpleRule{}
	ok := true

	rule.Goal, startIndex, ok = parser.parseRelation(tokens, startIndex)
	if ok {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_implication)
		if ok {
			rule.Pattern, startIndex, ok = parser.parseRelations(tokens, startIndex)
		}
	}

	return rule, startIndex, ok
}

func (parser *simpleInternalGrammarParser) parseQAPairs(tokens []SimpleToken, startIndex int) ([]mentalese.SimpleQAPair, int, bool) {

	qaPairs := []mentalese.SimpleQAPair{}
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)

	for startIndex < len(tokens) {
		qaPair := mentalese.SimpleQAPair{}
		qaPair, startIndex, ok = parser.parseQAPair(tokens, startIndex)
		if ok {
			qaPairs = append(qaPairs, qaPair)
		} else {
			break;
		}
	}

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_bracket)

	return qaPairs, startIndex, ok
}

func (parser *simpleInternalGrammarParser) parseQAPair(tokens []SimpleToken, startIndex int) (mentalese.SimpleQAPair, int, bool) {

	qaPair := mentalese.SimpleQAPair{}
	ok, formFound, posFound := true, false, false

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_brace)

	for ok {
		field := ""
		field, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_colon)
			if ok {
				switch field {
				case field_question:
					qaPair.Question, startIndex, ok = parser.parseRelations(tokens, startIndex)
					formFound = true
				case field_answer:
					qaPair.Answer, startIndex, ok = parser.parseRelations(tokens, startIndex)
					posFound = true;
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

	return qaPair, startIndex, ok
}

func (parser *simpleInternalGrammarParser) parseLexicon(tokens []SimpleToken, startIndex int) (*parse.SimpleLexicon, int, bool) {

	lexicon := parse.NewSimpleLexicon()
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

func (parser *simpleInternalGrammarParser) parseLexItem(tokens []SimpleToken, startIndex int) (parse.SimpleLexItem, int, bool) {

	lexItem := parse.SimpleLexItem{}
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

func (parser *simpleInternalGrammarParser) parseGrammar(tokens []SimpleToken, startIndex int) (*parse.SimpleGrammar, int, bool) {

	grammar := parse.NewSimpleGrammar()
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)
	for ok {
		lexItem, newStartIndex, ruleFound := parser.parseGrammarRule(tokens, startIndex)
		if ruleFound {
			grammar.AddRule(lexItem)
			startIndex = newStartIndex
		} else {
			ok = false
		}
	}

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_bracket)

	return grammar, startIndex, ok
}

func (parser *simpleInternalGrammarParser) parseGrammarRule(tokens []SimpleToken, startIndex int) (parse.SimpleGrammarRule, int, bool) {

	rule := parse.SimpleGrammarRule{}
	ok, ruleFound := true, false

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_brace)

	for ok {
		field := ""
		field, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_colon)
			if ok {
				switch field {
				case field_rule:
					rule.SyntacticCategories, rule.EntityVariables, startIndex, ok = parser.parseSyntacticRewriteRule(tokens, startIndex)
					ruleFound = true
				case field_sense:
					rule.Sense, startIndex, ok = parser.parseRelations(tokens, startIndex)
				default:
					ok = false
				}
			}
		}
	}

	// required fields
	if ruleFound {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_brace)
	} else {
		ok = false
	}

	return rule, startIndex, ok
}

func (parser *simpleInternalGrammarParser) parseSyntacticRewriteRule(tokens []SimpleToken, startIndex int) ([]string, []string, int, bool) {

	ok := false
	syntacticCategories := []string{}
	entityVariables := []string{}

	rule, startIndex, ok := parser.parseTransformation(tokens, startIndex)

	// check the constraints on this transformation
	if len(rule.Replacement) != 1 {
		ok = false
	} else if len(rule.Replacement[0].Arguments) != 1 {
		ok = false
	} else {
		for _, patternRelation := range rule.Pattern {
			if len(patternRelation.Arguments) != 1 {
				ok = false
			} else if !patternRelation.Arguments[0].IsVariable() {
				ok = false
			}
		}
	}

	if ok {
		syntacticCategories = append(syntacticCategories, rule.Replacement[0].Predicate)
		entityVariables = append(entityVariables, rule.Replacement[0].Arguments[0].TermValue)

		for _, patternRelation := range rule.Pattern {
			syntacticCategories = append(syntacticCategories, patternRelation.Predicate)
			entityVariables = append(entityVariables, patternRelation.Arguments[0].TermValue)
		}
	}

	return syntacticCategories, entityVariables, startIndex, ok
}

func (parser *simpleInternalGrammarParser) parseRelations(tokens []SimpleToken, startIndex int) ([]mentalese.SimpleRelation, int, bool) {

	relations := []mentalese.SimpleRelation{}
	ok := true
	commaFound := false

	for ok {

		if len(relations) > 0 {
			_, startIndex, commaFound = parser.parseSingleToken(tokens, startIndex, t_comma)
			if !commaFound {
				break;
			}
		}

		relation := mentalese.SimpleRelation{}
		relation, startIndex, ok = parser.parseRelation(tokens, startIndex)
		if ok {
			relations = append(relations, relation)
		}
	}

	return relations, startIndex, ok
}

func (parser *simpleInternalGrammarParser) parseRelation(tokens []SimpleToken, startIndex int) (mentalese.SimpleRelation, int, bool) {

	relation := mentalese.SimpleRelation{}
	ok := true
	commaFound, argumentFound := false, false
	argument := mentalese.SimpleTerm{}
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

func (parser *simpleInternalGrammarParser) parseArgument(tokens []SimpleToken, startIndex int) (mentalese.SimpleTerm, int, bool) {

	ok := false
	tokenValue := ""
	term := mentalese.SimpleTerm{}

	tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
	if ok {
		term.TermType = mentalese.Term_predicateAtom
		term.TermValue = tokenValue
	} else {
		tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_variable)
		if ok {
			term.TermType = mentalese.Term_variable
			term.TermValue = tokenValue
		} else {
			tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_number)
			if ok {
				term.TermType = mentalese.Term_number
				term.TermValue = tokenValue
			} else {
				tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_stringConstant)
				if ok {
					term.TermType = mentalese.Term_stringConstant
					term.TermValue = tokenValue
				} else {
					tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_anonymousVariable)
					if ok {
						term.TermType = mentalese.Term_anonymousVariable
						term.TermValue = tokenValue
					}
				}
			}
		}
	}

	return term, startIndex, ok
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