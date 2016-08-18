package example1

type simpleGrammar struct {
    rules map[string][][]string
}

func NewSimpleGrammar(rules map[string][][]string) *simpleGrammar {
    return &simpleGrammar{rules: rules}
}

// returns rules, ok (where rules is an array of string-arrays)
func (grammar *simpleGrammar) FindRules(antecedent string) [][]string {
    rules, ok := grammar.rules[antecedent];

    if ok {
        return rules
    } else {
        return [][]string{}
    }
}
