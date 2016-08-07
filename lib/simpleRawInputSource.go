package lib

type simpleRawInputSource struct {
    rawInput string
}

func NewSimpleRawInputSource(rawInput string) *simpleRawInputSource {
    return &simpleRawInputSource{rawInput}
}

func (source *simpleRawInputSource) GetRawInput() string {
    return source.rawInput
}
