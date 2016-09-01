package example2

type SimpleGrammar struct {
	rules map[string][]SimpleGrammarRule
}

func NewSimpleGrammar(rules map[string][]SimpleGrammarRule) *SimpleGrammar {
	return &SimpleGrammar{rules: rules}
}

// returns rules, ok (where rules is an array of string-arrays)
func (grammar *SimpleGrammar) FindRules(antecedent string) []SimpleGrammarRule {
	rules, ok := grammar.rules[antecedent]

	if ok {
		return rules
	} else {
		return []SimpleGrammarRule{}
	}
}
