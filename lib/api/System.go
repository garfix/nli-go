package api

import "nli-go/lib/mentalese"

type System interface {
	HandleRequest(request mentalese.Request)
	RunRelationSet(resource string, relationSet mentalese.RelationSet) mentalese.BindingSet
	RunRelationSetString(resource string, relationSet string) mentalese.BindingSet
	GetClientConnector() ClientConnector
}
