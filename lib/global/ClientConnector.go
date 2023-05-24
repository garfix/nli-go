package global

import (
	"nli-go/lib/central"
	"nli-go/lib/mentalese"

	"golang.org/x/net/websocket"
)

type ClientConnector struct {
	conn   *websocket.Conn
	system *System
}

func (c *ClientConnector) SendToProcess(processType string, message mentalese.RelationSet) {
	first := message[0]
	// todo: remove this!
	if first.Predicate == mentalese.PredicateRespond {
		c.system.processRunner.StartProcess(
			central.LANGUAGE_PROCESS,
			message,
			mentalese.NewBinding())
	} else {
		c.system.processRunner.RunRelationSet(processType, message)
	}
}

func (c *ClientConnector) SendToClient(processType string, messageType string, message interface{}) {
	println("client connector sending! " + messageType)
	response := mentalese.Response{
		ProcessType: processType,
		MessageType: messageType,
		Success:     true,
		ErrorLines:  []string{},
		Productions: []string{},
		Message:     message,
	}
	websocket.JSON.Send(c.conn, response)
}
