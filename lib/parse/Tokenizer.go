package parse

import (
	"nli-go/lib/common"
	"regexp"
)

// A simple tokenizer.
// It treats any whitespace as token separator
// It also treats any non-word-character as a single-character token
// Returns an array of string tokens from rawInput.
type Tokenizer struct {
	log *common.SystemLog
}

func NewTokenizer(log *common.SystemLog) *Tokenizer {
	return &Tokenizer{log: log}
}

func (tok *Tokenizer) Process(rawInput string) []string {

	expression, _ := regexp.Compile("([\\w]+|[^\\s])")
	tokens := expression.FindAllString(rawInput, -1)

	return tokens
}
