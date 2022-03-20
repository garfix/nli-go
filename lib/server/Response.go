package server

import "nli-go/lib/mentalese"

type Response struct {
	Success     bool
	ErrorLines  []string
	Productions []string
	Message     mentalese.RelationSet
}

type ResponseAnswer struct {
	Success     bool
	ErrorLines  []string
	Productions []string
	Answer      string
}
