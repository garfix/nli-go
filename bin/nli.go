package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"nli-go/lib/common"
	"nli-go/lib/server"
	"os"
	"strings"
)

func main() {

	const optApplication = "a"
	const optSession = "s"
	const optOutput = "o"
	const optPort = "p"

	const txtApplication = "Directory of the application (example: resources/blocks)"
	const txtSession = "A session identifier of your choice"
	const txtOutput = "Directory to store variable working data"
	const txtPort = "Server port"

	workingDir, _ := os.Getwd()
	defaultOutputDir := workingDir + "/var"

	answerCmd := flag.NewFlagSet("answer", flag.ExitOnError)
	answerApp := answerCmd.String(optApplication, "", txtApplication)
	answerSes := answerCmd.String(optSession, "", txtSession)
	answerOut := answerCmd.String(optOutput, defaultOutputDir, txtOutput)
	answerPort := answerCmd.String(optPort, "3333", txtPort)

	interCmd := flag.NewFlagSet("inter", flag.ExitOnError)
	interApp := interCmd.String(optApplication, "", txtApplication)
	interSes := interCmd.String(optSession, "", txtSession)
	interOut := interCmd.String(optOutput, defaultOutputDir, txtOutput)
	interPort := interCmd.String(optPort, "3333", txtPort)

	resetCmd := flag.NewFlagSet("reset", flag.ExitOnError)
	resetApp := resetCmd.String(optApplication, "", txtApplication)
	resetSes := resetCmd.String(optSession, "", txtSession)
	resetOut := resetCmd.String(optOutput, defaultOutputDir, txtOutput)
	resetPort := resetCmd.String(optPort, "3333", txtPort)

	testCmd := flag.NewFlagSet("test", flag.ExitOnError)
	testApp := testCmd.String(optApplication, "", txtApplication)
	testSes := testCmd.String(optSession, "", txtSession)
	testOut := testCmd.String(optOutput, defaultOutputDir, txtOutput)
	testPort := testCmd.String(optPort, "3333", txtPort)

	queryCmd := flag.NewFlagSet("query", flag.ExitOnError)
	queryApp := queryCmd.String(optApplication, "", txtApplication)
	querySes := queryCmd.String(optSession, "", txtSession)
	queryOut := queryCmd.String(optOutput, defaultOutputDir, txtOutput)
	queryPort := queryCmd.String(optPort, "3333", txtPort)

	flag.Parse()

	if len(os.Args) < 2 {
		showUsage()
	}

	switch os.Args[1] {
	case "answer":
		answerCmd.Parse(os.Args[2:])
		if *answerApp == "" || answerCmd.Arg(0) == "" {
			showUsage()
		}
		answer(*answerApp, *answerSes, *answerOut, *answerPort, answerCmd.Arg(0))
	case "inter":
		interCmd.Parse(os.Args[2:])
		if *interApp == "" {
			showUsage()
		}
		goInteractive(*interApp, *interSes, *interOut, *interPort)
	case "query":
		queryCmd.Parse(os.Args[2:])
		if *queryApp == "" || queryCmd.Arg(0) == "" {
			showUsage()
		}
		performQuery(*queryApp, *querySes, *queryOut, *queryPort, queryCmd.Arg(0))
	case "reset":
		resetCmd.Parse(os.Args[2:])
		if *resetApp == "" || *resetSes == "" {
			showUsage()
		}
		resetSession(*resetApp, *resetSes, *resetOut, *resetPort)
	case "test":
		testCmd.Parse(os.Args[2:])
		if *testApp == "" || *testSes == "" {
			showUsage()
		}
		test(*testApp, *testSes, *testOut, *testPort)
	default:
		println("Unknown command: " + os.Args[1])
	}
}

func showUsage() {
	fmt.Println("NLI-GO client")
	fmt.Println("")
	fmt.Println("Answer a single question/command:")
	fmt.Println("    bin/nli answer <options -a,-s,-o,-p> \"Does the green block support a pyramid?\"")
	fmt.Println("")
	fmt.Println("Start an interactive session:")
	fmt.Println("    bin/nli inter <options -a,-s,-o,-p>")
	fmt.Println("")
	fmt.Println("Low level query:")
	fmt.Println("    bin/nli query <options -a,-s,-o,-p> \"dom:at(E, X, Z, Y) dom:type(E, Type) dom:color(E, Color) dom:size(E, Width, Length, Height)\"")
	fmt.Println("")
	fmt.Println("Reset dialog context:")
	fmt.Println("    bin/nli reset <options -a,-s,-o,-p>")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -a </path/to/application>")
	fmt.Println("  -s <session_id>")
	fmt.Println("  -o </path/to/generated/output>")
	fmt.Println("  -p <server_port>")
	fmt.Println("")
	os.Exit(1)
}

func answer(appDir string, sessionId string, workDir string, port string, sentence string) {

	cwd, _ := os.Getwd()
	request := server.Request{
		SessionId:      sessionId,
		ApplicationDir: common.AbsolutePath(cwd, appDir),
		WorkDir:        common.AbsolutePath(cwd, workDir),
		Command:        "answer",
		Query:          sentence,
	}

	requestString, _ := json.Marshal(request)

	connection, err := net.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Fprintf(connection, string(requestString)+"\n")

	resultString, _ := ioutil.ReadAll(connection)

	result := server.ResponseAnswer{}
	json.Unmarshal(resultString, &result)

	if result.Success {
		println(result.Answer)
	} else {
		println(strings.Join(result.ErrorLines, "\n"))
	}
}

func performQuery(appDir string, sessionId string, workDir string, port string, query string) {

	cwd, _ := os.Getwd()
	request := server.Request{
		SessionId:      sessionId,
		ApplicationDir: common.AbsolutePath(cwd, appDir),
		WorkDir:        common.AbsolutePath(cwd, workDir),
		Command:        "query",
		Query:          query,
	}

	requestString, _ := json.Marshal(request)

	connection, err := net.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Fprintf(connection, string(requestString)+"\n")

	resultString, err2 := ioutil.ReadAll(connection)
	if err2 != nil {
		println(err2.Error())
		os.Exit(1)
	}

	println(string(resultString))
}

func goInteractive(appDir string, sessionId string, workDir string, port string) {

	sentence := ""
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("NLI-GO session with " + appDir + ". Type 'exit' to stop, 'reset' to start fresh.")

	for true {

		fmt.Print("\n> ")
		sentence, _ = reader.ReadString('\n')
		sentence = strings.Trim(sentence, "\n")

		if sentence == "exit" {
			break
		}
		if sentence == "reset" {
			resetSession(appDir, sessionId, workDir, port)
			println("OK")
			continue
		}

		answer(appDir, sessionId, workDir, port, sentence)
	}
}

func resetSession(appDir string, sessionId string, workDir string, port string) {

	cwd, _ := os.Getwd()
	request := server.Request{
		SessionId:      sessionId,
		ApplicationDir: common.AbsolutePath(cwd, appDir),
		WorkDir:        common.AbsolutePath(cwd, workDir),
		Command:        "reset",
	}

	requestString, _ := json.Marshal(request)

	connection, err := net.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Fprintf(connection, string(requestString)+"\n")

	_, err2 := ioutil.ReadAll(connection)
	if err2 != nil {
		println(err2.Error())
		os.Exit(1)
	}

	println("OK")
}

func test(appDir string, sessionId string, workDir string, port string) {

	cwd, _ := os.Getwd()
	request := server.Request{
		SessionId:      sessionId,
		ApplicationDir: common.AbsolutePath(cwd, appDir),
		WorkDir:        common.AbsolutePath(cwd, workDir),
		Command:        "test",
	}

	requestString, _ := json.Marshal(request)

	connection, err := net.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Fprintf(connection, string(requestString)+"\n")

	resultString, err2 := ioutil.ReadAll(connection)
	if err2 != nil {
		println(err2.Error())
		os.Exit(1)
	}

	result := server.ResponseAnswer{}
	json.Unmarshal(resultString, &result)

	if result.Success {
		println(result.Answer)
	} else {
		println(strings.Join(result.ErrorLines, "\n"))
	}
}
