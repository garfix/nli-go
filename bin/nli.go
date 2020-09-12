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
	Success      bool
	ErrorLines   []string
	Productions  []string
	Answer       string
	OptionKeys   []string
	OptionValues []string
}

// This application takes a sentence as its parameter and executes a given command.
func main() {

	var sessionId = ""
	var absSessionPath = ""
	var configPath = ""
	var returnType = ""

	answer := ""
	options := common.NewOptions()

	flag.StringVar(&sessionId, "s", "", "Session id: an arbitrary identifier for current user's dialog session")
	flag.StringVar(&configPath, "c", "", "Config path: (relative) path to the root directory of an application")
	flag.StringVar(&returnType, "r", "text", "Return type (text / json)")

	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Usage: nli [-s <session_id>] [-r JSON] -c </path/to/application> <full sentence>")
		fmt.Println("")
		fmt.Println("Example:")
		fmt.Println("    nli -s 73926642 -c fox/config.yml \"Did the quick brown fox jump over the lazy dog?\"")
		fmt.Println("")
		return
	}

	sentence := flag.Arg(0)
	workingDir, _ := os.Getwd()
	absConfigPath := common.AbsolutePath(workingDir, configPath)
	log := common.NewSystemLog(false)
	system := global.NewSystem(absConfigPath, log)

	// load dialog context
	if sessionId != "" {

		executable, _ := os.Executable()
		executablePath := filepath.Dir(executable)

		absSessionPath = common.AbsolutePath(executablePath, "sessions/" + sessionId + ".json")
		system.PopulateDialogContext(absSessionPath, true)
	}

	if !log.IsOk() {
		goto done
	}

	// the actual system call
	answer, options = system.Answer(sentence)

	// store dialog context for next call
	if sessionId != "" {
		system.StoreDialogContext(absSessionPath)
	}

	done:

	response := ""

	if returnType == "json" || returnType == "JSON" {
		result := Result{
			Success:      log.IsOk(),
			ErrorLines:   log.GetErrors(),
			Productions:  log.GetProductions(),
			Answer:       answer,
			OptionKeys:   options.GetKeys(),
			OptionValues: options.GetValues(),
		}

		responseRaw, _ := json.MarshalIndent(result, "", "    ")
		response = string(responseRaw) + "\n"
	} else {
		if log.IsOk() {
			response = answer + "\n"
		} else {
			for _, err := range log.GetErrors() {
				response += err + "\n"
			}
		}
	}

	fmt.Printf(response)

	if !log.IsOk() {
		os.Exit(1)
	}
}
