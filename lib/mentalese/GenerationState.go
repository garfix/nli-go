package mentalese

type GenerationState struct {
	generated []Term
}

func NewGenerationState() *GenerationState {
	return &GenerationState{
		generated: []Term{},
	}
}

func (s *GenerationState) MarkGenerated(id Term) {
	s.generated = append(s.generated, id)
}

func (s *GenerationState) IsGenerated(id Term) bool {
	found := false
	for _, anId := range s.generated {
		if anId.Equals(id) {
			found = true
			break
		}
	}
	return found
}

func (s *GenerationState) Clear() {
	s.generated = []Term{}
}
