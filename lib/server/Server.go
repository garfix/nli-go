package server

// idea from: https://gist.github.com/miguelmota/301340db93de42b537df5588c1380863

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"nli-go/lib/common"
	"nli-go/lib/global"
	"path/filepath"
	"runtime/debug"
	"strings"
)

type Server struct {
	port    string
	systems map[string]*global.System
}

func NewServer(port string) *Server {
	return &Server{
		port:    port,
		systems: map[string]*global.System{},
	}
}

func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "localhost", server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		server.handleRequest(conn)
	}
}

func (server *Server) handleRequest(conn net.Conn) {

	if r := recover(); r != nil {
		errorString := fmt.Sprintf("%s\n%s", r, debug.Stack())
		response := Response{
			Success:    false,
			ErrorLines: strings.Split(errorString, "\n"),
		}
		responseJSON, _ := json.Marshal(response)
		conn.Write(responseJSON)
		conn.Close()
		return
	}

	reader := bufio.NewReader(conn)
	requestJson, err := reader.ReadBytes('\n')
	if err != nil {
		response := Response{
			Success:    false,
			ErrorLines: []string{err.Error()},
		}
		responseJSON, _ := json.Marshal(response)
		conn.Write(responseJSON)
		conn.Close()
		return
	}

	request := Request{}
	err = json.Unmarshal(requestJson, &request)
	if err != nil {
		response := Response{
			Success:    false,
			ErrorLines: []string{"Request could not be parsed"},
		}
		responseJSON, _ := json.Marshal(response)
		conn.Write(responseJSON)
		conn.Close()
		return
	}

	fmt.Printf("%s\n", request)

	system := server.getSystem(request)
	if !system.GetLog().IsOk() {
		response := Response{
			Success:    false,
			ErrorLines: system.GetLog().GetErrors(),
		}
		responseJSON, _ := json.Marshal(response)
		conn.Write(responseJSON)
		conn.Close()
		return
	}

	switch request.Command {
	case "send":
		client := &RequestHandler{conn: conn}
		go client.handleSend(system, request.Message)

	case "reset":
		delete(server.systems, request.SessionId)
		response := Response{
			Success: true,
		}
		responseJSON, _ := json.Marshal(response)
		conn.Write(responseJSON)
		conn.Close()

	case "query":
		client := &RequestHandler{conn: conn}
		go client.handleQuery(system, request.Query)

	case "answer":
		client := &RequestHandler{conn: conn}
		go client.handleAnswer(system, request.Query)

	case "test":
		client := &RequestHandler{conn: conn}
		go client.performTests(system, request.ApplicationDir)

	default:
		response := Response{
			Success:    false,
			ErrorLines: []string{"Unknown command: " + request.Command},
		}
		responseJSON, _ := json.Marshal(response)
		conn.Write(responseJSON)
		conn.Close()
	}
}

func (server *Server) getSystem(request Request) *global.System {
	system, found := server.systems[request.SessionId]
	if !found {
		system = buildSystem(request.WorkDir, request.ApplicationDir, request.SessionId, request.WorkDir)
		server.systems[request.SessionId] = system
	}
	return system
}

func buildSystem(workingDir string, applicationPath string, sessionId string, outputDir string) *global.System {
	absApplicationPath := common.AbsolutePath(workingDir, applicationPath)

	systemLog := common.NewSystemLog()

	outputDir, _ = filepath.Abs(outputDir)
	outputDir = filepath.Clean(outputDir)

	system := global.NewSystem(absApplicationPath, sessionId, outputDir, systemLog)

	return system
}
