package tests

import (
	"nli-go/lib/example1"
	"strings"
	"testing"
)

func TestSimpleTokenizer(test *testing.T) {

	rawInputSource := example1.NewSimpleRawInputSource("How old is Byron?")
	tokenizer := example1.NewSimpleTokenizer()

	wordArray := tokenizer.Process(rawInputSource)

	wordString := strings.Join(wordArray, "/")
	if wordString != "How/old/is/Byron/?" {
		test.Error(wordString)
	}
}
