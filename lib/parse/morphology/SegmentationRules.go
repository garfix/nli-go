package morphology


type SegmentationRules struct {
	index map[string][]SegmentationRule
}

func NewSegmentationRules() *SegmentationRules {
	return &SegmentationRules{
		index: map[string][]SegmentationRule{},
	}
}

func (rules *SegmentationRules) Add(rule SegmentationRule) {

	antecedent := rule.GetAntecedent()
	category := antecedent.GetCategory()

	_, found := rules.index[category]
	if !found {
		rules.index[category] = []SegmentationRule{}
	}

	rules.index[category] = append(rules.index[category], rule)
}

func (rules *SegmentationRules) FindRules(word string, category string) ([]SegmentationRule, []map[string]string) {

	foundRules := []SegmentationRule{}
	bindings := []map[string]string{}

	segmentationRules, found := rules.index[category]
	if !found {
		return foundRules, bindings
	}

	for _, aRule := range segmentationRules {
		if aRule.antecedent.category != category {
			continue
		}

		someResults, match := aRule.Matches(word)
		if match {
			binding, ok := rules.findBinding(someResults, aRule.antecedent.GetPattern())
			if !ok {
				continue
			}
			bindings = append(bindings, binding)
			foundRules = append(foundRules, aRule)
		}
	}

	return foundRules, bindings
}

func (rules *SegmentationRules) findBinding(results []string, pattern []SegmentPatternCharacter) (map[string]string, bool) {

	binding := map[string]string{}
	ok := true

	for i, character := range pattern {
		if character.characterType == CharacterTypeClass {
			variable := character.GetVariable()
			value := results[i + 1]

			// check if the variable exists and has the same value
			existing, found := binding[variable]
			if found && existing != value {
				ok = false
				break
			}

			binding[variable] = value
		} else if character.characterType == CharacterTypeRest {
			value := results[i + 1]
			binding[CharacterTypeRest] = value
		}
	}

	return binding, ok
}