package api

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type System interface {
	GetLog() *common.SystemLog
	HandleRequest(request mentalese.Request)
	RunRelationSet(resource string, relationSet mentalese.RelationSet) mentalese.BindingSet
	RunRelationSetString(resource string, relationSet string) mentalese.BindingSet
	GetClientConnector() ClientConnector
}
