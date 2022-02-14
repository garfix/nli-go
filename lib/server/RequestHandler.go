package server

import (
	"encoding/json"
	"net"
	"nli-go/lib/global"
	"nli-go/lib/mentalese"
)

type RequestHandler struct {
	conn net.Conn
}

func (handler *RequestHandler) handleMessage(system *global.System, inMessage mentalese.RelationSet) {

	log := system.GetLog()
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
	result := system.Query(query)

	handler.conn.Write([]byte(result.ToJson()))
	handler.conn.Close()
}
