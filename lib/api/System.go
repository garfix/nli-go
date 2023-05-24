package api

import "nli-go/lib/mentalese"

type System interface {
	HandleRequest(request mentalese.Request)
	RunRelationSet(processType string, relationSet mentalese.RelationSet) mentalese.BindingSet
	RunRelationSetString(processType string, relationSet string) mentalese.BindingSet
	GetClientConnector() ClientConnector
}
