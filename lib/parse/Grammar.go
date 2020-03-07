package parse

type Grammar struct {
	rules map[string]map[int][]GrammarRule
}

func NewGrammar() *Grammar {
	return &Grammar{rules: map[string]map[int][]GrammarRule{}}
}

func (grammar *Grammar) AddRule(rule GrammarRule) {

	antecedent := rule.GetAntecedent()
	argumentCount := len(rule.GetAntecedentVariables())

	_, found := grammar.rules[antecedent]
	if !found {
		grammar.rules[antecedent] = map[int][]GrammarRule{}
	}
	_, found = grammar.rules[antecedent][argumentCount]
	if !found {
		grammar.rules[antecedent][argumentCount] = []GrammarRule{}
	}

	grammar.rules[antecedent][argumentCount] = append(grammar.rules[antecedent][argumentCount], rule)
}

// returns rules, ok (where rules is an array of string-arrays)
func (grammar *Grammar) FindRules(antecedent string, argumentCount int) []GrammarRule {
	rules, ok := grammar.rules[antecedent][argumentCount]

	if ok {
		return rules
	} else {
		return []GrammarRule{}
	}
}

func (grammar *Grammar) ImportFrom(fromGrammar *Grammar) {
	for _, rulesPerCategory := range fromGrammar.rules {
		for _, rulesPerCount := range rulesPerCategory {
			for _, rule := range rulesPerCount {
				grammar.AddRule(rule)
			}
		}
	}
}
