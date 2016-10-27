package generate

import "nli-go/lib/mentalese"

type GenerationGrammar struct {
	rules map[string][]GenerationGrammarRule
}

func NewGenerationGrammar() *GenerationGrammar {
	return &GenerationGrammar{rules: map[string][]GenerationGrammarRule{}}
}

func (grammar *GenerationGrammar) AddRule(rule GenerationGrammarRule) {

	antecedent := rule.Antecedent.Predicate

	grammar.rules[antecedent] = append(grammar.rules[antecedent], rule)
}

// returns rules, ok (where rules is an array of string-arrays)
func (grammar *GenerationGrammar) FindRules(antecedent mentalese.Relation) []GenerationGrammarRule {
	rules, ok := grammar.rules[antecedent.Predicate]

	if ok {
		return rules
	} else {
		return []GenerationGrammarRule{}
	}
}
