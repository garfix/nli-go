package central

import "nli-go/lib/mentalese"

type RelationGenerator interface {
	generate(template mentalese.Relation, bindings []mentalese.Binding) (mentalese.RelationSet, bool)
}
