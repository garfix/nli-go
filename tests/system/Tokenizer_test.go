package tests

import (
	"strings"
	"testing"
	"nli-go/lib/parse"
)

func TestTokenizer(test *testing.T) {

	tokenizer := parse.NewTokenizer()
	wordArray := tokenizer.Process("How old is Byron?")

	wordString := strings.Join(wordArray, "/")
	if wordString != "How/old/is/Byron/?" {
		test.Error(wordString)
	}
}
