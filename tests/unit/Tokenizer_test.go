package tests

import (
	"nli-go/lib/parse"
	"strings"
	"testing"
)

func TestTokenizer(test *testing.T) {

	tokenizer := parse.NewTokenizer(parse.DefaultTokenizerExpression)

	tests := []struct {
		input            string
		expected         string
	}{
		{"How old is Byron?", "How/old/is/Byron/?"},
		{"25 years C64", "25/years/C64"},
		{"Karen Spärck Jones", "Karen/Spärck/Jones"},
		{"Düsseldorf, Köln, Москва, 北京市, إسرائيل !@#$", "Düsseldorf/,/Köln/,/Москва/,/北京市/,/إسرائيل/!/@/#/$"},
	}

	for _, aTest := range tests {
		wordArray := tokenizer.Process(aTest.input)
		wordString := strings.Join(wordArray, "/")
		if wordString != aTest.expected {
			test.Error(wordString)
		}
	}
}
