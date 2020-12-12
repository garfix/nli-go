package morphology

import "nli-go/lib/mentalese"

type CharacterClass struct {
	name string
	characters mentalese.TermList
}

func NewCharacterClass(name string, list mentalese.TermList) CharacterClass {
	return CharacterClass{
		name:       name,
		characters: list,
	}
}