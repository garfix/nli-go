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

func (handler *RequestHandler) handleMessage(system *global.System, message mentalese.RelationSet) {

	result := system.HandleMessage(message)
	log := system.GetLog()

	response := Response{
		Success:     log.IsOk(),
		ErrorLines:  log.GetErrors(),
		Productions: log.GetProductions(),
		Message:     result,
	}

	responseRaw, _ := json.MarshalIndent(response, "", "    ")
	responseString := string(responseRaw) + "\n"
	handler.conn.Write([]byte(responseString))
	handler.conn.Close()
	print(responseString)
}

func (handler *RequestHandler) handleQuery(system *global.System, query string) {
	result := system.Query(query)

	handler.conn.Write([]byte(result.ToJson()))
	handler.conn.Close()
}
