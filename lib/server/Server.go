package server

import (
	"context"
	"fmt"
	"net/http"
	"nli-go/lib/api"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/global"
	"nli-go/lib/mentalese"
	"nli-go/resources/blocks"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	httpServer *http.Server
	baseDir    string
	workdDir   string
}

func NewServer(port string, appDir string, workdDir string) *Server {
	server := &http.Server{Addr: ":" + port}
	return &Server{
		httpServer: server,
		baseDir:    appDir,
		workdDir:   workdDir,
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

		server.HandleSingleRequest(conn, request, systems)
	}
}

func (server *Server) HandleSingleRequest(conn *websocket.Conn, request mentalese.Request, systems map[string]api.System) {
	defer func() {
		if err := recover(); err != nil {
			message := fmt.Sprintf("%v", err)
			println(message)

			response := mentalese.Response{
				Resource:    central.NO_RESOURCE,
				MessageType: mentalese.MessageError,
				Message:     message,
			}

			websocket.JSON.Send(conn, response)
		}
	}()

	server.logRequest(request)

	if request.MessageType == mentalese.MessageReset {
		delete(systems, request.System)
		response := mentalese.Response{
			Resource:    central.NO_RESOURCE,
			MessageType: mentalese.MessageAcknowledge,
			Message:     "",
		}

		websocket.JSON.Send(conn, response)
	} else if request.MessageType == mentalese.MessageSendLog {
		system := server.getSystem(conn, request, &systems)
		response := mentalese.Response{
			Resource:    central.NO_RESOURCE,
			MessageType: mentalese.MessageAcknowledge,
			Message:     strings.Join(system.GetLog().GetErrors(), "\n"),
		}
		websocket.JSON.Send(conn, response) //
	} else if request.MessageType == mentalese.MessageDebug {
		system := server.getSystem(conn, request, &systems)
		if request.Message == "on" {
			system.GetLog().SetDebug(true)
		} else {
			system.GetLog().SetDebug(false)
		}

	} else {
		system := server.getSystem(conn, request, &systems)
		system.HandleRequest(request)
	}
}

func (server *Server) logRequest(request mentalese.Request) {

	year, month, _ := time.Now().Date()
	filename := server.workdDir + "/log/" + strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-queries.log"

	if request.Resource == central.RESOURCE_LANGUAGE && request.Message != "" {
		text := fmt.Sprintf("%s %s", time.Now().Local(), request.Message.(string)+"\n")

		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			f.WriteString(text)
		}
		f.Close()
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

	sessionId := common.CreateUuid()
	systemLog := common.NewSystemLog()
	appDir := server.baseDir + "/" + request.System

	system = global.NewSystem(appDir, server.workdDir, sessionId, systemLog, conn)
	if request.System == "blocks" {
		system = blocks.CreateBlocksSystem(system)
	}

	if system.GetLog().IsOk() {
		(*systems)[request.System] = system
	}

	return system
}
