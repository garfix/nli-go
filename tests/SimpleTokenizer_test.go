package tests

import (
	"strings"
	"testing"
	"nli-go/lib/natlang"
)

func TestSimpleTokenizer(test *testing.T) {

	tokenizer := natlang.NewSimpleTokenizer()
	wordArray := tokenizer.Process("How old is Byron?")

	wordString := strings.Join(wordArray, "/")
	if wordString != "How/old/is/Byron/?" {
		test.Error(wordString)
	}
}
