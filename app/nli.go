package main

import (
	"encoding/json"
	"fmt"
	"nli-go/lib/common"
	"nli-go/lib/global"
	"os"
	"path/filepath"
)

type Result struct {
	Success    bool
	ErrorLines []string
	Productions []string
	Value      []string
}

// This application takes a sentence as its parameter and executes a given command.
func main() {

	if len(os.Args) != 4 {
		fmt.Println("Usage: nli <command> -s <session-id> </path/to/config.json> <full sentence>")
		fmt.Println("")
		fmt.Println("Example:")
		fmt.Println("    nli answer fox/config.json \"Did the quick brown fox jump over the lazy dog?\"")
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

	log := common.NewSystemLog(false)
	system := global.NewSystem(path, log)
	config := system.ReadConfig(configPath, log)

	value := []string{}
	errorLines := []string{}

	if log.IsOk() {

		builder := global.NewSystemBuilder(filepath.Dir(configPath), log)
		builder.BuildFromConfig(system, config)

		if log.IsOk() {
			switch command {
			case "answer":
				value = []string{system.Answer(sentence)}
			case "suggest":
				value = system.Suggest(sentence)
			default:
				errorLines = []string{fmt.Sprintf("%s is not valid command.\n", os.Args[1])}
			}
		}

		if log.IsOk() {
			builder.SaveDialogContextFromPath(system, config.DialogContextPath)
		}

		errorLines = append(errorLines, log.GetErrors()...)
	}

	productions := log.GetProductions()

	result := Result{
		Success:    log.IsOk(),
		ErrorLines: errorLines,
		Productions: productions,
		Value:      value,
	}

	jsonString, _ := json.Marshal(result)
	fmt.Printf(string(jsonString) + "\n")

	if !log.IsOk() {
		os.Exit(1)
	}
}
