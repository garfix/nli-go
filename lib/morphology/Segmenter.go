package morphology

type Segmenter struct {

}

func NewSegmenter() *Segmenter {
	return &Segmenter{}
}

func (segmenter *Segmenter) Segment(characterClasses []CharacterClass, segmentationRules []SegmentationRule, word string, category string) []string {
	return []string{}
}
