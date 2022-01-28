package server

import "nli-go/lib/mentalese"

type Request struct {
	SessionId      string
	ApplicationDir string
	WorkDir        string
	Command        string
	Message        mentalese.RelationSet
	Query          string
}
