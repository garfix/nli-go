package main

import (
    "fmt"
    "os"
    "nli-go/lib/global"
    "nli-go/lib/common"
    "strings"
)
// This application takes a partial sentence as its parameter and returns a JSON array of suggested words.
func main()  {

    if len(os.Args) != 3 {
        fmt.Println("NLI-GO Suggest")
        fmt.Println("Returns a list of suggested next words in a sentence.")
        fmt.Println("Use:")
        fmt.Println("\tnligo-suggest /path/to/config.json \"<partial sentence>\"")
        fmt.Println("Example:")
        fmt.Println("\tnligo-suggest fox/config.json \"The quick brown fox jumps\"")
        return
    }

    configPath := os.Args[1]
    sentence := os.Args[2]

    path := configPath
    if len(path) > 0 && path[0] != os.PathSeparator {
        path = common.Dir() + string(os.PathSeparator) + configPath
    }

    log := global.NewSystemLog()
    system := global.NewSystem(path, log)
    suggests := []string{}

    if log.IsOk() {
        suggests = system.Suggest(sentence)
    }

    fmt.Fprintln(os.Stdout, strings.Join(suggests, "\n"))
    fmt.Fprintln(os.Stderr, strings.Join(log.GetLogLines(), "\n"))
}
