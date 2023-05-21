package blocks

import (
	"nli-go/lib/api"
	"nli-go/lib/mentalese"
)

type BlocksSystem struct {
	base api.System
}

func CreateBlocksSystem(base api.System) *BlocksSystem {
	return &BlocksSystem{
		base: base,
	}
}

func (system *BlocksSystem) HandleRequest(request mentalese.Request) {
	system.base.HandleRequest(request)
}
