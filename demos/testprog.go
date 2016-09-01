package main

import (
	"nli-go/lib/example3"
	"fmt"
)

func main() {

	parser := example3.NewSimpleInternalGrammarParser()
	ok := true

	lexicon, _, ok := parser.CreateLexicon("" +
		"[" +
		"\t{ form: 'boek'\npos: noun }" +
		"]")
	if !ok {
		fmt.Print("Parse error")
	}

	_, ok = lexicon.GetLexItem("boek", "noun")
	if !ok {
		fmt.Print("Parse error")
	}
}