package morphology

type Segmenter struct {
	segmentationRules *SegmentationRules
}

func NewSegmenter(segmentationRules *SegmentationRules) *Segmenter {
	return &Segmenter{
		segmentationRules: segmentationRules,
	}
}

func (segmenter *Segmenter) Segment(word string, category string) []string {

	segments := []string{}

	rules, bindings := segmenter.segmentationRules.FindRules(word, category)

	for i, rule := range rules {

		singleRuleSegments := []string{}
		binding := bindings[i]

		if rule.IsTerminal() {
			segments = []string{ word }
			break
		}

		ok := true
		for _, consequent := range rule.GetConsequents() {
			consequentSegments := segmenter.buildSegments(consequent, binding)
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

func (segmenter *Segmenter) buildSegments(segmentNode SegmentNode, binding map[string]string) []string {

	word := segmenter.buildWord(segmentNode.GetPattern(), binding)

	return segmenter.Segment(word, segmentNode.category)
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
