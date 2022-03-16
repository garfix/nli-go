package api

import "nli-go/lib/mentalese"

type Process interface {
	Advance()
	SetMutableVariable(variable string, value mentalese.Term)
}
