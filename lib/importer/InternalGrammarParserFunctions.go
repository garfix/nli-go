package importer

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
)

func (parser *InternalGrammarParser) parseRelationSet(tokens []Token, startIndex int) (mentalese.RelationSet, int, bool) {

	relationSet := mentalese.RelationSet{}
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)

	for startIndex < len(tokens) {
		relation := mentalese.Relation{}
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

func (parser *InternalGrammarParser) parseTransformations(tokens []Token, startIndex int) ([]mentalese.RelationTransformation, int, bool) {

	transformations := []mentalese.RelationTransformation{}
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)

	for startIndex < len(tokens) {
		transformation := mentalese.RelationTransformation{}
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

func (parser *InternalGrammarParser) parseTransformation(tokens []Token, startIndex int) (mentalese.RelationTransformation, int, bool) {

	transformation := mentalese.RelationTransformation{}
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

func (parser *InternalGrammarParser) parseRules(tokens []Token, startIndex int) ([]mentalese.Rule, int, bool) {

	rules := []mentalese.Rule{}
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)

	for startIndex < len(tokens) {
		rule := mentalese.Rule{}
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

func (parser *InternalGrammarParser) parseRule(tokens []Token, startIndex int) (mentalese.Rule, int, bool) {

	rule := mentalese.Rule{}
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

func (parser *InternalGrammarParser) parseQAPairs(tokens []Token, startIndex int) ([]mentalese.QAPair, int, bool) {

	qaPairs := []mentalese.QAPair{}
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)

	for startIndex < len(tokens) {
		qaPair := mentalese.QAPair{}
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

func (parser *InternalGrammarParser) parseQAPair(tokens []Token, startIndex int) (mentalese.QAPair, int, bool) {

	qaPair := mentalese.QAPair{}
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

func (parser *InternalGrammarParser) parseLexicon(tokens []Token, startIndex int) (*parse.Lexicon, int, bool) {

	lexicon := parse.NewLexicon()
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

func (parser *InternalGrammarParser) parseLexItem(tokens []Token, startIndex int) (parse.LexItem, int, bool) {

	lexItem := parse.LexItem{}
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

func (parser *InternalGrammarParser) parseGrammar(tokens []Token, startIndex int) (*parse.Grammar, int, bool) {

	grammar := parse.NewGrammar()
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

func (parser *InternalGrammarParser) parseGrammarRule(tokens []Token, startIndex int) (parse.GrammarRule, int, bool) {

	rule := parse.GrammarRule{}
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

func (parser *InternalGrammarParser) parseSyntacticRewriteRule(tokens []Token, startIndex int) ([]string, []string, int, bool) {

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

func (parser *InternalGrammarParser) parseRelations(tokens []Token, startIndex int) ([]mentalese.Relation, int, bool) {

	relations := []mentalese.Relation{}
	ok := true
	commaFound := false

	for ok {

		if len(relations) > 0 {
			_, startIndex, commaFound = parser.parseSingleToken(tokens, startIndex, t_comma)
			if !commaFound {
				break;
			}
		}

		relation := mentalese.Relation{}
		relation, startIndex, ok = parser.parseRelation(tokens, startIndex)
		if ok {
			relations = append(relations, relation)
		}
	}

	return relations, startIndex, ok
}

func (parser *InternalGrammarParser) parseRelation(tokens []Token, startIndex int) (mentalese.Relation, int, bool) {

	relation := mentalese.Relation{}
	ok := true
	commaFound, argumentFound := false, false
	argument := mentalese.Term{}
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
					argument, startIndex, ok = parser.parseTerm(tokens, startIndex)
					if ok {
						relation.Arguments = append(relation.Arguments, argument)
					}
				}

			} else {

				// first argument (there may not be one, zero arguments are allowed)
				argument, newStartIndex, argumentFound = parser.parseTerm(tokens, startIndex)
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

// {A: 13, B: 'banaan'}
// {}
func (parser *InternalGrammarParser) parseBinding(tokens []Token, startIndex int) (mentalese.Binding, int, bool) {

	binding := mentalese.Binding{}
	ok := true
	commaFound := false
	variable := ""
	value := mentalese.Term{}

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_brace)
	for ok {
		if len(binding) > 0 {
			// second and further bindings
			_, startIndex, commaFound = parser.parseSingleToken(tokens, startIndex, t_comma)
			if !commaFound {
				break;
			}
		} else {
			// check for zero bindings
			_, _, ok = parser.parseSingleToken(tokens, startIndex, t_closing_brace)
			if ok {
				break;
			}
		}

		variable, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_variable)
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_colon)
			if ok {
				value, startIndex, ok = parser.parseTerm(tokens, startIndex)
				if ok {
					binding[variable] = value
				}
			}
		}
	}
	if ok {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_brace)
	}

	return binding, startIndex, ok
}

// [{A:1, B:2} {C:'hello', D:'goodbye'}]
func (parser *InternalGrammarParser) parseBindings(tokens []Token, startIndex int) ([]mentalese.Binding, int, bool) {

	bindings := []mentalese.Binding{}
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)

	for ok {
		binding, newStartIndex, bindingFound := parser.parseBinding(tokens, startIndex)
		if bindingFound {
			bindings = append(bindings, binding)
			startIndex = newStartIndex
		} else {
			break
		}
	}

	if ok {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_bracket)
	}

	return bindings, startIndex, ok
}

func (parser *InternalGrammarParser) parseTerm(tokens []Token, startIndex int) (mentalese.Term, int, bool) {

	ok := false
	tokenValue := ""
	term := mentalese.Term{}

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
func (parser *InternalGrammarParser) parseSingleToken(tokens []Token, startIndex int, tokenId int) (string, int, bool) {

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