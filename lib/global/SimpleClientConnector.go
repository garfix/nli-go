package global

import "nli-go/lib/mentalese"

type SimpleClientConnector struct {
}

func CreateSimpleClientConnector() *SimpleClientConnector {
	return &SimpleClientConnector{}
}

func (c *SimpleClientConnector) SendToProcess(processType string, message mentalese.RelationSet) {

}
func (c *SimpleClientConnector) SendToClient(processType string, message mentalese.RelationSet) {

}
