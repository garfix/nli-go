package tests

import "testing"
import "strings"
import "nli-go/lib"

func TestSimpleTokenizer(test *testing.T) {

    rawInputSource := lib.NewSimpleRawInputSource("How old is Byron?")
    tokenizer := lib.NewSimpleTokenizer()

    wordArray := tokenizer.Process(rawInputSource)

    wordString := strings.Join(wordArray, "/")
    if wordString != "How/old/is/Byron/?" {
        test.Error(wordString)
    }
}
