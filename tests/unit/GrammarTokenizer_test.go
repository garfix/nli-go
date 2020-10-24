package tests

import (
	"fmt"
	"nli-go/lib/importer"
	"testing"
)

func TestGrammarTokenizer(test *testing.T) {

	tok := importer.NewGrammarTokenizer()

	// test token values

	tokens, _, _ := tok.Tokenize("dog(fido)")

	if len(tokens) != 4 {
		test.Error(fmt.Sprintf("Incorrect number of tokens: %d", len(tokens)))
	}
	if tokens[0].TokenValue != "dog" {
		test.Error(fmt.Sprintf("Error in value: %s", tokens[0].TokenValue))
	}
	if tokens[1].TokenValue != "(" {
		test.Error(fmt.Sprintf("Error in value: %s", tokens[1].TokenValue))
	}
	if tokens[2].TokenValue != "fido" {
		test.Error(fmt.Sprintf("Error in value: %s", tokens[2].TokenValue))
	}
	if tokens[3].TokenValue != ")" {
		test.Error(fmt.Sprintf("Error in value: %s", tokens[3].TokenValue))
	}

	tokens, _, _ = tok.Tokenize("name(Dog, 'Fido_dido2')")

	if len(tokens) != 6 {
		test.Error(fmt.Sprintf("Incorrect number of tokens: %d", len(tokens)))
	}
	if tokens[2].TokenValue != "Dog" {
		test.Error(fmt.Sprintf("Error in value: %s", tokens[2].TokenValue))
	}
	if tokens[4].TokenValue != "Fido_dido2" {
		test.Error(fmt.Sprintf("Error in value: %s", tokens[4].TokenValue))
	}

	tokens, _, _ = tok.Tokenize("name(Dog, 'Fido_dido2')")

	if len(tokens) != 6 {
		test.Error(fmt.Sprintf("Incorrect number of tokens: %d", len(tokens)))
	}
	if tokens[2].TokenValue != "Dog" {
		test.Error(fmt.Sprintf("Error in value: %s", tokens[2].TokenValue))
	}
	if tokens[4].TokenValue != "Fido_dido2" {
		test.Error(fmt.Sprintf("Error in value: %s", tokens[4].TokenValue))
	}

	tokens, _, _ = tok.Tokenize("mature(X) :- age(X, 18)")

	if len(tokens) != 11 {
		test.Error(fmt.Sprintf("Incorrect number of tokens: %d", len(tokens)))
	}
	if tokens[4].TokenValue != ":-" {
		test.Error(fmt.Sprintf("Error in value: %s", tokens[4].TokenValue))
	}
	if tokens[9].TokenValue != "18" {
		test.Error(fmt.Sprintf("Error in value: %s", tokens[9].TokenValue))
	}

	tokens, _, _ = tok.Tokenize("strange('it\\'s', '\\\\', '')")

	if len(tokens) != 8 {
		test.Error(fmt.Sprintf("Incorrect number of tokens: %d", len(tokens)))
	}
	if tokens[2].TokenValue != "it's" {
		test.Error(fmt.Sprintf("Error in value: %s", tokens[2].TokenValue))
	}
	if tokens[4].TokenValue != "\\" {
		test.Error(fmt.Sprintf("Error in value: %s", tokens[4].TokenValue))
	}
	if tokens[6].TokenValue != "" {
		test.Error(fmt.Sprintf("Error in value: %s", tokens[6].TokenValue))
	}

	tokens2, line, _ := tok.Tokenize(" \t spacey('it',\t'really',\n'is')\t \t")
	if len(tokens2) != 8 {
		test.Error(fmt.Sprintf("Incorrect number of tokens: %d", len(tokens2)))
	}
	if line != 2 {
		test.Error(fmt.Sprintf("Incorrect number of lines: %d", line))
	}

	_, line2, _ := tok.Tokenize("/*\nmulti line comment\n*/")
	if line2 != 3 {
		test.Error(fmt.Sprintf("Incorrect number of lines: %d", line))
	}
}
