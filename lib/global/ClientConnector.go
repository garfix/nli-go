package global

import (
	"nli-go/lib/mentalese"

	"golang.org/x/net/websocket"
)

type ClientConnector struct {
	conn   *websocket.Conn
	system *System
}

func (c *ClientConnector) SendToClient(resource string, messageType string, message interface{}) {
	response := mentalese.Response{
		Resource:    resource,
		MessageType: messageType,
		Success:     true,
		ErrorLines:  []string{},
		Productions: []string{},
		Message:     message,
	}
	// fmt.Printf("%v sent:     %s\n", &c.conn, messageType)
	websocket.JSON.Send(c.conn, response)
}
