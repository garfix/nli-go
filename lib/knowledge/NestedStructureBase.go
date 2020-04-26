package knowledge

import "nli-go/lib/mentalese"

// nested query structures (quant, or)
type NestedStructureBase interface {
	KnowledgeBase
	SolveNestedStructure(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings
}
