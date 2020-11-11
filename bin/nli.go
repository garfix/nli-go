package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"nli-go/lib/common"
	"nli-go/lib/global"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Result struct {
	Success      bool
	ErrorLines   []string
	Productions  []string
	Answer       string
	OptionKeys   []string
	OptionValues []string
}

func main() {

	var interactive = false
	var sessionId = ""
	var absSessionPath = ""
	var configPath = ""
	var variableDir = ""
	var returnType = ""

	flag.BoolVar(&interactive, "i", false, "Interative: start a session with NLI-GO")
	flag.StringVar(&sessionId, "s", "", "Session id: an arbitrary identifier for current user's dialog session")
	flag.StringVar(&configPath, "c", "", "Config path: (relative) path to the root directory of an application")
	flag.StringVar(&variableDir, "d", "", "Directory: (relative) path to the directory where caches and sessions are stored (default ./var)")
	flag.StringVar(&returnType, "r", "text", "Return type: text / json")

	flag.Parse()

	if len(flag.Args()) == 0 && flag.NFlag() == 0 {
		fmt.Println("Use:")
		fmt.Println(" nli -c </path/to/application> [-i] [-s <session_id>] [-r JSON] [-d </path/to/generated/output>] <full sentence>")
		fmt.Println("")
		fmt.Println("Single question:")
		fmt.Println("    bin/nli -c resources/blocks \"Does the green block support a pyramid?\"")
		fmt.Println("")
		fmt.Println("Interactive:")
		fmt.Println("    bin/nli -i -c resources/blocks")
		fmt.Println("")
		fmt.Println("Type `nli/go --help` for more information.")
		fmt.Println("")
		return
	}

	sentence := flag.Arg(0)
	workingDir, _ := os.Getwd()
	absConfigPath := common.AbsolutePath(workingDir, configPath)

	if variableDir == "" {
		variableDir = workingDir + "/var"
	} else {
		variableDir, _ = filepath.Abs(variableDir)
		variableDir = filepath.Clean(variableDir)
	}

	log := common.NewSystemLog()
	system := global.NewSystem(absConfigPath, variableDir, log)

	// load dialog context
	if sessionId != "" {
		system.PopulateDialogContext(sessionId, true)
	}

	if log.IsOk() {
		if interactive {
			goInteractive(system, log, absSessionPath, configPath, returnType)
		} else {
			singleLine(system, log, sentence, sessionId, returnType)
		}
	}

	if !log.IsOk() {
		os.Exit(1)
	}
}

func singleLine(system *global.System, log *common.SystemLog, sentence string, sessionId string, returnType string) (string, *common.Options) {

	// the actual system call
	answer, options := system.Answer(sentence)

	// store dialog context for next call
	if sessionId != "" {
		system.StoreDialogContext(sessionId)
	}

	response := createResponseString(log, answer, options, returnType)

	fmt.Printf(response)

	return answer, options
}

func goInteractive(system *global.System, log *common.SystemLog, absSessionPath string, configPath string, returnType string) {

	sentence := ""
	options := &common.Options{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("NLI-GO session with " + configPath + ". Type 'exit' to stop.")

	for true {

		fmt.Print("\n> ")
		sentence, _ = reader.ReadString('\n')
		sentence = strings.Trim(sentence, "\n")

		if sentence == "exit" {
			break
		}

		if options.HasOptions() {
			index, err := strconv.Atoi(sentence)
			if err == nil {
				sentence = optionIndexToOptionKey(index, options)
			}
		}

		_, options = singleLine(system, log, sentence, absSessionPath, returnType)

	}
}

func optionIndexToOptionKey(index int, options *common.Options) string {

	for i, key := range options.GetKeys() {
		if i == index - 1 {
			return key
		}
	}

	return ""
}

func createResponseString(log *common.SystemLog, answer string, options *common.Options, returnType string) string {

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
			values := options.GetValues()
			for i := range options.GetKeys() {
				response += strconv.Itoa(i + 1) + ") " + values[i] + "\n"
			}

		} else {
			for _, err := range log.GetErrors() {
				response += err + "\n"
			}
		}
	}

	return response
}