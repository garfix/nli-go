package parse

import (
	"fmt"
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
)

type Parser struct {
	grammar         *Grammar
	lexicon         *Lexicon
	varIndexCounter map[string]int
}

func NewParser(grammar *Grammar, lexicon *Lexicon) *Parser {
	return &Parser{grammar: grammar, lexicon: lexicon, varIndexCounter: map[string]int{}}
}

// Returns a new variable name
func (parser *Parser) getNewVariable(formalVariable string) string {

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
func (parser *Parser) createVariableMap(actualAntecedent string, formalVariables []string) map[string]string {

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
func (parser *Parser) Process(tokens []string) (mentalese.RelationSet, int, bool) {

	length, relationList, ok := parser.parseAllRules("s", tokens, 0, parser.getNewVariable("Sentence"))

	set := mentalese.RelationSet{}
	set = append(set, relationList...)

	return set, length, ok
}

// Parses tokens, starting from start, using all rules with given antecedent
func (parser *Parser) parseAllRules(antecedent string, tokens []string, start int, antecedentVariable string) (int, []mentalese.Relation, bool) {

	common.LogTree("parseAllRules", antecedent)

	rules := parser.grammar.FindRules(antecedent)
	cursor := 0
	ok := false
	relations := []mentalese.Relation{}

	for _, rule := range rules {
		cursor, relations, ok = parser.parse(rule, tokens, start, antecedentVariable)

		if ok {
			break
		}
	}

	common.LogTree("parseAllRules", cursor, relations, ok)

	return cursor, relations, ok
}

// Try to parse tokens using the rule given
func (parser *Parser) parse(rule GrammarRule, tokens []string, start int, antecedentVariable string) (int, []mentalese.Relation, bool) {

	cursor := start
	syntacticCategories := rule.SyntacticCategories
	relations := []mentalese.Relation{}
	success := true

	common.LogTree("parse", rule)

	// create a map of formal variables to actual variables (new variables are created)
	variableMap := parser.createVariableMap(antecedentVariable, rule.EntityVariables)

	// non-leaf node relations
	ruleRelations := parser.createGrammarRuleRelations(rule.Sense, variableMap)
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
			relations = []mentalese.Relation{}
			success = false
			break
		}
	}

	common.LogTree("parse", cursor, relations, success)

	return cursor, relations, success
}

// Try to parse tokens given a single syntactic category
// Returns the index to the token following the parsed sequence
func (parser *Parser) parseSingleConsequent(syntacticCategory string, tokens []string, start int, v string) (int, []mentalese.Relation, bool) {

	cursor := 0
	relations := []mentalese.Relation{}
	ok := false

	common.LogTree("parseSingleConsequent", syntacticCategory, tokens, start, v)

	// if the sentence has run out of tokens, fail
	if start < len(tokens) {

		token := tokens[start]

		// leaf node?
		lexItem, found := parser.lexicon.GetLexItem(token, syntacticCategory)
		if found {

			// leaf node relations
			relations = parser.createLexItemRelations(lexItem.RelationTemplates, v)
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

// Create actual relations given a set of templates and a variable map (formal to actual variables)
func (parser *Parser) createGrammarRuleRelations(relationTemplates []mentalese.Relation, variableMap map[string]string) []mentalese.Relation {

	relations := []mentalese.Relation{}

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
func (parser *Parser) createLexItemRelations(relationTemplates []mentalese.Relation, variable string) []mentalese.Relation {

	relations := []mentalese.Relation{}

	for _, relationTemplate := range relationTemplates {

		relation := mentalese.Relation{}
		relation.Predicate = relationTemplate.Predicate

		for _, argument := range relationTemplate.Arguments {

			relationArgument := argument

			if argument.TermType == mentalese.Term_predicateAtom && argument.TermValue == "this" {

				relationArgument.TermType = mentalese.Term_variable
				relationArgument.TermValue = variable
			}

			relation.Arguments = append(relation.Arguments, relationArgument)
		}

		relations = append(relations, relation)
	}

	return relations
}
