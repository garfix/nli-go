package importer

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"strings"
)

func (parser *InternalGrammarParser) parseRelationSet(tokens []Token, startIndex int) (mentalese.RelationSet, int, bool) {

	relationSet := mentalese.RelationSet{}
	ok := false

	for startIndex < len(tokens) {
		relation, newStartIndex, found := parser.parseRelation(tokens, startIndex)
		if found {
			relationSet = append(relationSet, relation)
			startIndex = newStartIndex
			ok = true
		} else {
			break
		}
	}

	if !ok {
		tokenValue, newStartIndex, found := parser.parseSingleToken(tokens, startIndex, t_predicate)
		if found {
			if tokenValue == mentalese.PredicateNone {
				startIndex = newStartIndex
				ok = true
			}
		}
	}

	return relationSet, startIndex, ok
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
			break
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
			break
		}
	}

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_bracket)

	return solutions, startIndex, ok
}

func (parser *InternalGrammarParser) parseMap(tokens []Token, startIndex int, parseCustomValue func(tokens []Token, startIndex int, key string) (int, bool, bool)) (int, bool) {

	ok, done, allRequiredItemsFound := true, false, false

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_brace)

	for ok && !done {
		field := ""
		field, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_colon)
			if ok {
				startIndex, ok, allRequiredItemsFound = parseCustomValue(tokens, startIndex, field)
				if ok {
					_, newStartIndex, separatorFound := parser.parseSingleToken(tokens, startIndex, t_comma)
					if separatorFound {
						startIndex = newStartIndex
					} else {
						done = true
					}
				}
			}
		}
	}

	if ok {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_brace)
	}

	if !allRequiredItemsFound {
		ok = false
	}

	return startIndex, ok
}

func (parser *InternalGrammarParser) parseSolution(tokens []Token, startIndex int) (mentalese.Solution, int, bool) {

	solution := mentalese.Solution{}
	solution.Transformations = []mentalese.Rule{}
	conditionFound, responsesFound, resultFound := false, false, false

	callback := func(tokens []Token, startIndex int, key string) (int, bool, bool) {

		ok := true

		switch key {
		case field_condition:
			solution.Condition, startIndex, ok = parser.parseRelations(tokens, startIndex)
			ok = ok && !conditionFound
			conditionFound = true
		case field_transformations:
			solution.Transformations, startIndex, ok = parser.parseRules(tokens, startIndex)
		case field_result:
			solution.Result, startIndex, ok = parser.parseTerm(tokens, startIndex)
			ok = ok && solution.Result.IsVariable() || solution.Result.IsAnonymousVariable()
			ok = ok && !resultFound
			resultFound = true
		case field_responses:
			solution.Responses, startIndex, ok = parser.parseResponses(tokens, startIndex)
			ok = ok && !responsesFound
			responsesFound = true
		default:
			ok = false
		}

		return startIndex, ok, conditionFound && responsesFound && resultFound
	}

	startIndex, ok := parser.parseMap(tokens, startIndex, callback)

	return solution, startIndex, ok
}

func (parser *InternalGrammarParser) parseResponses(tokens []Token, startIndex int) ([]mentalese.ResultHandler, int, bool) {

	handlers := []mentalese.ResultHandler{}
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)

	for startIndex < len(tokens) {
		rule := mentalese.ResultHandler{}
		rule, startIndex, ok = parser.parseResultHandler(tokens, startIndex)
		if ok {
			handlers = append(handlers, rule)
		} else {
			break
		}
	}

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_bracket)

	return handlers, startIndex, ok
}

func (parser *InternalGrammarParser) parseResultHandler(tokens []Token, startIndex int) (mentalese.ResultHandler, int, bool) {

	resultHandler := mentalese.ResultHandler{}
	preparationFound, answerFound := false, false

	callback := func(tokens []Token, startIndex int, key string) (int, bool, bool) {

		ok := true

		switch key {
		case field_condition:
			resultHandler.Condition, startIndex, ok = parser.parseRelations(tokens, startIndex)
		case field_preparation:
			resultHandler.Preparation, startIndex, ok = parser.parseRelations(tokens, startIndex)
			ok = ok && !preparationFound
			preparationFound = true
		case field_answer:
			resultHandler.Answer, startIndex, ok = parser.parseRelations(tokens, startIndex)
			ok = ok && !answerFound
			answerFound = true
		default:
			ok = false
		}

		return startIndex, ok, answerFound
	}

	startIndex, ok := parser.parseMap(tokens, startIndex, callback)

	return resultHandler, startIndex, ok
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

func (parser *InternalGrammarParser) parseGenerationGrammar(tokens []Token, startIndex int) (*parse.Grammar, int, bool) {

	grammar := parse.NewGrammar()
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
	ruleFound, senseFound := false, false

	callback := func(tokens []Token, startIndex int, key string) (int, bool, bool) {

		ok := true

		switch key {
			case field_rule:
				rule.SyntacticCategories, rule.EntityVariables, rule.PositionTypes, startIndex, ok = parser.parseSyntacticRewriteRule(tokens, startIndex)
				ok = ok && !ruleFound
				ruleFound = true
			case field_sense:
				rule.Sense, startIndex, ok = parser.parseRelations(tokens, startIndex)
				ok = ok && !senseFound
				senseFound = true
			default:
				ok = false
		}

		return startIndex, ok, ruleFound
	}

	startIndex, ok := parser.parseMap(tokens, startIndex, callback)

	return rule, startIndex,  ok
}

func (parser *InternalGrammarParser) parseGenerationGrammarRule(tokens []Token, startIndex int) (parse.GrammarRule, int, bool) {

	rule := parse.GrammarRule{}
	ruleFound, conditionFound := false, false

	callback := func(tokens []Token, startIndex int, key string) (int, bool, bool) {

		ok := true

		switch key {
			case field_rule:
				rule.SyntacticCategories, rule.EntityVariables, rule.PositionTypes, startIndex, ok = parser.parseSyntacticRewriteRule(tokens, startIndex)
				ok = ok && !ruleFound
				ruleFound = true
			case field_condition:
				rule.Sense, startIndex, ok = parser.parseRelations(tokens, startIndex)
				ok = ok && !conditionFound
				conditionFound = true
			default:
				ok = false
		}

		return startIndex, ok, ruleFound
	}

	startIndex, ok := parser.parseMap(tokens, startIndex, callback)

	return rule, startIndex,  ok
}

func (parser *InternalGrammarParser) parseSyntacticRewriteRule(tokens []Token, startIndex int) ([]string, [][]string, []string, int, bool) {

	syntacticCategories := []string{}
	entityVariables := [][]string{}
	positionTypes := []string{}
	list := []string{}
	ok := true

	headRelation := mentalese.Relation{}
	headRelation, startIndex, ok = parser.parseRelation(tokens, startIndex)
	if ok {
		syntacticCategories = append(syntacticCategories, headRelation.Predicate)
		positionTypes = append(positionTypes, parse.PosTypeRelation)

		list = []string{}
		for _, argument := range headRelation.Arguments {
			list = append(list, argument.TermValue)
		}
		entityVariables = append(entityVariables, list)

		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_rewrite)
		for ok {
			tailRelation, newStartIndex, isRelation := parser.parseRelation(tokens, startIndex)
			if isRelation {
				startIndex = newStartIndex
				syntacticCategories = append(syntacticCategories, tailRelation.Predicate)
				positionTypes = append(positionTypes, parse.PosTypeRelation)

				list = []string{}
				for _, argument := range tailRelation.Arguments {
					list = append(list, argument.TermValue)
				}
				entityVariables = append(entityVariables, list)
			} else {
				tailString, newStartIndex, isString := parser.parseSingleToken(tokens, startIndex, t_stringConstant)
				if isString {
					startIndex = newStartIndex
					syntacticCategories = append(syntacticCategories, tailString)
					positionTypes = append(positionTypes, parse.PosTypeWordForm)
					entityVariables = append(entityVariables, []string{})
				} else {
					tailRegExp, newStartIndex, isRegExp := parser.parseSingleToken(tokens, startIndex, t_regExp)
					if isRegExp {
						startIndex = newStartIndex
						syntacticCategories = append(syntacticCategories, tailRegExp)
						positionTypes = append(positionTypes, parse.PosTypeRegExp)
						entityVariables = append(entityVariables, []string{})
					} else {
						break
					}
				}
			}
		}

		ok = ok && len(syntacticCategories) > 1
	}

	return syntacticCategories, entityVariables, positionTypes, startIndex, ok
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

	ok := true
	commaFound, argumentFound := false, false
	argument := mentalese.Term{}
	arguments := []mentalese.Term{}
	predicate := ""
	newStartIndex := 0
	positive := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_negative)
	if ok {
		positive = false
	}
	predicate, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
	if ok {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_parenthesis)
		for ok {
			if len(arguments) > 0 {

				// second and further arguments
				_, startIndex, commaFound = parser.parseSingleToken(tokens, startIndex, t_comma)
				if !commaFound {
					break
				} else {
					argument, startIndex, ok = parser.parseTerm(tokens, startIndex)
					if ok {
						arguments = append(arguments, argument)
					}
				}

			} else {

				// first argument (there may not be one, zero arguments are allowed)
				argument, newStartIndex, argumentFound = parser.parseTerm(tokens, startIndex)
				if !argumentFound {
					break
				} else {
					arguments = append(arguments, argument)
					startIndex = newStartIndex
				}

			}
		}
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_parenthesis)
		}
	}

	relation := mentalese.NewRelation(positive, predicate, arguments)

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
				break
			}
		} else {
			// check for zero bindings
			_, _, ok = parser.parseSingleToken(tokens, startIndex, t_closing_brace)
			if ok {
				break
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
func (parser *InternalGrammarParser) parseBindings(tokens []Token, startIndex int) (mentalese.Bindings, int, bool) {

	bindings := mentalese.Bindings{}
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

func (parser *InternalGrammarParser) parseId(tokens []Token, startIndex int) (string, string, int, bool) {

	id := ""
	entityType := ""

	token, newStartIndex, ok := parser.parseSingleToken(tokens, startIndex, t_id)
	if ok {
		i := strings.Index(token, ":")
		if i == -1 {
			ok = false
		} else {
			startIndex = newStartIndex
			entityType = token[0:i]
			id = token[i+1:]
		}
	}

	return id, entityType, startIndex, ok
}

func (parser *InternalGrammarParser) parseTermList(tokens []Token, startIndex int) (mentalese.TermList, int, bool) {

	list := mentalese.TermList{}
	term := mentalese.Term{}
	ok := false
	newStartIndex := 0

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_opening_bracket)
	if ok {
		term, newStartIndex, ok = parser.parseTerm(tokens, startIndex)
		if ok {
			list = append(list, term)
			startIndex = newStartIndex

			for startIndex < len(tokens) {
				_, newStartIndex, ok = parser.parseSingleToken(tokens, startIndex, t_comma)
				if ok {
					startIndex = newStartIndex
					term, newStartIndex, ok = parser.parseTerm(tokens, startIndex)
					if ok {
						list = append(list, term)
						startIndex = newStartIndex
					} else {
						goto end
					}
				} else {
					break
				}
			}
		}
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_closing_bracket)
	}

	end:

	return list, startIndex, ok
}

func (parser *InternalGrammarParser) parseTerm(tokens []Token, startIndex int) (mentalese.Term, int, bool) {

	ok := false
	tokenValue := ""
	term := mentalese.Term{}
	newStartIndex := 0

	tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_variable)
	if ok {
		term.TermType = mentalese.TermTypeVariable
		term.TermValue = tokenValue
	} else {
		tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_number)
		if ok {
			term.TermType = mentalese.TermTypeStringConstant
			term.TermValue = tokenValue
		} else {
			tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_stringConstant)
			if ok {
				term.TermType = mentalese.TermTypeStringConstant
				term.TermValue = tokenValue
			} else {
				tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_anonymousVariable)
				if ok {
					term.TermType = mentalese.TermTypeAnonymousVariable
					term.TermValue = tokenValue
				} else {
					id := ""
					entityType := ""
					id, entityType, startIndex, ok = parser.parseId(tokens, startIndex)
					if ok {
						term.TermType = mentalese.TermTypeId
						term.TermValue = id
						term.TermEntityType = entityType
					} else {
						rule := mentalese.Rule{}
						rule, newStartIndex, ok = parser.parseRule(tokens, startIndex)
						if ok {
							term.TermType = mentalese.TermTypeRule
							term.TermValueRule = rule
							startIndex = newStartIndex
						} else {
							relationSet := mentalese.RelationSet{}
							relationSet, newStartIndex, ok = parser.parseRelationSet(tokens, startIndex)
							if ok {
								term.TermType = mentalese.TermTypeRelationSet
								term.TermValueRelationSet = relationSet
								startIndex = newStartIndex
							} else {
								tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_predicate)
								if ok {
									term.TermType = mentalese.TermTypePredicateAtom
									term.TermValue = tokenValue
								} else {
									list := mentalese.TermList{}
									list, newStartIndex, ok = parser.parseTermList(tokens, startIndex)
									if ok {
										term.TermType = mentalese.TermTypeList
										term.TermValueList = list
										startIndex = newStartIndex
									} else {
										tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, t_regExp)
										if ok {
											term.TermType = mentalese.TermTypeRegExp
											term.TermValue = tokenValue
										}
									}
								}
							}
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

		ok = token.TokenId == tokenId
		if ok {
			tokenValue = token.TokenValue
			startIndex++
		}
	}

	return tokenValue, startIndex, ok
}
