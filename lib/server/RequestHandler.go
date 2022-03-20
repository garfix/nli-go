package server

import (
	"encoding/json"
	"fmt"
	"net"
	"nli-go/lib/global"
	"nli-go/lib/mentalese"
	"runtime/debug"
	"strings"
)

type RequestHandler struct {
	conn net.Conn
}

func (handler *RequestHandler) panicHandler() {
	if r := recover(); r != nil {
		errorString := fmt.Sprintf("%s\n%s", r, debug.Stack())
		response := Response{
			Success:    false,
			ErrorLines: strings.Split(errorString, "\n"),
		}
		responseJSON, _ := json.Marshal(response)
		handler.conn.Write(responseJSON)
		handler.conn.Close()
	}
}

func (handler *RequestHandler) handleSend(system *global.System, inMessage mentalese.RelationSet) {

	defer handler.panicHandler()

	log := system.GetLog()
	log.Clear()

	outMMessage := system.SendAndWaitForResponse(inMessage)

	response := Response{
		Success:     log.IsOk(),
		ErrorLines:  log.GetErrors(),
		Productions: log.GetProductions(),
		Message:     outMMessage,
	}

	responseRaw, _ := json.MarshalIndent(response, "", "    ")
	responseString := string(responseRaw) + "\n"
	handler.conn.Write([]byte(responseString))
	handler.conn.Close()
}

func (handler *RequestHandler) handleQuery(system *global.System, query string) {

	defer handler.panicHandler()

	resultJson := system.Query(query).ToJson()
	handler.conn.Write([]byte(resultJson))
	handler.conn.Close()
}

func (handler *RequestHandler) handleAnswer(system *global.System, query string) {

	defer handler.panicHandler()

	log := system.GetLog()
	result, options := system.Answer(query)

	message := result
	if options.HasOptions() {
		message = options.String()
	}

	response := ResponseAnswer{
		Success:     log.IsOk(),
		ErrorLines:  log.GetErrors(),
		Productions: log.GetProductions(),
		Answer:      message,
	}
	responseJSON, _ := json.Marshal(response)

	handler.conn.Write(responseJSON)
	handler.conn.Close()
}
