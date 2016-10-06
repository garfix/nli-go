package tests

import (
)
import "testing"

func TestSimpleParser(test *testing.T) {

	//rules := map[string][][]string{
	//	"S": {
	//		{"NP", "VP"},
	//	},
	//	"NP": {
	//		{"NBar"},
	//		{"det", "NBar"},
	//	},
	//	"NBar": {
	//		{"noun"},
	//		{"adj", "NBar"},
	//	},
	//	"VP": {
	//		{"verb"},
	//	},
	//}
	//
	//lexItems := map[string][]string{
	//	"the":   {"det"},
	//	"a":     {"det"},
	//	"shy":   {"adj"},
	//	"small": {"adj"},
	//	"boy":   {"noun"},
	//	"girl":  {"noun"},
	//	"cries": {"verb"},
	//	"sings": {"verb"},
	//}
	//
	//grammar := natlang.NewSimpleGrammar()
	//for _, rule := range rules {
	//	grammar.AddRule(rule)
	//}
	//
	//rawInput := "the small shy girl sings"
	//tokenizer := natlang.NewSimpleTokenizer()
	//
	//parser := natlang.NewSimpleParser(grammar, example1.NewSimpleLexicon(lexItems))
	//
	//wordArray := tokenizer.Process(rawInput)
	//length, relations, ok := parser.Process(wordArray)
	//
	//if !ok {
	//	test.Error("Parse failed")
	//}
	//if relations != 3 {
	//	test.Error(fmt.Sprintf("Length not equal to 5: %d", length))
	//}
}
