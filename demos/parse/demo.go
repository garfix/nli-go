package main

import "fmt"
import "strings"
import "os"
import "nli-go/lib/example1"

// Provide a sentence as command line parameters (or as a single parameter within quotes)
// and this app will provide the tokens, separated by slashes
func main() {

	rawInput := strings.Join(os.Args[1:], " ")

	if len(rawInput) == 0 {
		fmt.Print("use: parser \"Provide a sentence here\"\n")
		return
	}

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
		"the":   {"det"},
		"a":     {"det"},
		"shy":   {"adj"},
		"small": {"adj"},
		"boy":   {"noun"},
		"girl":  {"noun"},
		"cries": {"verb"},
		"sings": {"verb"},
	}

	inputSource := example1.NewSimpleRawInputSource(rawInput)
	tokenizer := example1.NewSimpleTokenizer()
	parser := example1.NewSimpleParser(example1.NewSimpleGrammar(rules), example1.NewSimpleLexicon(lexItems))

	wordArray := tokenizer.Process(inputSource)

	length, parseTree, ok := parser.Process(wordArray)

	if ok {
		fmt.Print("ok")

		fmt.Print(" (")
		fmt.Print(length)
		fmt.Print(")\n")

		fmt.Print(parseTree)
		fmt.Print("\n")
		fmt.Print("\n")
	} else {
		fmt.Print("not ok\n")
	}
}
