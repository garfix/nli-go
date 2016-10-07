package parse

type SimpleGrammar struct {
	rules map[string][]SimpleGrammarRule
}

func NewSimpleGrammar() *SimpleGrammar {
	return &SimpleGrammar{rules: map[string][]SimpleGrammarRule{}}
}

func (grammar *SimpleGrammar) AddRule(rule SimpleGrammarRule) {

	antecedent := rule.SyntacticCategories[0]

	grammar.rules[antecedent] = append(grammar.rules[antecedent], rule)
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
