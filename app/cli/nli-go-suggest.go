package main

import (
    "fmt"
    "os"
    "nli-go/lib/parse/earley"
    "nli-go/lib/importer"
    "nli-go/lib/parse"
    "encoding/json"
)
// This application takes a partial sentence as its parameter and returns a JSON array of suggested words.
func main()  {

    if len(os.Args) != 2 {
        fmt.Println("NLI-GO Suggest")
        fmt.Println("Returns a list of suggested next words in a sentence.")
        fmt.Println("Example:")
        fmt.Println("\tnli-go-suggest \"The quick brown fox jumps\"")
        return
    }

    internalGrammarParser := importer.NewInternalGrammarParser()
    grammar := internalGrammarParser.CreateGrammar(internalGrammarParser.LoadText("../../resources/english-1.grammar"))
    lexicon := internalGrammarParser.CreateLexicon(internalGrammarParser.LoadText("../../resources/english-1.lexicon"))

    tokenizer := parse.NewTokenizer()
    parser := earley.NewParser(grammar, lexicon)

    sentence := os.Args[1]
    tokens := tokenizer.Process(sentence)
    suggests := parser.Suggest(tokens)

    b, err := json.Marshal(suggests)
    if err != nil {
        fmt.Println("error:", err)
    } else {
        fmt.Println(string(b))
    }
}