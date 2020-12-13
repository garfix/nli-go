package morphology

type SegmentNode struct {
	category string
	pattern []SegmentPatternCharacter
}

func NewSegmentNode(category string, pattern []SegmentPatternCharacter) SegmentNode {
	return SegmentNode{
		category: category,
		pattern:  pattern,
	}
}

func (s SegmentNode) GetPattern() []SegmentPatternCharacter {
	return s.pattern
}
