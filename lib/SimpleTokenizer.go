package lib

import "regexp"

type SimpleTokenizer struct {

}

// A simple tokenizer.
// It treats any whitespace as token separater
// It also treats any non-word-character as a single-character token
// Returns an array of string tokens from rawInput.
func (*SimpleTokenizer) Process(rawInput string) []string {

    expression, _ := regexp.Compile("([\\w]+|[^\\s])")
    tokens := expression.FindAllString(rawInput, -1)

    return tokens
}
