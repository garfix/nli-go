package api

import "nli-go/lib/mentalese"

type SimpleMessenger interface {
	SetOutBinding(variable string, value mentalese.Term)
	GetOutBinding() mentalese.Binding
}
