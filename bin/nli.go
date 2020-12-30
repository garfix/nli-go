package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"nli-go/lib/common"
	"nli-go/lib/global"
	"nli-go/lib/mentalese"
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

	const optApplication = "a"
	const optSession = "s"
	const optReturnType = "r"
	const optOutput = "o"

	const txtApplication = "Directory of the application (example: resources/blocks)"
	const txtSession = "A session identifier of your choice"
	const txtReturnType = "Return type (text|json)"
	const txtOutput = "Directory to store variable working data"

	workingDir, _ := os.Getwd()
	defaultOutputDir := workingDir + "/var"

	answerCmd := flag.NewFlagSet("answer", flag.ExitOnError)
	answerApp := answerCmd.String(optApplication, "", txtApplication)
	answerSes := answerCmd.String(optSession, "", txtSession)
	answerRet := answerCmd.String(optReturnType, "text", txtReturnType)
	answerOut := answerCmd.String(optOutput, defaultOutputDir, txtOutput)

	interCmd := flag.NewFlagSet("inter", flag.ExitOnError)
	interApp := interCmd.String(optApplication, "", txtApplication)
	interSes := interCmd.String(optSession, "", txtSession)
	interOut := interCmd.String(optOutput, defaultOutputDir, txtOutput)

	resetCmd := flag.NewFlagSet("reset", flag.ExitOnError)
	resetApp := resetCmd.String(optApplication, "", txtApplication)
	resetSes := resetCmd.String(optSession, "", txtSession)
	resetOut := resetCmd.String(optOutput, defaultOutputDir, txtOutput)

	queryCmd := flag.NewFlagSet("query", flag.ExitOnError)
	queryApp := queryCmd.String(optApplication, "", txtApplication)
	querySes := queryCmd.String(optSession, "", txtSession)
	queryOut := queryCmd.String(optOutput, defaultOutputDir, txtOutput)

	flag.Parse()

	if len(os.Args) < 2 {
		showUsage()
	}

	log := common.NewSystemLog()

	if log.IsOk() {
		switch os.Args[1] {
		case "answer":
			answerCmd.Parse(os.Args[2:])
			if *answerApp == "" || answerCmd.Arg(0) == "" {
				showUsage()
			}
			system := buildSystem(log, workingDir, *answerApp, *answerSes, *answerOut)
			answer(system, log, answerCmd.Arg(0), *answerRet)
		case "inter":
			interCmd.Parse(os.Args[2:])
			if *interApp == "" {
				showUsage()
			}
			system := buildSystem(log, workingDir, *interApp, *interSes, *interOut)
			goInteractive(system, log, *interApp)
		case "query":
			queryCmd.Parse(os.Args[2:])
			if *queryApp == "" || queryCmd.Arg(0) == "" {
				showUsage()
			}
			system := buildSystem(log, workingDir, *queryApp, *querySes, *queryOut)
			performQuery(system, queryCmd.Arg(0))
		case "reset":
			resetCmd.Parse(os.Args[2:])
			if *resetApp == "" || *resetSes == "" {
				showUsage()
			}
			system := buildSystem(log, workingDir, *resetApp, *resetSes, *resetOut)
			resetDialog(system)
		default:
			println("Unknown command: " + os.Args[1])
		}
	} else {
		for _, error := range log.GetErrors() {
			fmt.Println(error)
		}
	}

	if !log.IsOk() {
		os.Exit(1)
	}
}

func showUsage()  {
	fmt.Println("NLI-GO")
	fmt.Println("")
	fmt.Println("Answer a single question/command:")
	fmt.Println("    bin/nli answer <options -a,-s,-r,-o> \"Does the green block support a pyramid?\"")
	fmt.Println("")
	fmt.Println("Reset dialog context:")
	fmt.Println("    bin/nli reset <options -a,-s,-o>")
	fmt.Println("")
	fmt.Println("Start an interactive session:")
	fmt.Println("    bin/nli inter <options -a,-s,-o>")
	fmt.Println("")
	fmt.Println("Low level query:")
	fmt.Println("    bin/nli query <options -a,-s,-o> \"dom:at(E, X, Z, Y) dom:type(E, Type) dom:color(E, Color) dom:size(E, Width, Length, Height)\"")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -a </path/to/application>")
	fmt.Println("  -s <session_id>")
	fmt.Println("  -r <JSON|text>")
	fmt.Println("  -o </path/to/generated/output>")
	fmt.Println("")
	os.Exit(1)
}

func buildSystem(log *common.SystemLog, workingDir string, applicationPath string, sessionId string, outputDir string) *global.System {
	absApplicationPath := common.AbsolutePath(workingDir, applicationPath)

	outputDir, _ = filepath.Abs(outputDir)
	outputDir = filepath.Clean(outputDir)

	system := global.NewSystem(absApplicationPath, sessionId, outputDir, log)

	return system
}

func performQuery(system *global.System, query string)  {

	// the actual system call
	bindings := system.Query(query)

	response := bindingsToJson(bindings) + "\n"

	fmt.Printf(response)
}


func bindingsToJson(set mentalese.BindingSet) string {

	type aMap = map[string]string
	type array = []aMap

	arr := array{}

	for _, item := range set.GetAll() {
		i := aMap{}
		for k, v := range item.GetAll() {
			i[k] = v.String()
		}
		arr = append(arr, i)
	}


	responseRaw, _ := json.MarshalIndent(arr, "", "    ")

	return string(responseRaw)
}

func answer(system *global.System, log *common.SystemLog, sentence string, returnType string) (string, *common.Options) {

	// the actual system call
	answer, options := system.Answer(sentence)

	response := createResponseString(log, answer, options, returnType)

	fmt.Printf(response)

	return answer, options
}

func goInteractive(system *global.System, log *common.SystemLog, configPath string) {

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

		_, options = answer(system, log, sentence, "text")

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

func resetDialog(system *global.System) {
	system.ClearDialogContext()
}