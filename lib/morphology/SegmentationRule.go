package morphology

type SegmentationRule struct {
	antecedent SegmentNode
	consequents []SegmentNode
}

func NewSegmentationRule(antecedent SegmentNode, consequents []SegmentNode) SegmentationRule {
	return SegmentationRule{
		antecedent:  antecedent,
		consequents: consequents,
	}
}