package example2

type simpleGrammar struct {
	rules map[string][]SimpleGrammarRule
}

func NewSimpleGrammar(rules map[string][]SimpleGrammarRule) *simpleGrammar {
	return &simpleGrammar{rules: rules}
}

// returns rules, ok (where rules is an array of string-arrays)
func (grammar *simpleGrammar) FindRules(antecedent string) []SimpleGrammarRule {
	rules, ok := grammar.rules[antecedent]

	if ok {
		return rules
	} else {
		return []SimpleGrammarRule{}
	}
}
