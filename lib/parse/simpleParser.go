package parse

import (
	"fmt"
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
)

type simpleParser struct {
	grammar         *SimpleGrammar
	lexicon         *SimpleLexicon
	varIndexCounter map[string]int
}

func NewSimpleParser(grammar *SimpleGrammar, lexicon *SimpleLexicon) *simpleParser {
	return &simpleParser{grammar: grammar, lexicon: lexicon, varIndexCounter: map[string]int{}}
}

// Returns a new variable name
func (parser *simpleParser) getNewVariable(formalVariable string) string {

	initial := formalVariable[0:1]

	_, present := parser.varIndexCounter[initial]
	if !present {
		parser.varIndexCounter[initial] = 1
	} else {
		parser.varIndexCounter[initial]++
	}

	return fmt.Sprint(initial, parser.varIndexCounter[initial])
}

// Creates a map of formal variables to actual variables (new variables are created)
func (parser *simpleParser) createVariableMap(actualAntecedent string, formalVariables []string) map[string]string {

	m := map[string]string{}
	antecedentVariable := formalVariables[0]

	for i := 1; i < len(formalVariables); i++ {

		consequentVariable := formalVariables[i]

		if consequentVariable == antecedentVariable {

			// the consequent variable matches the antecedent variable, inherit its actual variable
			m[consequentVariable] = actualAntecedent

		} else {

			// we're going to add a new actual variable, unless we already have
			_, present := m[consequentVariable]
			if !present {
				m[consequentVariable] = parser.getNewVariable(consequentVariable)
			}
		}
	}

	return m
}

// Parses tokens using parser.grammar and parser.lexicon
func (parser *simpleParser) Process(tokens []string) (int, mentalese.SimpleRelationSet, bool) {

	length, _, relationList, ok := parser.parseAllRules("s", tokens, 0, parser.getNewVariable("Sentence"))

	set := mentalese.SimpleRelationSet{}
	set = append(set, relationList...)
	// TODO: remove parse tree nodes?
	return length, set, ok
}

// Parses tokens, starting from start, using all rules with given antecedent
func (parser *simpleParser) parseAllRules(antecedent string, tokens []string, start int, antecedentVariable string) (int, SimpleParseTreeNode, []mentalese.SimpleRelation, bool) {

	common.Logf("parseAllRules: %s\n", antecedent)

	rules := parser.grammar.FindRules(antecedent)
	node := SimpleParseTreeNode{SyntacticCategory: antecedent}

	for _, rule := range rules {

		cursor, childNodes, relations, ok := parser.parse(rule, tokens, start, antecedentVariable)

		if ok {
			node.Children = childNodes

			common.Log("parseAllRules end 1\n")

			return cursor, node, relations, true
		}
	}

	common.Log("parseAllRules end 2\n")

	return 0, node, []mentalese.SimpleRelation{}, false
}

// Try to parse tokens using the rule given
func (parser *simpleParser) parse(rule SimpleGrammarRule, tokens []string, start int, antecedentVariable string) (int, []SimpleParseTreeNode, []mentalese.SimpleRelation, bool) {

	cursor := start
	childNodes := []SimpleParseTreeNode{}
	syntacticCategories := rule.SyntacticCategories
	relations := []mentalese.SimpleRelation{}

	common.Logf("parse %v\n", rule)

	// create a map of formal variables to actual variables (new variables are created)
	variableMap := parser.createVariableMap(antecedentVariable, rule.EntityVariables)

	// non-leaf node relations
	ruleRelations := parser.createGrammarRuleRelations(rule.Sense, variableMap)
	relations = append(relations, ruleRelations...)

	// parse each of the children
	for i := 1; i < len(syntacticCategories); i++ {

		consequentVariable := variableMap[rule.EntityVariables[i]]
		newCursor, childNode, childRelations, ok := parser.parseSingleConsequent(syntacticCategories[i], tokens, cursor, consequentVariable)
		if ok {
			childNodes = append(childNodes, childNode)
			relations = append(relations, childRelations...)
			cursor = newCursor
		} else {
			return 0, childNodes, []mentalese.SimpleRelation{}, false
		}
	}

	return cursor, childNodes, relations, true
}

// Try to parse tokens given a single syntactic category
// Returns the index to the token following the parsed sequence
func (parser *simpleParser) parseSingleConsequent(syntacticCategory string, tokens []string, start int, v string) (int, SimpleParseTreeNode, []mentalese.SimpleRelation, bool) {

	node := SimpleParseTreeNode{SyntacticCategory: syntacticCategory}

	// if the sentence has run out of tokens, fail
	if start >= len(tokens) {
		return 0, node, []mentalese.SimpleRelation{}, false
	}

	token := tokens[start]

	common.Logf("parseSingleConsequent: %s (%s)\n", token, syntacticCategory)

	// leaf node?
	lexItem, found := parser.lexicon.GetLexItem(token, syntacticCategory)
	if found {
		node.Word = tokens[start]

		common.Logf("Leaf node: %s\n", token)

		relations := []mentalese.SimpleRelation{}

		// leaf node relations
		relations = parser.createLexItemRelations(lexItem.RelationTemplates, v)

		return start + 1, node, relations, true

	} else {

		// non leaf-node
		newCursor, node, relations, ok := parser.parseAllRules(syntacticCategory, tokens, start, v)
		return newCursor, node, relations, ok

	}
}

// Create actual relations given a set of templates and a variable map (formal to actual variables)
func (parser *simpleParser) createGrammarRuleRelations(relationTemplates []mentalese.SimpleRelation, variableMap map[string]string) []mentalese.SimpleRelation {

	relations := []mentalese.SimpleRelation{}

	for _, relation := range relationTemplates {
		for a, argument := range relation.Arguments {

			relation.Arguments[a].TermType = mentalese.Term_variable
			relation.Arguments[a].TermValue = variableMap[argument.TermValue]
		}

		relations = append(relations, relation)
	}

	return relations
}

// Create actual relations given a set of templates and an actual variable to replace any * positions
func (parser *simpleParser) createLexItemRelations(relationTemplates []mentalese.SimpleRelation, variable string) []mentalese.SimpleRelation {

	relations := []mentalese.SimpleRelation{}

	for _, relation := range relationTemplates {
		for a, argument := range relation.Arguments {
			if argument.TermType == mentalese.Term_predicateAtom && argument.TermValue == "this" {

				relation.Arguments[a].TermType = mentalese.Term_variable
				relation.Arguments[a].TermValue = variable
			}
		}

		relations = append(relations, relation)
	}

	return relations
}
