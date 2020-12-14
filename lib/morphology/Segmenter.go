package morphology

type Segmenter struct {

}

func NewSegmenter() *Segmenter {
	return &Segmenter{}
}

func (segmenter *Segmenter) Segment(segmentationRules []SegmentationRule, word string, category string) []string {

	segments := []string{}

	rules, bindings := segmenter.findRules(segmentationRules, word, category)

	for i, rule := range rules {

		singleRuleSegments := []string{}
		binding := bindings[i]

		if rule.IsTerminal() {
			segments = []string{ word }
			break
		}

		ok := true
		for _, consequent := range rule.GetConsequents() {
			consequentSegments := segmenter.buildSegments(consequent, binding, segmentationRules)
			if len(consequentSegments) == 0 {
				ok = false
				break
			}
			singleRuleSegments = append(singleRuleSegments, consequentSegments...)
		}

		if ok {
			segments = singleRuleSegments
			break
		}
	}

	return segments
}

func (segmenter *Segmenter) buildSegments(segmentNode SegmentNode, binding map[string]string, segmentationRules []SegmentationRule) []string {

	word := segmenter.buildWord(segmentNode.GetPattern(), binding)

	return segmenter.Segment(segmentationRules, word, segmentNode.category)
}

func (segmenter *Segmenter) buildWord(pattern []SegmentPatternCharacter, binding map[string]string) string {

	word := ""

	for _, character := range pattern {
		if character.characterType == CharacterTypeRest {
			word += binding[CharacterTypeRest]
		} else if character.characterType == CharacterTypeClass {
			word += binding[character.GetVariable()]
		} else {
			word += character.characterValue
		}
	}

	return word
}

func (segmenter *Segmenter) findRules(segmentationRules []SegmentationRule, word string, category string) ([]SegmentationRule, []map[string]string) {

	rules := []SegmentationRule{}
	bindings := []map[string]string{}

	for _, aRule := range segmentationRules {
		if aRule.antecedent.category != category {
			continue
		}

		someresults, match := aRule.Matches(word)
		if match {
			binding, ok := segmenter.findBinding(someresults, aRule.antecedent.GetPattern())
			if !ok {
				continue
			}
			bindings = append(bindings, binding)
			rules = append(rules, aRule)
		}
	}

	return rules, bindings
}

func (segmenter *Segmenter) findBinding(results []string, pattern []SegmentPatternCharacter) (map[string]string, bool) {

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