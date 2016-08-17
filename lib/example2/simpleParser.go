package example2

import (
    "strings"
)

type simpleParser struct {
    grammar *simpleGrammar
    lexicon *simpleLexicon
}

func NewSimpleParser(grammar *simpleGrammar, lexicon *simpleLexicon) *simpleParser {
    return &simpleParser{grammar: grammar, lexicon: lexicon}
}

// Parses tokens using parser.grammar and parser.lexicon
func (parser *simpleParser) Process(tokens []string) (int, []SimpleRelation, bool) {

    length, _, relationList, ok := parser.parseAllRules("S", tokens, 0)
// TODO: remove parse tree nodes?
    return length, relationList, ok
}

// Parses tokens, starting from start, using all rules with given antedecent
func (parser *simpleParser) parseAllRules(antecedent string, tokens []string, start int) (int, SimpleParseTreeNode, []SimpleRelation, bool) {

    rules := parser.grammar.FindRules(antecedent)
    node := SimpleParseTreeNode{SyntacticCategory: antecedent}

    for i := 0; i < len(rules); i++ {

        rule := rules[i]
        cursor, childNodes, relations, ok := parser.parse(rule, tokens, start)

        if ok {
            node.Children = childNodes
            return cursor, node, relations, true
        }
    }

    return 0, node, []SimpleRelation{}, false
}

// Try to parse tokens using the rule given
func (parser *simpleParser) parse(rule SimpleGrammarRule, tokens []string, start int)  (int, []SimpleParseTreeNode, []SimpleRelation, bool) {

    cursor := start
    childNodes := []SimpleParseTreeNode{}
    syntacticCategories := rule.SyntacticCategories
    relations := []SimpleRelation{}

    relationTemplates := rule.RelationTemplates
    relations = append(relations, relationTemplates...)

    for i := 1; i < len(syntacticCategories); i++ {

        newCursor, childNode, childRelations, ok := parser.parseSingleConsequent(syntacticCategories[i], tokens, cursor)
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
func (parser *simpleParser) parseSingleConsequent(syntacticCategory string, tokens []string, start int)  (int, SimpleParseTreeNode, []SimpleRelation, bool) {

    node := SimpleParseTreeNode{SyntacticCategory:syntacticCategory}

    // if the sentence has run out of tokens, fail
    if start >= len(tokens) {
        return 0, node, []SimpleRelation{}, false
    }

    if strings.ToLower(syntacticCategory) == syntacticCategory {

        if parser.lexicon.CheckPartOfSpeech(tokens[start], syntacticCategory) {
            node.Word = tokens[start]

            relations := []SimpleRelation{}

            lexItem, found := parser.lexicon.GetLexItem(tokens[start], syntacticCategory)
            if found {
                relations = lexItem.RelationTemplates
            }

            return start + 1, node, relations, true

        } else {
            return 0, node, []SimpleRelation{}, false
        }

    } else {

        newCursor, node, relations, ok := parser.parseAllRules(syntacticCategory, tokens, start)
        return newCursor, node, relations, ok

    }
}
