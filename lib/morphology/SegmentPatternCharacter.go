package morphology

const CharacterTypeRest = "rest"
const CharacterTypeClass = "character-class"
const CharacterTypeLiteral = "literal"

type SegmentPatternCharacter struct {
	characterType string
	characterValue string
	index int
}

func NewSegmentPatterCharacter(characterType string, characterValue string, index int) SegmentPatternCharacter {
	return SegmentPatternCharacter{
		characterType:  characterType,
		characterValue: characterValue,
		index: index,
	}
}
