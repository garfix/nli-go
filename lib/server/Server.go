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

		reader := bufio.NewReader(conn)
		requestJson, err := reader.ReadBytes('\n')
		if err != nil {
			conn.Close()
			return
		}

		request := Request{}
		err = json.Unmarshal(requestJson, &request)
		if err != nil {
			conn.Close()
			return
		}

		println(request.SessionId + ": " + request.Command)

		switch request.Command {
		case "send":
			system := server.getSystem(request)
			client := &RequestHandler{conn: conn}
			go client.handleMessage(system, request.Message)

		case "reset":
			_, found := server.systems[request.SessionId]
			if found {
				delete(server.systems, request.SessionId)
			}
			conn.Write([]byte("{\"result\": \"OK\"}"))
			conn.Close()

		case "query":
			system := server.getSystem(request)
			client := &RequestHandler{conn: conn}
			go client.handleQuery(system, request.Query)

		default:
			conn.Write([]byte("Unknown command: " + request.Command))
			conn.Close()
		}
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
