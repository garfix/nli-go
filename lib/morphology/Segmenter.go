package morphology

import "strconv"

type Segmenter struct {

}

func NewSegmenter() *Segmenter {
	return &Segmenter{}
}

func (segmenter *Segmenter) Segment(segmentationRules []SegmentationRule, word string, category string) []string {

	segments := []string{}

	_, found, _ := segmenter.findRule(segmentationRules, word, category)

	if found {

	}

	return segments
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
			variable := character.characterValue + strconv.Itoa(character.index)
			value := results[i + 1]

			// check if the variable exists and has the same value
			existing, found := binding[variable]
			if found && existing != value {
				ok = false
				break
			}

			binding[variable] = value
		}
	}

	return binding, ok
}