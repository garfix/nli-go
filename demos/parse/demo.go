package main

import "fmt"
import "strings"
import "os"
import "nli-go/lib"

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
        "the": {"det"},
        "a": {"det"},
        "shy": {"adj"},
        "small": {"adj"},
        "boy": {"noun"},
        "girl": {"noun"},
        "cries": {"verb"},
        "sings": {"verb"},
    }

    inputSource := lib.NewSimpleRawInputSource(rawInput)
    tokenizer := lib.NewSimpleTokenizer()
    parser := lib.NewSimpleParser(lib.NewSimpleGrammar(rules), lib.NewSimpleLexicon(lexItems))

    wordArray := tokenizer.Process(inputSource)

    success := parser.Process(wordArray)
    //parseTree :=

    if (success) {
        fmt.Print("ok")
    } else {
        fmt.Print("not ok")
    }
    fmt.Print("\n")
}
