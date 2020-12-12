package morphology

type SegmentPatternCharacter struct {
	characterType string
	characterValue string
}

func NewSegmentPatterCharacter(characterType string, characterValue string) SegmentPatternCharacter {
	return SegmentPatternCharacter{
		characterType:  characterType,
		characterValue: characterValue,
	}
}
