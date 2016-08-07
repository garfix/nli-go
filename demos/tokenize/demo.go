package main

import "fmt"
import "strings"
import "os"
import "nli/lib"

// Provide a sentence as command line parameters (or as a single parameter within quotes)
// and this app will provide the tokens, separarated by slashes
func main() {

    rawInput := strings.Join(os.Args[1:], " ")

    if len(rawInput) == 0 {
        fmt.Print("use: tokenizer \"Provide a sentence here\"\n")
        return
    }

    tokenizer := new(lib.SimpleTokenizer)
    wordArray := tokenizer.Process(rawInput)

    fmt.Print(strings.Join(wordArray, "/"))
    fmt.Print("\n")
}
