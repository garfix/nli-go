package tests

import "testing"
import "strings"
import "nli-go/lib"

func TestSimpleTokenizer(t *testing.T) {

    rawInput := string("How old is Byron?")

    tokenizer := lib.NewSimpleTokenizer()
    wordArray := tokenizer.Process(rawInput)

    wordString := strings.Join(wordArray, "/")
    if wordString != "How/old/is/Byron/?" {
        t. Error(wordString)
    }
}
