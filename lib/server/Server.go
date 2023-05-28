package server

import (
	"context"
	"net/http"
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/global"
	"nli-go/lib/mentalese"
	"nli-go/resources/blocks"
	"path/filepath"

	"golang.org/x/net/websocket"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string) *Server {
	server := &http.Server{Addr: ":" + port}
	return &Server{
		httpServer: server,
	}
}

func (server *Server) Run() {

	http.Handle("/", websocket.Handler(server.HandleSingleConnection))

	println("NLI-GO server listening on " + server.httpServer.Addr + "\n")

	err := server.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic("ListenAndServe: " + err.Error())
	}

	println("Server closed")
}

func (server *Server) RunInBackground() {
	go server.Run()
}

func (server *Server) HandleSingleConnection(conn *websocket.Conn) {

	systems := map[string]api.System{}

	for {
		request := mentalese.Request{}

		err := websocket.JSON.Receive(conn, &request)
		if err != nil {
			println(err.Error())
			break
		}

		system := server.getSystem(conn, request, &systems)
		system.HandleRequest(request)
	}
}

func (server *Server) Close() {
	server.httpServer.Shutdown(context.TODO())
}

func (server *Server) getSystem(conn *websocket.Conn, request mentalese.Request, systems *map[string]api.System) api.System {

	system, found := (*systems)[request.System]
	if found {
		return system
	}

	applicationDir := common.Dir() + "/../../resources/" + request.System
	workDir := common.Dir() + "/../../var"

	sessionId := common.CreateUuid()

	system = buildSystem(workDir, applicationDir, sessionId, workDir, conn)
	if request.System == "blocks" {
		system = blocks.CreateBlocksSystem(system)
	}

	(*systems)[request.System] = system

	return system
}

func buildSystem(workingDir string, applicationPath string, sessionId string, outputDir string, conn *websocket.Conn) api.System {
	absApplicationPath := common.AbsolutePath(workingDir, applicationPath)

	systemLog := common.NewSystemLog()

	outputDir, _ = filepath.Abs(outputDir)
	outputDir = filepath.Clean(outputDir)

	system := global.NewSystem(absApplicationPath, sessionId, outputDir, systemLog, conn)

	return system
}
