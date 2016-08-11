package lib

import (
    "nli-go/types"
    "strings"
)

type simpleParser struct {
    grammar types.Grammar
    lexicon types.Lexicon
}

func NewSimpleParser(grammar types.Grammar, lexicon types.Lexicon) *simpleParser {
    return &simpleParser{grammar: grammar, lexicon: lexicon}
}

// Parses tokens using parser.grammar and parser.lexicon
func (parser *simpleParser) Process(tokens []string) (int, SimpleParseTreeNode, bool) {

    length, parseTree, ok := parser.parseAllRules("S", tokens, 0)
    return length, parseTree, ok
}

// Parses tokens, starting from start, using all rules with given antedecent
func (parser *simpleParser) parseAllRules(antecedent string, tokens []string, start int) (int, SimpleParseTreeNode, bool) {

    rules := parser.grammar.FindRules(antecedent)
    node := SimpleParseTreeNode{SyntacticCategory: antecedent}

    for i := 0; i < len(rules); i++ {

        consequents := rules[i]
        cursor, childNodes, ok := parser.parse(consequents, tokens, start)

        if ok {
            node.Children = childNodes
            return cursor, node, true
        }
    }

    return 0, node, false
}

// Try to parse tokens using the rule given in consequents
// Return true/false for success
func (parser *simpleParser) parse(consequents []string, tokens []string, start int)  (int, []SimpleParseTreeNode, bool) {

    cursor := start
    childNodes := []SimpleParseTreeNode{}

    for i := 0; i < len(consequents); i++ {

        newCursor, childNode, ok := parser.parseSingleConsequent(consequents[i], tokens, cursor)
        if ok {
            childNodes = append(childNodes, childNode)
            cursor = newCursor
        } else {
            return 0, childNodes, false;
        }
    }

    return cursor, childNodes, true
}

// Try to parse tokens given a single syntactic category
// Returns the index to the token following the parsed sequence
func (parser *simpleParser) parseSingleConsequent(syntacticCategory string, tokens []string, start int)  (int, SimpleParseTreeNode, bool) {

    node := SimpleParseTreeNode{SyntacticCategory:syntacticCategory}

    // if the sentence has run out of tokens, fail
    if start >= len(tokens) {
        return 0, node, false
    }

    if strings.ToLower(syntacticCategory) == syntacticCategory {

        if parser.lexicon.CheckPartOfSpeech(tokens[start], syntacticCategory) {
            node.Word = tokens[start]
            return start + 1, node, true
        } else {
            return 0, node, false
        }

    } else {

        newCursor, node, ok := parser.parseAllRules(syntacticCategory, tokens, start)
        return newCursor, node, ok

    }
}
