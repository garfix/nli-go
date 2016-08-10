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
func (parser *simpleParser) Process(tokens []string) bool {

    _, ok := parser.parseAllRules("S", tokens, 0)
    return ok
}

// Parses tokens, starting from start, using all rules with given antedecent
func (parser *simpleParser) parseAllRules(antecedent string, tokens []string, start int) (int, bool) {

    rules := parser.grammar.FindRules(antecedent)

    for i := 0; i < len(rules); i++ {

        consequents := rules[i]
        cursor, ok := parser.parse(consequents, tokens, start)

        if ok {
            return cursor, true
        }
    }

    return 0, false
}

// Try to parse tokens using the rule given in consequents
// Return true/false for success
func (parser *simpleParser) parse(consequents []string, tokens []string, start int)  (int, bool) {

    cursor := start

    for i := 0; i < len(consequents); i++ {

        newCursor, ok := parser.parseSingleConsequent(consequents[i], tokens, cursor)
        if ok {
            cursor = newCursor
        } else {
            return 0, false;
        }
    }

    return cursor, true
}

// Try to parse tokens given a single syntactic category
// Returns the index to the token following the parsed sequence
func (parser *simpleParser) parseSingleConsequent(syntacticCategory string, tokens []string, start int)  (int, bool) {

    // if the sentence has run out of tokens, fail
    if start >= len(tokens) {
        return 0, false
    }

    if strings.ToLower(syntacticCategory) == syntacticCategory {

        if parser.lexicon.CheckPartOfSpeech(tokens[start], syntacticCategory) {
            return start + 1, true
        } else {
            return 0, false
        }

    } else {

        newCursor, ok := parser.parseAllRules(syntacticCategory, tokens, start)
        return newCursor, ok

    }
}