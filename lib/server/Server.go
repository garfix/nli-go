package server

// idea from: https://gist.github.com/miguelmota/301340db93de42b537df5588c1380863

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"nli-go/lib/common"
	"nli-go/lib/global"
	"nli-go/lib/mentalese"
	"path/filepath"

	"golang.org/x/net/websocket"
)

type Server struct {
	httpServer *http.Server
	systems    map[string]*global.System
}

func NewServer(port string) *Server {
	server := &http.Server{Addr: ":" + port}
	return &Server{
		httpServer: server,
		systems:    map[string]*global.System{},
	}
}

func (server *Server) Run() {

	http.Handle("/", websocket.Handler(server.HandleSingleConnection))

	go func() {

		println("NLI-GO server listening on " + server.httpServer.Addr + "\n")

		err := server.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic("ListenAndServe: " + err.Error())
		}

		println("Server closed")
	}()
}

func (server *Server) HandleSingleConnection(conn *websocket.Conn) {

	// println(request.Message.String())

	//io.Copy(ws, ws)

	// response := Response{
	// 	ErrorLines: []string{"niet ok"},
	// }
	// websocket.JSON.Send(ws, response)

	for true {

		request := mentalese.Request{}

		websocket.JSON.Receive(conn, &request)

		fmt.Printf("%s\t%s\t%s\n", request.SessionId, request.Command, request.Query+request.Message.String())

		system := server.getSystem(request, conn)
		switch request.Command {
		case "send":
			// client := &RequestHandler{conn: conn}
			// go client.handleSend(system, request.Message)
			system.GetClientConnector().SendToProcess(request.ProcessType, request.Message)

		case "reset":
			delete(server.systems, request.SessionId)
			response := mentalese.Response{
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
			// client := &RequestHandler{conn: conn}
			// go client.performTests(system, request.ApplicationDir)

		default:
			response := mentalese.Response{
				Success:    false,
				ErrorLines: []string{"Unknown command: " + request.Command},
			}
			responseJSON, _ := json.Marshal(response)
			conn.Write(responseJSON)
			conn.Close()
		}

	}
}

func (server *Server) Close() {
	server.httpServer.Shutdown(context.TODO())
}

// func (server *Server) Run() {
// 	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "localhost", server.port))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer listener.Close()

// 	println("NLI-GO server listening on port " + server.port + "\n")

// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		server.handleRequest(conn)
// 	}
// }

// func (server *Server) handleRequest(conn net.Conn) {

// 	if r := recover(); r != nil {
// 		errorString := fmt.Sprintf("%s\n%s", r, debug.Stack())
// 		response := Response{
// 			Success:    false,
// 			ErrorLines: strings.Split(errorString, "\n"),
// 		}
// 		responseJSON, _ := json.Marshal(response)
// 		conn.Write(responseJSON)
// 		conn.Close()
// 		return
// 	}

// 	reader := bufio.NewReader(conn)
// 	requestJson, err := reader.ReadBytes('\n')
// 	if err != nil {
// 		response := Response{
// 			Success:    false,
// 			ErrorLines: []string{err.Error()},
// 		}
// 		responseJSON, _ := json.Marshal(response)
// 		conn.Write(responseJSON)
// 		conn.Close()
// 		return
// 	}

// 	request := Request{}
// 	err = json.Unmarshal(requestJson, &request)
// 	if err != nil {
// 		response := Response{
// 			Success:    false,
// 			ErrorLines: []string{"Request could not be parsed"},
// 		}
// 		responseJSON, _ := json.Marshal(response)
// 		conn.Write(responseJSON)
// 		conn.Close()
// 		return
// 	}

// 	fmt.Printf("%s\t%s\t%s\n", request.SessionId, request.Command, request.Query+request.Message.String())

// 	system := server.getSystem(request)
// 	if !system.GetLog().IsOk() {
// 		response := Response{
// 			Success:    false,
// 			ErrorLines: system.GetLog().GetErrors(),
// 		}
// 		responseJSON, _ := json.Marshal(response)
// 		conn.Write(responseJSON)
// 		conn.Close()
// 		return
// 	}

// 	switch request.Command {
// 	case "send":
// 		client := &RequestHandler{conn: conn}
// 		go client.handleSend(system, request.Message)

// 	case "reset":
// 		delete(server.systems, request.SessionId)
// 		response := Response{
// 			Success: true,
// 		}
// 		responseJSON, _ := json.Marshal(response)
// 		conn.Write(responseJSON)
// 		conn.Close()

// 	case "query":
// 		client := &RequestHandler{conn: conn}
// 		go client.handleQuery(system, request.Query)

// 	case "answer":
// 		client := &RequestHandler{conn: conn}
// 		go client.handleAnswer(system, request.Query)

// 	case "test":
// 		client := &RequestHandler{conn: conn}
// 		go client.performTests(system, request.ApplicationDir)

// 	default:
// 		response := Response{
// 			Success:    false,
// 			ErrorLines: []string{"Unknown command: " + request.Command},
// 		}
// 		responseJSON, _ := json.Marshal(response)
// 		conn.Write(responseJSON)
// 		conn.Close()
// 	}
// }

func (server *Server) getSystem(request mentalese.Request, conn *websocket.Conn) *global.System {
	system, found := server.systems[request.SessionId]
	if !found {

		applicationDir := common.Dir() + "/../../resources/blocks"
		workDir := common.Dir() + "/../../var"

		system = buildSystem(workDir, applicationDir, request.SessionId, workDir, conn)
		server.systems[request.SessionId] = system
	}
	return system
}

func buildSystem(workingDir string, applicationPath string, sessionId string, outputDir string, conn *websocket.Conn) *global.System {
	absApplicationPath := common.AbsolutePath(workingDir, applicationPath)

	systemLog := common.NewSystemLog()

	outputDir, _ = filepath.Abs(outputDir)
	outputDir = filepath.Clean(outputDir)

	system := global.NewSystem(absApplicationPath, sessionId, outputDir, systemLog, conn)

	return system
}
