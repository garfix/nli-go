package main

import (
    "fmt"
    "os"
    "nli-go/lib/global"
    "nli-go/lib/common"
    "strings"
)
// This application takes a sentence as its parameter and returns a JSON array with a suggested answer.
func main()  {

    if len(os.Args) != 3 {
        fmt.Println("NLI-GO Answer")
        fmt.Println("Returns an answer to a user utterance.")
        fmt.Println("Use:")
        fmt.Println("\tnligo-answer /path/to/config.json <full sentence>")
        fmt.Println("Example:")
        fmt.Println("\tnligo-answer fox/config.json \"Did the quick brown jump over the lazy dog?\"")
        return
    }

    configPath := os.Args[1]
    sentence := os.Args[2]
    answer := ""

    path := configPath
    if len(path) > 0 && path[0] != os.PathSeparator {
        path = common.Dir() + string(os.PathSeparator) + configPath
    }

    log := global.NewSystemLog()
    system := global.NewSystem(path, log)

    if log.IsOk() {
        answer = system.Answer(sentence)
    }

    fmt.Fprintln(os.Stdout, answer)
    fmt.Fprintln(os.Stderr, strings.Join(log.GetLogLines(), "\n"))
}
