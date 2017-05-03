package main

import (
    "fmt"
    "os"
    "nli-go/lib/global"
    "nli-go/lib/common"
    "encoding/json"
)

type Result struct {
    Error bool
    ErrorLines []string
    Value []string
}

// This application takes a sentence as its parameter and executes a given command.
func main()  {

    if len(os.Args) != 4 {
        fmt.Println("Usage: nli <command> </path/to/config.json> <full sentence>")
        fmt.Println("")
        fmt.Println("Example:")
        fmt.Println("    nli answer fox/config.json \"Did the quick brown jump over the lazy dog?\"")
        fmt.Println("")
        fmt.Println("Commands:")
        fmt.Println("    answer     Return an answer to <full sentence>")
        fmt.Println("    suggest    Returns next word suggestions")
        return
    }

    command := os.Args[1]
    configPath := os.Args[2]
    sentence := os.Args[3]

    path := configPath
    if len(path) > 0 && path[0] != os.PathSeparator {
        path = common.Dir() + string(os.PathSeparator) + configPath
    }

    log := global.NewSystemLog()
    system := global.NewSystem(path, log)

    ok := log.IsOk()
    value := []string{}
    errorLines := []string{}

    if ok {
        switch command {
        case "answer":
            value = []string{ system.Answer(sentence) }
        case "suggest":
            value = system.Suggest(sentence)
        default:
            errorLines = []string{ fmt.Sprintf("%s is not valid command.\n", os.Args[1]) }
            ok = false
        }
    }

    if !log.IsOk() {
        errorLines = append(errorLines, log.GetLogLines()...)
    }

    result := Result{
        Error:      !ok,
        ErrorLines: errorLines,
        Value:      value,
    }

    jsonString, _ := json.Marshal(result)
    fmt.Printf(string(jsonString))

    if !ok {
        os.Exit(1)
    }
}
