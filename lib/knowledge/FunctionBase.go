package knowledge

import "nli-go/lib/mentalese"

type FunctionBase interface {
	KnowledgeBase
	Execute(input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool, bool)
}
