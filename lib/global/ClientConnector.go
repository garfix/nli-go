package global

import (
	"nli-go/lib/mentalese"

	"golang.org/x/net/websocket"
)

type ClientConnector struct {
	conn   *websocket.Conn
	system *System
}

// func CreatClientConnector(conn *websocket.Conn, processRunner *central.ProcessRunner) *ClientConnector {
// 	return &ClientConnector{
// 		conn:          conn,
// 		processRunner: processRunner,
// 	}
// }

func (c *ClientConnector) SendToProcess(processType string, message mentalese.RelationSet) {
	first := message[0]
	if first.Predicate == mentalese.PredicateRespond {

	}
}

func (c *ClientConnector) SendToClient(processType string, message mentalese.RelationSet) {
	response := mentalese.Response{
		Success:     true,
		ErrorLines:  []string{},
		Productions: []string{},
		Message:     message,
	}
	websocket.JSON.Send(c.conn, response)
}
