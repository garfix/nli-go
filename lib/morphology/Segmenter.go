package morphology

type Segmenter struct {

}

func NewSegmenter() *Segmenter {
	return &Segmenter{}
}

func (segmenter *Segmenter) Segment(segmentationRules []SegmentationRule, word string, category string) []string {

	segments := []string{}

	rule, found, binding := segmenter.findRule(segmentationRules, word, category)

	if found {
		for _, consequent := range rule.GetConsequents() {
			segments = append(segments, segmenter.buildSegments(consequent, binding, segmentationRules)...)
		}
	} else {
		segments = []string{ word }
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

func (segmenter *Segmenter) findRule(segmentationRules []SegmentationRule, word string, category string) (SegmentationRule, bool, map[string]string) {

	rule := SegmentationRule{}
	found := false
	ok := false
	binding := map[string]string{}

	for _, aRule := range segmentationRules {
		if aRule.antecedent.category != category {
			continue
		}

		someresults, match := aRule.Matches(word)
		if match {
			rule = aRule
			binding, ok = segmenter.findBinding(someresults, aRule.antecedent.GetPattern())
			if !ok {
				continue
			}
			found = true
			break
		}
	}

	return rule, found, binding
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