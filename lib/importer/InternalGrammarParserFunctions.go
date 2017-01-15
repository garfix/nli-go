package importer

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"nli-go/lib/generate"
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
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_semicolon)

		if ok {
			transformations = append(transformations, transformation)
		} else {
			break;
		}
	}

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_bracket)

	return transformations, startIndex, ok
}

// a(A) b(B) := c(A) d(B)
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
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_semicolon)
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

func (parser *InternalGrammarParser) parseGenerationLexicon(tokens []Token, startIndex int) (*generate.GenerationLexicon, int, bool) {

	lexicon := generate.NewGenerationLexicon()
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)
	for ok {
		lexItem, newStartIndex, lexItemFound := parser.parseGenerationLexItem(tokens, startIndex)
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
	ok, done, formFound, posFound := true, false, false, false

	for ok && !done {
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
				if ok {
					_, newStartIndex, separatorFound := parser.parseSingleToken(tokens, startIndex, t_comma)
					if separatorFound {
						startIndex = newStartIndex
					} else {
						_, newStartIndex, separatorFound := parser.parseSingleToken(tokens, startIndex, t_semicolon)
						if separatorFound {
							startIndex = newStartIndex
							done = true
						}
					}
				}

			}
		}
	}

	// required fields
	if !formFound || !posFound {
		ok = false
	}

	return lexItem, startIndex, ok
}

func (parser *InternalGrammarParser) parseGenerationLexItem(tokens []Token, startIndex int) (generate.GenerationLexeme, int, bool) {

	lexItem := generate.GenerationLexeme{}
	ok, done, formFound, posFound := true, false, false, false

	for ok && !done {
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
				case field_condition:
					lexItem.Condition, startIndex, ok = parser.parseRelations(tokens, startIndex)
				default:
					ok = false
				}
				if ok {
					_, newStartIndex, separatorFound := parser.parseSingleToken(tokens, startIndex, t_comma)
					if separatorFound {
						startIndex = newStartIndex
					} else {
						_, newStartIndex, separatorFound := parser.parseSingleToken(tokens, startIndex, t_semicolon)
						if separatorFound {
							startIndex = newStartIndex
							done = true
						}
					}
				}

			}
		}
	}

	// required fields
	if !formFound || !posFound {
		ok = false
	}

	return lexItem, startIndex, ok
}

func (parser *InternalGrammarParser) parseGrammar(tokens []Token, startIndex int) (*parse.Grammar, int, bool) {

	grammar := parse.NewGrammar()
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)
	for ok {
		rule, newStartIndex, ruleFound := parser.parseGrammarRule(tokens, startIndex)
		if ruleFound {
			grammar.AddRule(rule)
			startIndex = newStartIndex
		} else {
			ok = false
		}
	}

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_bracket)

	return grammar, startIndex, ok
}

func (parser *InternalGrammarParser) parseGenerationGrammar(tokens []Token, startIndex int) (*generate.GenerationGrammar, int, bool) {

	grammar := generate.NewGenerationGrammar()
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)
	for ok {
		rule, newStartIndex, ruleFound := parser.parseGenerationGrammarRule(tokens, startIndex)
		if ruleFound {
			grammar.AddRule(rule)
			startIndex = newStartIndex
		} else {
			ok = false
		}
	}

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_bracket)

	return grammar, startIndex, ok
}

// rule: S(S) -> NP(E) VP(S), sense: declaration(S) object(S, E);
func (parser *InternalGrammarParser) parseGrammarRule(tokens []Token, startIndex int) (parse.GrammarRule, int, bool) {

	rule := parse.GrammarRule{}
	ok, ruleFound, done := true, false, false

	for ok && !done {
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
				if ok {
					_, newStartIndex, separatorFound := parser.parseSingleToken(tokens, startIndex, t_comma)
					if separatorFound {
						startIndex = newStartIndex
					} else {
						_, newStartIndex, separatorFound := parser.parseSingleToken(tokens, startIndex, t_semicolon)
						if separatorFound {
							startIndex = newStartIndex
							done = true
						}
					}
				}
			}
		}
	}

	if !ruleFound || !done {
		ok = false
	}

	return rule, startIndex, ok
}

func (parser *InternalGrammarParser) parseGenerationGrammarRule(tokens []Token, startIndex int) (generate.GenerationGrammarRule, int, bool) {

	rule := generate.GenerationGrammarRule{}
	ok, ruleFound, done := true, false, false

	for ok && !done {
		field := ""
		field, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_colon)
			if ok {
				switch field {
				case field_rule:
					rule.Antecedent, rule.Consequents, startIndex, ok = parser.parseSyntacticRewriteRule2(tokens, startIndex)
					ruleFound = true
				case field_condition:
					rule.Condition, startIndex, ok = parser.parseRelations(tokens, startIndex)
				default:
					ok = false
				}
				if ok {
					_, newStartIndex, separatorFound := parser.parseSingleToken(tokens, startIndex, t_comma)
					if separatorFound {
						startIndex = newStartIndex
					} else {
						_, newStartIndex, separatorFound := parser.parseSingleToken(tokens, startIndex, t_semicolon)
						if separatorFound {
							startIndex = newStartIndex
							done = true
						}
					}
				}
			}
		}
	}

	if !ruleFound || !done {
		ok = false
	}

	return rule, startIndex, ok
}

func (parser *InternalGrammarParser) parseSyntacticRewriteRule(tokens []Token, startIndex int) ([]string, []string, int, bool) {

	syntacticCategories := []string{}
	entityVariables := []string{}
	ok := true

	headRelation := mentalese.Relation{}
	headRelation, startIndex, ok = parser.parseRelation(tokens, startIndex)
	if ok {
		ok = len(headRelation.Arguments) == 1
		if ok {

			syntacticCategories = append(syntacticCategories, headRelation.Predicate)
			entityVariables = append(entityVariables, headRelation.Arguments[0].TermValue)

			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_rewrite)
			if ok {
				tailRelations := []mentalese.Relation{}
				tailRelations, startIndex, ok = parser.parseRelations(tokens, startIndex)

				for _, patternRelation := range tailRelations {
					if len(patternRelation.Arguments) == 0 {
						patternRelation.Arguments = []mentalese.Term{{mentalese.Term_variable, "_"}}
					} else if len(patternRelation.Arguments) != 1 {
						ok = false
					} else if !patternRelation.Arguments[0].IsVariable() {
						ok = false
					}
					if ok {
						syntacticCategories = append(syntacticCategories, patternRelation.Predicate)
						entityVariables = append(entityVariables, patternRelation.Arguments[0].TermValue)
					}
				}

			}
		}
	}

	return syntacticCategories, entityVariables, startIndex, ok
}

func (parser *InternalGrammarParser) parseSyntacticRewriteRule2(tokens []Token, startIndex int) (mentalese.Relation, []mentalese.Relation, int, bool) {

	ok := false
	antecedent := mentalese.Relation{}
	consequents := mentalese.RelationSet{}

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
		antecedent = rule.Replacement[0]
		consequents = rule.Pattern
	}

	return antecedent, consequents, startIndex, ok
}

func (parser *InternalGrammarParser) parseRelations(tokens []Token, startIndex int) ([]mentalese.Relation, int, bool) {

	relations := []mentalese.Relation{}
	ok := true

	for ok {

		relation := mentalese.Relation{}
		relation, startIndex, ok = parser.parseRelation(tokens, startIndex)
		if ok {
			relations = append(relations, relation)
		}
	}

	ok = len(relations) > 0

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

		if tokens[startIndex].LineNumber > parser.lastParsedResult.LineNumber {
			parser.lastParsedResult.LineNumber = tokens[startIndex].LineNumber
		}

		ok = (token.TokenId == tokenId)
		if ok {
			tokenValue = token.TokenValue
			startIndex++
		}
	}

	return tokenValue, startIndex, ok
}