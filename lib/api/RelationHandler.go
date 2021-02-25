package api

import (
	"nli-go/lib/mentalese"
)

type RelationHandler func(messenger ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet