package server

import "nli-go/lib/mentalese"

type Response struct {
	Success     bool
	ErrorLines  []string
	Productions []string
	Message     mentalese.RelationSet
	Bindings    mentalese.BindingSet
}
