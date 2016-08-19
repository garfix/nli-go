package example1

// A class that provides a raw input string (that was typed in by the user)
type simpleRawInputSource struct {
	rawInput string
}

func NewSimpleRawInputSource(rawInput string) *simpleRawInputSource {
	return &simpleRawInputSource{rawInput: rawInput}
}

func (source *simpleRawInputSource) GetRawInput() string {
	return source.rawInput
}
