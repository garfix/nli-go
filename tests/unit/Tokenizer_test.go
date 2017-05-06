package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/parse"
	"strings"
	"testing"
)

func TestTokenizer(test *testing.T) {

	log := common.NewSystemLog(false)
	tokenizer := parse.NewTokenizer(log)
	wordArray := tokenizer.Process("How old is Byron?")

	wordString := strings.Join(wordArray, "/")
	if wordString != "How/old/is/Byron/?" {
		test.Error(wordString)
	}
}
