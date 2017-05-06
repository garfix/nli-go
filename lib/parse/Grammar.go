package parse

type Grammar struct {
	rules map[string][]GrammarRule
}

func NewGrammar() *Grammar {
	return &Grammar{rules: map[string][]GrammarRule{}}
}

func (grammar *Grammar) AddRule(rule GrammarRule) {

	antecedent := rule.SyntacticCategories[0]

	grammar.rules[antecedent] = append(grammar.rules[antecedent], rule)
}

// returns rules, ok (where rules is an array of string-arrays)
func (grammar *Grammar) FindRules(antecedent string) []GrammarRule {
	rules, ok := grammar.rules[antecedent]

	if ok {
		return rules
	} else {
		return []GrammarRule{}
	}
}

func (grammar *Grammar) ImportFrom(fromGrammar *Grammar) {
	for _, rules := range fromGrammar.rules {
		for _, rule := range rules {
			grammar.AddRule(rule)
		}
	}
}
