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
	if ok {
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
	}

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

func (parser *InternalGrammarParser) parseSolutions(tokens []Token, startIndex int) ([]mentalese.Solution, int, bool) {

	solutions := []mentalese.Solution{}
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)

	for startIndex < len(tokens) {
		solution := mentalese.Solution{}
		solution, startIndex, ok = parser.parseSolution(tokens, startIndex)
		if ok {
			solutions = append(solutions, solution)
		} else {
			break;
		}
	}

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_bracket)

	return solutions, startIndex, ok
}

func (parser *InternalGrammarParser) parseSolution(tokens []Token, startIndex int) (mentalese.Solution, int, bool) {

	solution := mentalese.Solution{}
	ok, done, conditionFound, prepationFound, answerFound := true, false, false, false, false

	for ok && !done {
		field := ""
		field, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_colon)
			if ok {
				switch field {
				case field_condition:
					solution.Condition, startIndex, ok = parser.parseRelations(tokens, startIndex)
					if conditionFound {
						ok = false
					}
					conditionFound = true
				case field_preparation:
					solution.Preparation, startIndex, ok = parser.parseRelations(tokens, startIndex)
					if prepationFound {
						ok = false
					}
					prepationFound = true;
				case field_answer:
					solution.Answer, startIndex, ok = parser.parseRelations(tokens, startIndex)
					if answerFound {
						ok = false
					}
					answerFound = true
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
						} else {
							ok = false
						}
					}
				}

			}
		}
	}

	// required fields
	if !conditionFound || !answerFound {
		ok = false
	}

	return solution, startIndex, ok
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
	ok, done, formFound, senseFound, posFound := true, false, false, false, false

	for ok && !done {
		field := ""
		field, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_colon)
			if ok {
				switch field {
				case field_form:
					if formFound {
						ok = false
					} else {
						formFound = true
						lexItem.Form, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_stringConstant)
						if ok {
							lexItem.IsRegExp = false
						} else {
							lexItem.Form, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_regExp)
							lexItem.IsRegExp = true
						}
					}
				case field_pos:
					lexItem.PartOfSpeech, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
					if posFound {
						ok = false
					}
					posFound = true;
				case field_sense:
					lexItem.RelationTemplates, startIndex, ok = parser.parseRelations(tokens, startIndex)
					if senseFound {
						ok = false
					}
					senseFound = true
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
						} else {
							ok = false
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
	ok, done, formFound, posFound, conditionFound := true, false, false, false, false

	for ok && !done {
		field := ""
		field, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_colon)
			if ok {
				switch field {
				case field_form:
					if formFound {
						ok = false
					} else {
						formFound = true
						lexItem.Form, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_stringConstant)
						if ok {
							lexItem.IsRegExp = false
						} else {
							lexItem.Form, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_regExp)
							lexItem.IsRegExp = true
						}
					}
				case field_pos:
					lexItem.PartOfSpeech, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
					if posFound {
						ok = false
					}
					posFound = true;
				case field_condition:
					lexItem.Condition, startIndex, ok = parser.parseRelations(tokens, startIndex)
					if conditionFound {
						ok = false
					}
					conditionFound = true
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
						} else {
							ok = false
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

// rule: s(S) -> np(E) vp(S), sense: declaration(S) object(S, E);
func (parser *InternalGrammarParser) parseGrammarRule(tokens []Token, startIndex int) (parse.GrammarRule, int, bool) {

	rule := parse.GrammarRule{}
	ok, ruleFound, senseFound, done := true, false, false, false

	for ok && !done {
		field := ""
		field, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_colon)
			if ok {
				switch field {
				case field_rule:
					rule.SyntacticCategories, rule.EntityVariables, startIndex, ok = parser.parseSyntacticRewriteRule(tokens, startIndex)
					if ruleFound {
						ok = false
					}
					ruleFound = true
				case field_sense:
					rule.Sense, startIndex, ok = parser.parseRelations(tokens, startIndex)
					if senseFound {
						ok = false
					}
					senseFound = true
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
						} else {
							ok = false
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
	ok, ruleFound, conditionFound, done := true, false, false, false

	for ok && !done {
		field := ""
		field, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_colon)
			if ok {
				switch field {
				case field_rule:
					rule.Antecedent, rule.Consequents, startIndex, ok = parser.parseSyntacticRewriteRule2(tokens, startIndex)
					if ruleFound {
						ok = false
					}
					ruleFound = true
				case field_condition:
					rule.Condition, startIndex, ok = parser.parseRelations(tokens, startIndex)
					if conditionFound {
						ok = false
					}
					conditionFound = true
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
						} else {
							ok = false
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

	antecedent, startIndex, ok = parser.parseRelation(tokens, startIndex)
	if ok {
		ok = len(antecedent.Arguments) == 1
		if ok {

			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_rewrite)
			if ok {
				consequents, startIndex, ok = parser.parseRelations(tokens, startIndex)

				for _, consequent := range consequents {
					if len(consequent.Arguments) == 0 {
						consequent.Arguments = []mentalese.Term{{mentalese.Term_variable, "_"}}
					} else if len(consequent.Arguments) != 1 {
						ok = false
					} else if !consequent.Arguments[0].IsVariable() {
						ok = false
					}
				}

			}
		}
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
					} else {
						tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_regExp)
						if ok {
							term.TermType = mentalese.Term_regExp
							term.TermValue = tokenValue
						}
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