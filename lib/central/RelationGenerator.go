package central

import "nli-go/lib/mentalese"

type RelationGenerator interface {
	generate(template mentalese.Relation, bindings mentalese.Bindings) (mentalese.RelationSet, bool)
}
