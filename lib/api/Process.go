package api

import "nli-go/lib/mentalese"

type Process interface {
	GetType() string
	GetChannel() chan mentalese.Request
	Advance()
}
