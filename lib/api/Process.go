package api

import "nli-go/lib/mentalese"

type Process interface {
	Advance()
	SetWaitingFor(set mentalese.RelationSet)
	GetWaitingFor() mentalese.RelationSet
}