package api

import "nli-go/lib/mentalese"

type ClientConnector interface {
	SendToProcess(processType string, message mentalese.RelationSet)
	SendToClient(processType string, message mentalese.RelationSet)
}
