package mentalese

type GrammarRules struct {
	index map[string]map[int][]GrammarRule
}

func NewGrammarRules() *GrammarRules {
	return &GrammarRules{
		index: map[string]map[int][]GrammarRule{},
	}
}

func (rules *GrammarRules) AddRule(rule GrammarRule) {

	antecedent := rule.GetAntecedent()
	argumentCount := len(rule.GetAntecedentVariables())

	_, found := rules.index[antecedent]
	if !found {
		rules.index[antecedent] = map[int][]GrammarRule{}
	}
	_, found = rules.index[antecedent][argumentCount]
	if !found {
		rules.index[antecedent][argumentCount] = []GrammarRule{}
	}

	rules.index[antecedent][argumentCount] = append(rules.index[antecedent][argumentCount], rule)
}

func (rules *GrammarRules) Merge(otherRules *GrammarRules) {
	for _, m := range otherRules.index {
		for _, r := range m {
			for _, rule := range r {
				rules.AddRule(rule)
			}
		}
	}
}

// returns index, ok (where index is an array of string-arrays)
func (grammar *GrammarRules) FindRules(antecedent string, argumentCount int) []GrammarRule {
	rules, ok := grammar.index[antecedent][argumentCount]

	if ok {
		return rules
	} else {
		return []GrammarRule{}
	}
}

func (grammar *GrammarRules) ImportFrom(fromGrammar *GrammarRules) {
	for _, rulesPerCategory := range fromGrammar.index {
		for _, rulesPerCount := range rulesPerCategory {
			for _, rule := range rulesPerCount {
				grammar.AddRule(rule)
			}
		}
	}
}

// Is word present in a rule `category(E) -> "word"`?
func (grammar *GrammarRules) WordOccurs(word string, category string) bool {
	found := false
	rules := grammar.FindRules(category, 1)
	for _, rule := range rules {
		if len(rule.SyntacticCategories) == 2 {
			if rule.PositionTypes[1] == PosTypeWordForm {
				if rule.SyntacticCategories[1] == word {
					return true
				}
			}
		}
	}
	return found
}
