package example2

import (
    "strings"
    "fmt"
)

type simpleParser struct {
    grammar *simpleGrammar
    lexicon *simpleLexicon
    varIndexCounter map[string]int
}

func NewSimpleParser(grammar *simpleGrammar, lexicon *simpleLexicon) *simpleParser {
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
func (parser *simpleParser) createVariableMap(variable string, entityVariables []string) map[string]string {

    m := map[string]string{}

    for i := 1; i < len(entityVariables); i++ {

        entityVariable := entityVariables[i]

        if entityVariable == entityVariables[0] {

            // the consequent variable matches the antecedent variable, inherit its actual variable
            m[entityVariable] = variable

        } else {

            // we're going to add a new actual variable, unless we already have
            _, present := m[entityVariable]
            if !present {
                m[entityVariable] = parser.getNewVariable(entityVariable)
            }
        }
    }

    return m
}

// Parses tokens using parser.grammar and parser.lexicon
func (parser *simpleParser) Process(tokens []string) (int, []SimpleRelation, bool) {

    length, _, relationList, ok := parser.parseAllRules("S", tokens, 0, parser.getNewVariable("sentence"))
// TODO: remove parse tree nodes?
    return length, relationList, ok
}

// Parses tokens, starting from start, using all rules with given antedecent
func (parser *simpleParser) parseAllRules(antecedent string, tokens []string, start int, antecedentVariable string) (int, SimpleParseTreeNode, []SimpleRelation, bool) {

    rules := parser.grammar.FindRules(antecedent)
    node := SimpleParseTreeNode{SyntacticCategory: antecedent}

    for i := 0; i < len(rules); i++ {

        rule := rules[i]

        cursor, childNodes, relations, ok := parser.parse(rule, tokens, start, antecedentVariable)

        if ok {
            node.Children = childNodes
            return cursor, node, relations, true
        }
    }

    return 0, node, []SimpleRelation{}, false
}

// Try to parse tokens using the rule given
func (parser *simpleParser) parse(rule SimpleGrammarRule, tokens []string, start int, antecedentVariable string)  (int, []SimpleParseTreeNode, []SimpleRelation, bool) {

    cursor := start
    childNodes := []SimpleParseTreeNode{}
    syntacticCategories := rule.SyntacticCategories
    relations := []SimpleRelation{}

    // create a map of formal variables to actual variables (new variables are created)
    variableMap := parser.createVariableMap(antecedentVariable, rule.EntityVariables)

    // non-leaf node relations
    ruleRelations := parser.createGrammarRuleRelations(rule.RelationTemplates, variableMap)
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
            return 0, childNodes, []SimpleRelation{}, false;
        }
    }

    return cursor, childNodes, relations, true
}

// Try to parse tokens given a single syntactic category
// Returns the index to the token following the parsed sequence
func (parser *simpleParser) parseSingleConsequent(syntacticCategory string, tokens []string, start int, v string)  (int, SimpleParseTreeNode, []SimpleRelation, bool) {

    node := SimpleParseTreeNode{SyntacticCategory:syntacticCategory}

    // if the sentence has run out of tokens, fail
    if start >= len(tokens) {
        return 0, node, []SimpleRelation{}, false
    }

    if strings.ToLower(syntacticCategory) == syntacticCategory {

        // leaf node
        if parser.lexicon.CheckPartOfSpeech(tokens[start], syntacticCategory) {
            node.Word = tokens[start]

            relations := []SimpleRelation{}

            lexItem, found := parser.lexicon.GetLexItem(tokens[start], syntacticCategory)
            if found {
                // leaf node relations
                relations = parser.createLexItemRelations(lexItem.RelationTemplates, v)
            }

            return start + 1, node, relations, true

        } else {
            return 0, node, []SimpleRelation{}, false
        }

    } else {

        // non leaf-node
        newCursor, node, relations, ok := parser.parseAllRules(syntacticCategory, tokens, start, v)
        return newCursor, node, relations, ok

    }
}

// Create actual relations given a set of templates and a variable map (formal to actual variables)
func (parser *simpleParser) createGrammarRuleRelations(relationTemplates []SimpleRelation, variableMap map[string]string) []SimpleRelation {

    relations := []SimpleRelation{}

    for i := 0; i < len(relationTemplates); i++ {
        relation := relationTemplates[i]

        for a := 0; a < len(relation.Arguments); a++ {
            argument := relation.Arguments[a]
            relation.Arguments[a] = variableMap[argument]
        }

        relations = append(relations, relation)
    }

    return relations
}

// Create actual relations given a set of templates and an actual variable to replace any * positions
func (parser *simpleParser) createLexItemRelations(relationTemplates []SimpleRelation, variable string) []SimpleRelation {

    relations := []SimpleRelation{}

    for i := 0; i < len(relationTemplates); i++ {
        relation := relationTemplates[i]

        for a := 0; a < len(relation.Arguments); a++ {
            argument := relation.Arguments[a]
            if argument == "*" {
                relation.Arguments[a] = variable
            }
        }

        relations = append(relations, relation)
    }

    return relations
}
