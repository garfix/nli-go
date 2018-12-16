package main

import (
	"encoding/json"
	"flag"
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

	var sessionId = ""
	var absSessionPath = ""
	var configPath = ""

	value := []string{}

	flag.StringVar(&sessionId, "s", "", "Session id: an arbitrary identifier for current user's dialog context")
	flag.StringVar(&configPath, "c", "", "Config path: (relative) path to a JSON nli-go config file")

	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Usage: nli [-s <session_id>] [-c </path/to/config.json>] <full sentence>")
		fmt.Println("")
		fmt.Println("Example:")
		fmt.Println("    nli -s 73926642 -c fox/config.json \"Did the quick brown fox jump over the lazy dog?\"")
		fmt.Println("")
		return
	}

	sentence := flag.Arg(0)
	absConfigPath := common.AbsolutePath(common.Dir(), configPath)
	log := common.NewSystemLog(false)
	system := global.NewSystem(absConfigPath, log)

	// load dialog context
	if sessionId != "" {

		executable, _ := os.Executable()
		executablePath := filepath.Dir(executable)

		absSessionPath = common.AbsolutePath(executablePath, "sessions/" + sessionId + ".json")
		system.PopulateDialogContext(absSessionPath)
	}

	if !log.IsOk() {
		goto done
	}

	// the actual system call
	value = []string{system.Answer(sentence)}

	// store dialog context for next call
	if sessionId != "" {
		system.StoreDialogContext(absSessionPath)
	}

	if log.IsOk() {
		goto done
	}

	done:

	result := Result{
		Success: log.IsOk(),
		ErrorLines: log.GetErrors(),
		Productions: log.GetProductions(),
		Value: value,
	}

	jsonString, _ := json.Marshal(result)
	fmt.Printf(string(jsonString) + "\n")

	if !log.IsOk() {
		os.Exit(1)
	}
}
