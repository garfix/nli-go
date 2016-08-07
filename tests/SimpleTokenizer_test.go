package tests

import "testing"
import "strings"
import "nli/lib"

func TestSimpleTokenizer(t *testing.T) {

    rawInput := string("How old is Byron?")

    tokenizer := new(lib.SimpleTokenizer)
    wordArray := tokenizer.Process(rawInput)

    wordString := strings.Join(wordArray, "/")
    if wordString != "How/old/is/Byron/?" {
        t. Error(wordString)
    }
}
