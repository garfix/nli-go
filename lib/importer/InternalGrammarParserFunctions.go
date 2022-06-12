package importer

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse/morphology"
	"regexp"
	"strconv"
	"strings"
)

func (parser *InternalGrammarParser) parseRules(tokens []Token, startIndex int) ([]mentalese.Rule, int, bool) {

	rules := []mentalese.Rule{}
	ok := true

	for startIndex < len(tokens) {
		rule := mentalese.Rule{}
		rule, startIndex, ok = parser.parseRule(tokens, startIndex)
		if !ok {
			break
		}
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tSemicolon)
		if !ok {
			break
		}
		rules = append(rules, rule)
	}

	return rules, startIndex, len(rules) > 0
}

func (parser *InternalGrammarParser) parseRule(tokens []Token, startIndex int) (mentalese.Rule, int, bool) {

	newStartIndex := 0
	rule := mentalese.Rule{}
	ok := true

	rule.Goal, startIndex, ok = parser.parseRelation(tokens, startIndex, true)
	if ok {
		_, newStartIndex, ok = parser.parseSingleToken(tokens, startIndex, tImplication)
		if ok {
			startIndex = newStartIndex
			rule.Pattern, startIndex, ok = parser.parseRelations(tokens, startIndex, true)
		} else {
			startIndex = newStartIndex
			ok = true
		}
	}
	if ok {
		for _, argument := range rule.Goal.Arguments {
			if argument.IsMutableVariable() {
				ok = false
			}
		}
	}

	return rule, startIndex, ok
}

func (parser *InternalGrammarParser) parseSolutions(tokens []Token, startIndex int) ([]mentalese.Solution, int, bool) {

	solutions := []mentalese.Solution{}
	ok := true

	for startIndex < len(tokens) {
		solution := mentalese.Solution{}
		solution, startIndex, ok = parser.parseSolution(tokens, startIndex)
		if ok {
			solutions = append(solutions, solution)
		} else {
			break
		}
	}

	return solutions, startIndex, len(solutions) > 0
}

func (parser *InternalGrammarParser) parseMap(tokens []Token, startIndex int, parseCustomValue func(tokens []Token, startIndex int, key string) (int, bool, bool)) (int, bool) {

	ok, done, allRequiredItemsFound := true, false, false

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tOpeningBrace)

	for ok && !done {
		field := ""
		field, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tPredicate)
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tColon)
			if ok {
				startIndex, ok, allRequiredItemsFound = parseCustomValue(tokens, startIndex, field)
				if ok {
					_, newStartIndex, separatorFound := parser.parseSingleToken(tokens, startIndex, tComma)
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
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tClosingBrace)
	}

	if !allRequiredItemsFound {
		ok = false
	}

	return startIndex, ok
}

func (parser *InternalGrammarParser) parseSolution(tokens []Token, startIndex int) (mentalese.Solution, int, bool) {

	solution := mentalese.Solution{}
	solution.Transformations = []mentalese.Rule{}
	conditionFound, responsesFound := false, false

	callback := func(tokens []Token, startIndex int, key string) (int, bool, bool) {

		ok := true

		switch key {
		case field_condition:
			solution.Condition, startIndex, ok = parser.parseRelations(tokens, startIndex, true)
			ok = ok && !conditionFound
			conditionFound = true
		case field_transformations:
			solution.Transformations, startIndex, ok = parser.parseRules(tokens, startIndex)
		case field_responses:
			solution.Responses, startIndex, ok = parser.parseResponses(tokens, startIndex)
			ok = ok && !responsesFound
			responsesFound = true
		default:
			ok = false
		}

		return startIndex, ok, conditionFound && responsesFound
	}

	startIndex, ok := parser.parseMap(tokens, startIndex, callback)

	return solution, startIndex, ok
}

func (parser *InternalGrammarParser) parseResponses(tokens []Token, startIndex int) ([]mentalese.ResultHandler, int, bool) {

	handlers := []mentalese.ResultHandler{}
	ok := true

	for startIndex < len(tokens) {
		rule := mentalese.ResultHandler{}
		rule, startIndex, ok = parser.parseResultHandler(tokens, startIndex)
		if ok {
			handlers = append(handlers, rule)
		} else {
			break
		}
	}

	return handlers, startIndex, len(handlers) > 0
}

func (parser *InternalGrammarParser) parseResultHandler(tokens []Token, startIndex int) (mentalese.ResultHandler, int, bool) {

	resultHandler := mentalese.ResultHandler{}
	preparationFound, answerFound, resultFound := false, false, false

	callback := func(tokens []Token, startIndex int, key string) (int, bool, bool) {

		ok := true

		switch key {
		case field_condition:
			resultHandler.Condition, startIndex, ok = parser.parseRelations(tokens, startIndex, true)
		case field_preparation:
			resultHandler.Preparation, startIndex, ok = parser.parseRelations(tokens, startIndex, true)
			ok = ok && !preparationFound
			preparationFound = true
		case field_answer:
			resultHandler.Answer, startIndex, ok = parser.parseRelations(tokens, startIndex, true)
			ok = ok && !answerFound
			answerFound = true
		case field_result:
			resultHandler.Result, startIndex, ok = parser.parseTerm(tokens, startIndex)
			ok = ok && resultHandler.Result.IsVariable() || resultHandler.Result.IsAnonymousVariable()
			ok = ok && !resultFound
			resultFound = true
		default:
			ok = false
		}

		return startIndex, ok, answerFound
	}

	startIndex, ok := parser.parseMap(tokens, startIndex, callback)

	return resultHandler, startIndex, ok
}

func (parser *InternalGrammarParser) parseGrammar(tokens []Token, startIndex int) (*mentalese.GrammarRules, int, bool) {

	grammar := mentalese.NewGrammarRules()
	ok := true
	found := false

	for ok {
		rule, newStartIndex, ruleFound := parser.parseGrammarRule(tokens, startIndex)
		if ruleFound {
			grammar.AddRule(rule)
			startIndex = newStartIndex
			found = true
		} else {
			ok = false
		}
	}

	return grammar, startIndex, found
}

func (parser *InternalGrammarParser) parseGenerationGrammar(tokens []Token, startIndex int) (*mentalese.GrammarRules, int, bool) {

	grammar := mentalese.NewGrammarRules()
	ok := true
	found := false

	for ok {
		rule, newStartIndex, ruleFound := parser.parseGenerationGrammarRule(tokens, startIndex)
		if ruleFound {
			grammar.AddRule(rule)
			startIndex = newStartIndex
			found = true
		} else {
			ok = false
		}
	}

	return grammar, startIndex, found
}

// rule: s(S) -> np(E) vp(S), sense: declaration(S) object(S, E);
func (parser *InternalGrammarParser) parseGrammarRule(tokens []Token, startIndex int) (mentalese.GrammarRule, int, bool) {

	rule := mentalese.GrammarRule{}
	ruleFound, senseFound, ellipsisFound, tagFound := false, false, false, false

	callback := func(tokens []Token, startIndex int, key string) (int, bool, bool) {

		ok := true

		switch key {
		case field_rule:
			rule.SyntacticCategories, rule.EntityVariables, rule.PositionTypes, startIndex, ok = parser.parseSyntacticRewriteRule(tokens, startIndex)
			ok = ok && !ruleFound
			ruleFound = true
		case field_sense:
			rule.Sense, startIndex, ok = parser.parseRelations(tokens, startIndex, true)
			ok = ok && !senseFound
			senseFound = true
		case field_ellipsis:
			rule.Ellipsis, startIndex, ok = parser.parseCategoryPaths(tokens, startIndex)
			ok = ok && !ellipsisFound
			ellipsisFound = true
		case field_tag:
			rule.Tag, startIndex, ok = parser.parseRelations(tokens, startIndex, true)
			ok = ok && !tagFound
			tagFound = true
		default:
			ok = false
		}

		return startIndex, ok, ruleFound
	}

	startIndex, ok := parser.parseMap(tokens, startIndex, callback)

	return rule, startIndex, ok
}

func (parser *InternalGrammarParser) parseGenerationGrammarRule(tokens []Token, startIndex int) (mentalese.GrammarRule, int, bool) {

	rule := mentalese.GrammarRule{}
	ruleFound, conditionFound := false, false

	callback := func(tokens []Token, startIndex int, key string) (int, bool, bool) {

		ok := true

		switch key {
		case field_rule:
			rule.SyntacticCategories, rule.EntityVariables, rule.PositionTypes, startIndex, ok = parser.parseSyntacticRewriteRule(tokens, startIndex)
			ok = ok && !ruleFound
			ruleFound = true
		case field_condition:
			rule.Sense, startIndex, ok = parser.parseRelations(tokens, startIndex, true)
			ok = ok && !conditionFound
			conditionFound = true
		default:
			ok = false
		}

		return startIndex, ok, ruleFound
	}

	startIndex, ok := parser.parseMap(tokens, startIndex, callback)

	return rule, startIndex, ok
}

func (parser *InternalGrammarParser) parseSyntacticRewriteRule(tokens []Token, startIndex int) ([]string, [][]string, []string, int, bool) {

	syntacticCategories := []string{}
	entityVariables := [][]string{}
	positionTypes := []string{}
	list := []string{}
	ok := true

	headRelation := mentalese.Relation{}
	headRelation, startIndex, ok = parser.parseRelation(tokens, startIndex, false)
	if ok {
		syntacticCategories = append(syntacticCategories, headRelation.Predicate)
		positionTypes = append(positionTypes, mentalese.PosTypeRelation)

		list = []string{}
		for _, argument := range headRelation.Arguments {
			list = append(list, argument.TermValue)
		}
		entityVariables = append(entityVariables, list)

		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tRewrite)
		for ok {
			tailRelation, newStartIndex, isRelation := parser.parseRelation(tokens, startIndex, false)
			if isRelation {
				startIndex = newStartIndex
				syntacticCategories = append(syntacticCategories, tailRelation.Predicate)
				positionTypes = append(positionTypes, mentalese.PosTypeRelation)

				list = []string{}
				for _, argument := range tailRelation.Arguments {
					list = append(list, argument.TermValue)
				}
				entityVariables = append(entityVariables, list)
			} else {
				tailString, newStartIndex, isString := parser.parseSingleToken(tokens, startIndex, tStringConstant)
				if isString {
					startIndex = newStartIndex
					syntacticCategories = append(syntacticCategories, tailString)
					positionTypes = append(positionTypes, mentalese.PosTypeWordForm)
					entityVariables = append(entityVariables, []string{})
				} else {
					tailRegExp, newStartIndex, isRegExp := parser.parseSingleToken(tokens, startIndex, tRegExp)
					if isRegExp {
						startIndex = newStartIndex
						syntacticCategories = append(syntacticCategories, tailRegExp)
						positionTypes = append(positionTypes, mentalese.PosTypeRegExp)
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

func (parser *InternalGrammarParser) parseCategoryPaths(tokens []Token, startIndex int) ([]mentalese.CategoryPath, int, bool) {
	paths := []mentalese.CategoryPath{}
	ok := true

	for ok {
		path, newStartIndex, found := parser.parseCategoryPath(tokens, startIndex)
		if found {
			paths = append(paths, path)
			startIndex = newStartIndex
		} else {
			ok = false
		}
	}

	return paths, startIndex, true
}

func (parser *InternalGrammarParser) parseCategoryPath(tokens []Token, startIndex int) (mentalese.CategoryPath, int, bool) {
	path := mentalese.CategoryPath{}
	slash := true

	node, newStartIndex, ok := parser.parseCategoryPathNode(tokens, startIndex)
	if ok {
		path = append(path, node)
		startIndex = newStartIndex

		for ok {
			_, newStartIndex, slash = parser.parseSingleToken(tokens, startIndex, tSlash)
			if slash {
				startIndex = newStartIndex
			} else {
				break
			}

			node, newStartIndex, ok = parser.parseCategoryPathNode(tokens, startIndex)
			if ok {
				path = append(path, node)
				startIndex = newStartIndex
			}
		}
	}

	return path, startIndex, ok
}

func (parser *InternalGrammarParser) parseCategoryPathNode(tokens []Token, startIndex int) (mentalese.CategoryPathNode, int, bool) {
	node := mentalese.CategoryPathNode{}
	nodeType := ""
	newStartIndex := 0
	ok := true
	categoryRequired := false
	category := ""
	variables := []string{}
	allowIndirect := false

	_, newStartIndex, ok = parser.parseSingleToken(tokens, startIndex, tUp)
	if ok {
		startIndex = newStartIndex
		nodeType = mentalese.NodeTypeParent
		goto category
	}

	_, newStartIndex, ok = parser.parseSingleToken(tokens, startIndex, tPositive)
	if ok {
		startIndex = newStartIndex
		nodeType = mentalese.NodeTypeNextSibling

		_, newStartIndex, ok = parser.parseSingleToken(tokens, startIndex, tNegative)
		if ok {
			startIndex = newStartIndex
			nodeType = mentalese.NodeTypeSibling
		}

		ok = true
		goto category
	}

	_, newStartIndex, ok = parser.parseSingleToken(tokens, startIndex, tNegative)
	if ok {
		startIndex = newStartIndex
		nodeType = mentalese.NodeTypePrevSibling
		goto category
	}

	_, newStartIndex, ok = parser.parseSingleToken(tokens, startIndex, tOpeningBracket)
	if ok {
		startIndex = newStartIndex
		nodeType, newStartIndex, ok = parser.parseSingleToken(tokens, startIndex, tPredicate)
		if ok {
			startIndex = newStartIndex
			_, newStartIndex, ok = parser.parseSingleToken(tokens, startIndex, tClosingBracket)
			if ok {
				startIndex = newStartIndex
				if !common.StringArrayContains([]string{mentalese.NodeTypePrevSentence}, nodeType) {
					ok = false
				}
				goto done
			}
		}
	}

	nodeType = mentalese.NodeTypeChild
	categoryRequired = true
	ok = true

	_, newStartIndex, allowIndirect = parser.parseSingleToken(tokens, startIndex, tSlash)
	if allowIndirect {
		startIndex = newStartIndex
	}

category:

	{
		relation, newStartIndex, ok := parser.parseRelation(tokens, startIndex, false)
		if ok {
			category = relation.Predicate
			startIndex = newStartIndex
			for _, argument := range relation.Arguments {
				variables = append(variables, argument.TermValue)
			}
		}
	}

done:

	node = mentalese.NewCategoryPathNode(nodeType, category, variables, allowIndirect)

	if categoryRequired && (category == "") {
		ok = false
	}

	return node, startIndex, ok
}

func (parser *InternalGrammarParser) parseRelations(tokens []Token, startIndex int, useAlias bool) ([]mentalese.Relation, int, bool) {

	relation := mentalese.Relation{}
	relationSet := mentalese.RelationSet{}
	newStartIndex := 0
	found := false
	ok := false

	for startIndex < len(tokens) {

		relation, newStartIndex, found = parser.parseRelationTag(tokens, startIndex)
		if found {
			relationSet = append(relationSet, relation)
			startIndex = newStartIndex
			ok = true
			continue
		}

		relation, newStartIndex, found = parser.parseRelation(tokens, startIndex, true)
		if found {
			relationSet = append(relationSet, relation)
			startIndex = newStartIndex
			ok = true
			continue
		}
		break
	}

	if !ok {
		tokenValue, newStartIndex, found := parser.parseSingleToken(tokens, startIndex, tPredicate)
		if found {
			if tokenValue == mentalese.AtomNone {
				startIndex = newStartIndex
				ok = true
			}
		}
	}

	_, _, found = parser.parseSingleToken(tokens, startIndex, tImplication)
	if found {
		ok = false
	}

	return relationSet, startIndex, ok
}

func (parser *InternalGrammarParser) parseRelationTag(tokens []Token, startIndex int) (mentalese.Relation, int, bool) {

	found := false
	relation := mentalese.Relation{}
	tag := ""

	_, startIndex, found = parser.parseSingleToken(tokens, startIndex, tDoubleOpeningBrace)
	if found {
		tag, startIndex, found = parser.parseVariable(tokens, startIndex)
		if found {
			relation = mentalese.NewRelation(false, mentalese.PredicateIncludeRelations, []mentalese.Term{
				mentalese.NewTermVariable(tag),
			})
			_, startIndex, found = parser.parseSingleToken(tokens, startIndex, tDoubleClosingBrace)
		}
	}

	return relation, startIndex, found
}

func (parser *InternalGrammarParser) parseSortRelations(tokens []Token, startIndex int) ([]mentalese.SortRelation, int, bool) {

	sortRelations := []mentalese.SortRelation{}
	sortRelation := mentalese.SortRelation{}
	newStartIndex := 0
	ok := true

	for true {
		sortRelation, newStartIndex, ok = parser.parseSortRelation(tokens, startIndex)
		if ok {
			sortRelations = append(sortRelations, sortRelation)
			startIndex = newStartIndex
		} else {
			break
		}
	}

	return sortRelations, startIndex, len(sortRelations) > 0
}

func (parser *InternalGrammarParser) parseSortRelation(tokens []Token, startIndex int) (mentalese.SortRelation, int, bool) {
	sortRelation := mentalese.SortRelation{}
	ok := true
	superSort := ""
	subSort := ""

	superSort, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tPredicate)
	if ok {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tGt)
		if ok {
			subSort, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tPredicate)
			if ok {
				sortRelation = mentalese.NewSortRelation(superSort, subSort)
			}
		}
	}

	return sortRelation, startIndex, ok
}

func (parser *InternalGrammarParser) parseKeyword(tokens []Token, startIndex int, keyword string) (int, bool) {
	value, newStartIndex, ok := parser.parseSingleToken(tokens, startIndex, tPredicate)
	if ok && value == keyword {
		return newStartIndex, true
	} else {
		return 0, false
	}
}

func (parser *InternalGrammarParser) parseKeywordRelation(tokens []Token, startIndex int, useAlias bool) (mentalese.Relation, int, bool) {

	keyword := ""
	ok := false
	ok1 := false
	ok2 := false
	ok3 := false
	ok4 := false
	ok5 := false
	relation := mentalese.Relation{}
	term1 := mentalese.Term{}
	term2 := mentalese.Term{}
	s1 := mentalese.RelationSet{}
	s2 := mentalese.RelationSet{}
	s3 := mentalese.RelationSet{}
	newStartIndex := 0

	keyword, newStartIndex, ok1 = parser.parseSingleToken(tokens, startIndex, tPredicate)
	if ok1 {
		startIndex = newStartIndex
		switch keyword {
		case "if":
			s1, startIndex, ok1 = parser.parseRelations(tokens, startIndex, useAlias)
			startIndex, ok2 = parser.parseKeyword(tokens, startIndex, "then")
			s2, startIndex, ok3 = parser.parseRelations(tokens, startIndex, useAlias)
			newStartIndex, ok4 = parser.parseKeyword(tokens, startIndex, "else")
			if ok4 {
				startIndex = newStartIndex
				s3, startIndex, ok4 = parser.parseRelations(tokens, startIndex, useAlias)
			}
			startIndex, ok5 = parser.parseKeyword(tokens, startIndex, "end")
			ok = ok1 && ok2 && ok3 && ok5
			if ok {
				if ok4 {
					relation = mentalese.NewRelation(false, mentalese.PredicateIfThenElse, []mentalese.Term{
						mentalese.NewTermRelationSet(s1),
						mentalese.NewTermRelationSet(s2),
						mentalese.NewTermRelationSet(s3),
					})
				} else {
					relation = mentalese.NewRelation(false, mentalese.PredicateIfThen, []mentalese.Term{
						mentalese.NewTermRelationSet(s1),
						mentalese.NewTermRelationSet(s2),
					})
				}
			}
		case "return":
			relation = mentalese.NewRelation(false, mentalese.PredicateReturn, []mentalese.Term{})
			ok = true
		case "fail":
			relation = mentalese.NewRelation(false, mentalese.PredicateFail, []mentalese.Term{})
			ok = true
		case "break":
			relation = mentalese.NewRelation(false, mentalese.PredicateBreak, []mentalese.Term{})
			ok = true
		case "cancel":
			relation = mentalese.NewRelation(false, mentalese.PredicateCancel, []mentalese.Term{})
			ok = true
		default:
			ok = false
		}
	} else {
		predicate := ""
		operators := map[int]string{
			tAssign:    mentalese.PredicateAssign,
			tEquals:    mentalese.PredicateEquals,
			tNotEquals: mentalese.PredicateNotEquals,
			tGt:        mentalese.PredicateGreaterThan,
			tGtEq:      mentalese.PredicateGreaterThanEquals,
			tLt:        mentalese.PredicateLessThan,
			tLtEq:      mentalese.PredicateLessThanEquals,
			tPositive:  mentalese.PredicateAdd,
			tNegative:  mentalese.PredicateSubtract,
			tMultiply:  mentalese.PredicateMultiply,
			tSlash:     mentalese.PredicateDivide,
		}
		_, newStartIndex, ok1 = parser.parseSingleToken(tokens, startIndex, tOpeningBracket)
		if ok1 {
			startIndex = newStartIndex
			term1, startIndex, ok2 = parser.parseTerm(tokens, startIndex)
			for operator, aPredicate := range operators {
				_, newStartIndex, ok3 = parser.parseSingleToken(tokens, startIndex, operator)
				if ok3 {
					startIndex = newStartIndex
					predicate = aPredicate

					// extra operand checks
					if predicate == mentalese.PredicateAssign {
						ok3 = term1.IsVariable()
					}
					break
				}
			}
			term2, startIndex, ok4 = parser.parseTerm(tokens, startIndex)
			_, startIndex, ok5 = parser.parseSingleToken(tokens, startIndex, tClosingBracket)
			ok = ok1 && ok2 && ok3 && ok4 && ok5
			if ok {
				if common.StringArrayContains(
					[]string{
						mentalese.PredicateAdd,
						mentalese.PredicateSubtract,
						mentalese.PredicateMultiply,
						mentalese.PredicateDivide},
					predicate) {
					relation = mentalese.NewRelation(false, predicate, []mentalese.Term{
						term1,
						term2,
						mentalese.NewTermAtom(mentalese.AtomReturnValue),
					})
				} else {
					relation = mentalese.NewRelation(false, predicate, []mentalese.Term{
						term1,
						term2,
					})
				}
			}
		}
	}

	return relation, startIndex, ok
}

func (parser *InternalGrammarParser) parseRelation(tokens []Token, startIndex int, useAlias bool) (mentalese.Relation, int, bool) {

	ok := true
	prefix := ""
	commaFound, argumentFound := false, false
	argument := mentalese.Term{}
	arguments := []mentalese.Term{}
	predicate := ""
	newStartIndex := 0
	negate := false
	relation := mentalese.Relation{}
	keywordRelation := mentalese.Relation{}

	keywordRelation, newStartIndex, ok = parser.parseKeywordRelation(tokens, startIndex, useAlias)
	if ok {
		return keywordRelation, newStartIndex, ok
	}

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tNegative)
	if ok {
		negate = true
	}

	relation, newStartIndex, ok = parser.parsePlaceholder(tokens, startIndex, negate)
	if ok {
		startIndex = newStartIndex
	} else {

		if useAlias {
			alias := ""
			possibleAlias := ""
			possibleAlias, newStartIndex, ok = parser.parseSingleToken(tokens, startIndex, tPredicate)
			if ok {
				predicate, newStartIndex, ok = parser.parseSingleToken(tokens, newStartIndex, tColon)
				if ok {
					alias = possibleAlias
					startIndex = newStartIndex
				}
			}

			applicationAlias, found := parser.aliasMap[alias]
			if !found {
				return relation, newStartIndex, false
			} else if applicationAlias != "" {
				prefix = applicationAlias + "_"
			}
		}

		predicate, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tPredicate)
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tOpeningParenthesis)
			for ok {
				if len(arguments) > 0 {

					// second and further arguments
					_, startIndex, commaFound = parser.parseSingleToken(tokens, startIndex, tComma)
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
				_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tClosingParenthesis)
			}

		}
		relation = mentalese.NewRelation(negate, prefix+predicate, arguments)
	}

	return relation, startIndex, ok
}

// {A: 13, B: 'banaan'}
// {}
func (parser *InternalGrammarParser) parseBinding(tokens []Token, startIndex int) (mentalese.Binding, int, bool) {

	binding := mentalese.NewBinding()
	ok := true
	commaFound := false
	variable := ""
	value := mentalese.Term{}

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tOpeningBrace)
	for ok {
		if len(binding.GetAll()) > 0 {
			// second and further bindings
			_, startIndex, commaFound = parser.parseSingleToken(tokens, startIndex, tComma)
			if !commaFound {
				break
			}
		} else {
			// check for zero bindings
			_, _, ok = parser.parseSingleToken(tokens, startIndex, tClosingBrace)
			if ok {
				break
			}
		}

		variable, startIndex, ok = parser.parseVariable(tokens, startIndex)
		if ok {
			_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tColon)
			if ok {
				value, startIndex, ok = parser.parseTerm(tokens, startIndex)
				if ok {
					binding.Set(variable, value)
				}
			}
		}
	}
	if ok {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tClosingBrace)
	}

	return binding, startIndex, ok
}

// [{A:1, B:2} {C:'hello', D:'goodbye'}]
func (parser *InternalGrammarParser) parseBindings(tokens []Token, startIndex int) (mentalese.BindingSet, int, bool) {

	bindings := mentalese.NewBindingSet()
	ok := true

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tOpeningBracket)

	for ok {
		binding, newStartIndex, bindingFound := parser.parseBinding(tokens, startIndex)
		if bindingFound {
			bindings.Add(binding)
			startIndex = newStartIndex
		} else {
			break
		}
	}

	if ok {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tClosingBracket)
	}

	return bindings, startIndex, ok
}

func (parser *InternalGrammarParser) parseId(tokens []Token, startIndex int) (string, string, int, bool) {

	id := ""
	sort := ""

	token, newStartIndex, ok := parser.parseSingleToken(tokens, startIndex, tId)
	if ok {
		i := strings.Index(token, ":")
		if i == -1 {
			ok = false
		} else {
			startIndex = newStartIndex
			sort = token[0:i]
			id = token[i+1:]
		}
	}

	return id, sort, startIndex, ok
}

func (parser *InternalGrammarParser) parsePlaceholder(tokens []Token, startIndex int, positive bool) (mentalese.Relation, int, bool) {

	tokenValue := ""
	placeholder := mentalese.Relation{}
	ok := false

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tPlaceholder)
	if ok {

		tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tPredicate)
		if ok {
			myRegex, _ := regexp.Compile("^([^\\d]+)(\\d*)$")
			result := myRegex.FindStringSubmatch(tokenValue)
			cat := result[1]
			index := result[2]
			if index == "" {
				index = "1"
			}

			placeholder = mentalese.NewRelation(positive, mentalese.PredicateSem, []mentalese.Term{
				mentalese.NewTermAtom(cat),
				mentalese.NewTermString(index),
			})
		}
	}

	return placeholder, startIndex, ok
}

func (parser *InternalGrammarParser) parseTermList(tokens []Token, startIndex int) (mentalese.TermList, int, bool) {

	list := mentalese.TermList{}
	term := mentalese.Term{}
	ok := false
	newStartIndex := 0

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tOpeningBracket)
	if ok {
		term, newStartIndex, ok = parser.parseTerm(tokens, startIndex)
		if ok {
			list = append(list, term)
			startIndex = newStartIndex

			for startIndex < len(tokens) {
				_, newStartIndex, ok = parser.parseSingleToken(tokens, startIndex, tComma)
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
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tClosingBracket)
	}

end:

	return list, startIndex, ok
}

func (parser *InternalGrammarParser) parseRuleReference(tokens []Token, startIndex int) (mentalese.Term, int, bool) {

	term := mentalese.Term{}
	ok := false
	newStartIndex := 0
	predicate := ""
	prefix := ""
	alias := ""
	possibleAlias := ""
	applicationAlias := ""
	found := false

	_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tAmpersand)
	if !ok {
		goto end
	}

	possibleAlias, newStartIndex, ok = parser.parseSingleToken(tokens, startIndex, tPredicate)
	if ok {
		predicate, newStartIndex, ok = parser.parseSingleToken(tokens, newStartIndex, tColon)
		if ok {
			alias = possibleAlias
			startIndex = newStartIndex
		}
	}

	applicationAlias, found = parser.aliasMap[alias]
	if !found {
		ok = false
		goto end
	} else if applicationAlias != "" {
		prefix = applicationAlias + "_"
	}

	predicate, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tPredicate)

	term = mentalese.NewTermAtom(prefix + predicate)

end:

	return term, startIndex, ok
}

func (parser *InternalGrammarParser) parseSegmentationRulesAndCharacterClasses(tokens []Token, startIndex int) (*morphology.SegmentationRules, int, bool) {

	characterClasses := []morphology.CharacterClass{}
	segmentationRules := []morphology.SegmentationRule{}
	done := false

	for !done {

		characterClass, newStartIndex, ok := parser.parseCharacterClass(tokens, startIndex)
		if ok {
			startIndex = newStartIndex
			characterClasses = append(characterClasses, characterClass)
		} else {
			segmentationRule, newStartIndex, ok := parser.parseSegmentationRule(tokens, startIndex)
			if ok {
				startIndex = newStartIndex
				segmentationRules = append(segmentationRules, segmentationRule)
			} else {
				done = true
			}
		}
	}

	compiledRules := morphology.NewSegmentationRules()
	for _, rule := range segmentationRules {
		regexp, ok := morphology.BuildRegexp(rule.GetAntecedent().GetPattern(), characterClasses)
		if !ok {
			done = false
			break
		}
		compiledRules.Add(morphology.NewSegmentationRule(rule.GetAntecedent(), rule.GetConsequents(), regexp))
	}

	return compiledRules, startIndex, done
}

// consonant: ['b', 'c', 'd']
func (parser *InternalGrammarParser) parseCharacterClass(tokens []Token, startIndex int) (morphology.CharacterClass, int, bool) {

	characterClass := morphology.CharacterClass{}
	ok := true
	name := ""
	list := mentalese.TermList{}

	name, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tPredicate)
	if ok {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tColon)
		if ok {
			list, startIndex, ok = parser.parseTermList(tokens, startIndex)
			termType, _ := list.GetTermType()
			ok = termType == mentalese.TermTypeStringConstant
			if ok {
				characterClass = morphology.NewCharacterClass(name, list)
			}
		}
	}

	return characterClass, startIndex, ok
}

// comp: '*{consonant1}{consonant1}er' -> adj: '*{consonant1}', suffix: 'er'
func (parser *InternalGrammarParser) parseSegmentationRule(tokens []Token, startIndex int) (morphology.SegmentationRule, int, bool) {
	ok := false
	found := true
	newStartIndex := 0
	rule := morphology.SegmentationRule{}
	antecedent := morphology.SegmentNode{}
	consequent := morphology.SegmentNode{}
	consequents := []morphology.SegmentNode{}

	antecedent, startIndex, ok = parser.parseSegmentationNode(tokens, startIndex)
	if ok {
		_, newStartIndex, found = parser.parseSingleToken(tokens, startIndex, tRewrite)
		if found {
			startIndex = newStartIndex
			for true {
				consequent, newStartIndex, ok = parser.parseSegmentationNode(tokens, startIndex)
				if ok {
					startIndex = newStartIndex
					consequents = append(consequents, consequent)
				} else {
					break
				}
				_, newStartIndex, found = parser.parseSingleToken(tokens, startIndex, tComma)
				if found {
					startIndex = newStartIndex
				} else {
					break
				}
			}
		}
		if ok {
			rule = morphology.NewSegmentationRule(antecedent, consequents, nil)
		}
	}

	return rule, startIndex, ok
}

func (parser *InternalGrammarParser) parseSegmentationNode(tokens []Token, startIndex int) (morphology.SegmentNode, int, bool) {
	ok := false
	category := ""
	text := ""
	pattern := []morphology.SegmentPatternCharacter{}
	node := morphology.SegmentNode{}

	category, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tPredicate)
	if ok {
		_, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tColon)
		if ok {
			text, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tStringConstant)
			if ok {
				pattern, ok = parser.parseSegmentPattern(text)
				if ok {
					node = morphology.NewSegmentNode(category, pattern)
				}
			}
		}
	}

	return node, startIndex, ok
}

func (parser *InternalGrammarParser) parseSegmentPattern(text string) ([]morphology.SegmentPatternCharacter, bool) {
	pattern := []morphology.SegmentPatternCharacter{}
	ok := true

	expression, _ := regexp.Compile("(\\*|\\{([a-z]+)([0-9]+)\\}|[^\\*\\{]+)")

	for _, results := range expression.FindAllStringSubmatch(text, -1) {

		result := results[1]
		characterType := ""
		characterValue := ""
		index := -1

		if result == "*" {
			characterType = morphology.CharacterTypeRest
			characterValue = ""
		} else if result[0] == '{' {
			characterType = morphology.CharacterTypeClass
			characterValue = results[2]
			index, _ = strconv.Atoi(results[3])
		} else {
			characterType = morphology.CharacterTypeLiteral
			characterValue = result
		}

		pattern = append(pattern, morphology.NewSegmentPatterCharacter(characterType, characterValue, index))
	}

	return pattern, ok
}

func (parser *InternalGrammarParser) parseVariable(tokens []Token, startIndex int) (string, int, bool) {
	variable := ""
	prefix := ""
	_, newStartIndex, ok := parser.parseSingleToken(tokens, startIndex, tColon)
	if ok {
		startIndex = newStartIndex
		prefix = ":"
	}
	variable, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tVariable)
	return prefix + variable, startIndex, ok
}

func (parser *InternalGrammarParser) parseTerm(tokens []Token, startIndex int) (mentalese.Term, int, bool) {

	ok := false
	tokenValue := ""
	term := mentalese.Term{}
	newStartIndex := 0

	tokenValue, startIndex, ok = parser.parseVariable(tokens, startIndex)
	if ok {
		term.TermType = mentalese.TermTypeVariable
		term.TermValue = tokenValue
		goto end
	}
	tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tNumber)
	if ok {
		term.TermType = mentalese.TermTypeStringConstant
		term.TermValue = tokenValue
		goto end
	}
	tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tStringConstant)
	if ok {
		term.TermType = mentalese.TermTypeStringConstant
		term.TermValue = tokenValue
		goto end
	}
	tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tAnonymousVariable)
	if ok {
		term.TermType = mentalese.TermTypeAnonymousVariable
		term.TermValue = tokenValue
		goto end
	}
	{
		id := ""
		sort := ""
		id, sort, startIndex, ok = parser.parseId(tokens, startIndex)
		if ok {
			term.TermType = mentalese.TermTypeId
			term.TermValue = id
			term.TermSort = sort
			goto end
		}
	}
	{
		relationSet := mentalese.RelationSet{}
		relationSet, newStartIndex, ok = parser.parseRelations(tokens, startIndex, true)
		if ok {
			term.TermType = mentalese.TermTypeRelationSet
			term.TermValueRelationSet = relationSet
			startIndex = newStartIndex
			goto end
		}
	}
	{
		rule := mentalese.Rule{}
		rule, newStartIndex, ok = parser.parseRule(tokens, startIndex)
		if ok {
			term.TermType = mentalese.TermTypeRule
			term.TermValueRule = &rule
			startIndex = newStartIndex
			goto end
		}
	}
	{
		tokenValue, startIndex, ok = parser.parseSingleToken(tokens, startIndex, tPredicate)
		if ok {
			term.TermType = mentalese.TermTypePredicateAtom
			term.TermValue = tokenValue
			goto end
		}
	}
	{
		reference := mentalese.Term{}
		reference, newStartIndex, ok = parser.parseRuleReference(tokens, startIndex)
		if ok {
			term = reference
			startIndex = newStartIndex
			goto end
		}
	}
	{
		list := mentalese.TermList{}
		list, newStartIndex, ok = parser.parseTermList(tokens, startIndex)
		if ok {
			term.TermType = mentalese.TermTypeList
			term.TermValueList = list
			startIndex = newStartIndex
			goto end
		}
		{
			tokenValue, newStartIndex, ok = parser.parseSingleToken(tokens, startIndex, tRegExp)
			if ok {
				term.TermType = mentalese.TermTypeRegExp
				term.TermValue = tokenValue
				startIndex = newStartIndex
			}
		}
	}

end:

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
