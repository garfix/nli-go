package api

type ClientConnector interface {
	SendToClient(processType string, messageType string, message interface{})
}
