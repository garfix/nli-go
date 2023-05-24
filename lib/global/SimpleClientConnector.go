package global

import "nli-go/lib/mentalese"

type SimpleClientConnector struct {
}

func CreateSimpleClientConnector() *SimpleClientConnector {
	return &SimpleClientConnector{}
}

func (c *SimpleClientConnector) SendToClient(processType string, messageType string, message mentalese.RelationSet) {

}
