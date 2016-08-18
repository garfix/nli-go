package tests

import (
    "testing"
    "fmt"
    "nli-go/lib/example1"
)

func TestSimpleParser(test *testing.T) {

    rules := map[string][][]string{
        "S": {
            {"NP", "VP"},
        },
        "NP": {
            {"NBar"},
            {"det", "NBar"},
        },
        "NBar": {
            {"noun"},
            {"adj", "NBar"},
        },
        "VP": {
            {"verb"},
        },
    }

    lexItems := map[string][]string{
        "the": {"det"},
        "a": {"det"},
        "shy": {"adj"},
        "small": {"adj"},
        "boy": {"noun"},
        "girl": {"noun"},
        "cries": {"verb"},
        "sings": {"verb"},
    }

    rawInput := "the small shy girl sings"
    inputSource := example1.NewSimpleRawInputSource(rawInput)
    tokenizer := example1.NewSimpleTokenizer()
    parser := example1.NewSimpleParser(example1.NewSimpleGrammar(rules), example1.NewSimpleLexicon(lexItems))

    wordArray := tokenizer.Process(inputSource)
    length, parseTree, ok := parser.Process(wordArray)

    if !ok {
        test.Error("Parse failed")
    }
    if length != 5 {
        test.Error(fmt.Sprintf("Length not equal to 5: %d", length))
    }
    if parseTree.SyntacticCategory != "S" {
        test.Error("Missing S")
    }
    if parseTree.Children[1].SyntacticCategory != "VP" {
        test.Error("Missing VP")
    }
    if parseTree.Children[0].Children[1].Children[0].SyntacticCategory != "adj" {
        test.Error("Missing adj")
    }
    if parseTree.Children[0].Children[1].Children[0].Word != "small" {
        test.Error("Wrong word")
    }
}