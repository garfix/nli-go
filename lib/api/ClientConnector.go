package api

type ClientConnector interface {
	SendToClient(resource string, messageType string, message interface{})
}
