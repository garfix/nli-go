package lib

import "regexp"
import "nli-go/types"

type simpleTokenizer struct {

}

func NewSimpleTokenizer() *simpleTokenizer {
    return &simpleTokenizer{}
}

// A simple tokenizer.
// It treats any whitespace as token separator
// It also treats any non-word-character as a single-character token
// Returns an array of string tokens from rawInput.
func (*simpleTokenizer) Process(rawInput types.RawInputSource) []string {

    expression, _ := regexp.Compile("([\\w]+|[^\\s])")
    tokens := expression.FindAllString(rawInput.GetRawInput(), -1)

    return tokens
}
