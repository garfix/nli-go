package example2

type simpleGrammar struct {
    rules map[string][]SimpleGrammarRule
}

func NewSimpleGrammar(rules map[string][]SimpleGrammarRule) *simpleGrammar {
    return &simpleGrammar{rules: rules}
}

// returns rules, ok (where rules is an array of string-arrays)
func (grammar *simpleGrammar) FindRules(antecedent string) []SimpleGrammarRule {
    rules, ok := grammar.rules[antecedent];

    if ok {
        return rules
    } else {
        return []SimpleGrammarRule{}
    }
}

// returns rules, ok (where rules is an array of string-arrays)
func (grammar *simpleGrammar) FindRule(syntacticCategories []string) (SimpleGrammarRule, bool) {

    rules, ok := grammar.rules[syntacticCategories[0]];
    if ok {

        for i := 0; i < len(rules); i++ {

            rule := rules[0]
            ruleCategories := rule.SyntacticCategories

            found := true
            for c := 1; c < len(ruleCategories); c++ {
                if ruleCategories[c] != syntacticCategories[c] {
                    found = false;
                    break;
                }
            }

            if found {
                return rule, true
            }
        }

    }

    return SimpleGrammarRule{}, false
}
