package parse

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
)

// A simple top-down parser
// Note: does not support left-recursive rewrite rules; infinite looping

type TopDownParser struct {
	grammar         *Grammar
	lexicon         *Lexicon
	senseBuilder    SenseBuilder
}

func NewTopDownParser(grammar *Grammar, lexicon *Lexicon) *TopDownParser {
	return &TopDownParser{grammar: grammar, lexicon: lexicon, senseBuilder: NewSenseBuilder()}
}

// Parses tokens using parser.grammar and parser.lexicon
func (parser *TopDownParser) Parse(tokens []string) (mentalese.RelationSet, int, bool) {

	length, relationList, ok := parser.parseAllRules("s", tokens, 0, parser.senseBuilder.GetNewVariable("Sentence"))

	return relationList, length, ok
}

// Parses tokens, starting from start, using all rules with given antecedent
func (parser *TopDownParser) parseAllRules(antecedent string, tokens []string, start int, antecedentVariable string) (int, mentalese.RelationSet, bool) {

	common.LogTree("parseAllRules", antecedent)

	rules := parser.grammar.FindRules(antecedent)
	cursor := 0
	ok := false
	relations := mentalese.RelationSet{}

	for _, rule := range rules {
		cursor, relations, ok = parser.parseWithRule(rule, tokens, start, antecedentVariable)

		if ok {
			break
		}
	}

	common.LogTree("parseAllRules", cursor, relations, ok)

	return cursor, relations, ok
}

// Try to parse tokens using the rule given
func (parser *TopDownParser) parseWithRule(rule GrammarRule, tokens []string, start int, antecedentVariable string) (int, mentalese.RelationSet, bool) {

	cursor := start
	syntacticCategories := rule.SyntacticCategories
	relations := mentalese.RelationSet{}
	success := true

	common.LogTree("parse", rule)

	// create a map of formal variables to actual variables (new variables are created)
	variableMap := parser.senseBuilder.CreateVariableMap(antecedentVariable, rule.EntityVariables)

	// non-leaf node relations
	ruleRelations := parser.senseBuilder.CreateGrammarRuleRelations(rule.Sense, variableMap)
	relations = append(relations, ruleRelations...)

	// parse each of the children
	for i := 1; i < len(syntacticCategories); i++ {

		consequentVariable := variableMap[rule.EntityVariables[i]]
		newCursor, childRelations, ok := parser.parseSingleConsequent(syntacticCategories[i], tokens, cursor, consequentVariable)
		if ok {
			relations = append(relations, childRelations...)
			cursor = newCursor
		} else {
			cursor = 0
			relations = mentalese.RelationSet{}
			success = false
			break
		}
	}

	common.LogTree("parse", cursor, relations, success)

	return cursor, relations, success
}

// Try to parse tokens given a single syntactic category
// Returns the index to the token following the parsed sequence
func (parser *TopDownParser) parseSingleConsequent(syntacticCategory string, tokens []string, start int, v string) (int, mentalese.RelationSet, bool) {

	cursor := 0
	relations := mentalese.RelationSet{}
	ok := false

	common.LogTree("parseSingleConsequent", syntacticCategory, tokens, start, v)

	// if the sentence has run out of tokens, fail
	if start < len(tokens) {

		token := tokens[start]

		// leaf node?
		lexItem, found := parser.lexicon.GetLexItem(token, syntacticCategory)
		if found {

			// leaf node relations
			relations = parser.senseBuilder.CreateLexItemRelations(lexItem.RelationTemplates, v)
			cursor = start + 1
			ok = true

		} else {

			// non leaf-node
			cursor, relations, ok = parser.parseAllRules(syntacticCategory, tokens, start, v)

		}
	}

	common.LogTree("parseSingleConsequent", cursor, relations, ok)

	return cursor, relations, ok
}
