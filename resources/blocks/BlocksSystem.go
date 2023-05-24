package blocks

import (
	"nli-go/lib/api"
	"nli-go/lib/central"
	"nli-go/lib/mentalese"
)

const MESSAGE_DESCRIBE = "describe"
const MESSAGE_DESCRIPTION = "description"

type BlocksSystem struct {
	base api.System
}

func CreateBlocksSystem(base api.System) *BlocksSystem {
	return &BlocksSystem{
		base: base,
	}
}

func (system *BlocksSystem) HandleRequest(request mentalese.Request) {
	switch request.MessageType {
	case MESSAGE_DESCRIBE:
		scene := "dom:at(E, X, Z, Y) go:has_sort(E, Type) dom:color(E, Color) dom:size(E, Width, Length, Height)"
		bindings := system.base.RunRelationSetString(central.SIMPLE_PROCESS, scene)
		system.base.GetClientConnector().SendToClient(central.SIMPLE_PROCESS, MESSAGE_DESCRIPTION, bindings.AsSimple())
	default:
		system.base.HandleRequest(request)
	}
}

func (system *BlocksSystem) RunRelationSet(processType string, relationSet mentalese.RelationSet) mentalese.BindingSet {
	return system.base.RunRelationSet(processType, relationSet)
}

func (system *BlocksSystem) RunRelationSetString(processType string, relationSet string) mentalese.BindingSet {
	return system.base.RunRelationSetString(processType, relationSet)
}

func (system *BlocksSystem) GetClientConnector() api.ClientConnector {
	return system.base.GetClientConnector()
}
