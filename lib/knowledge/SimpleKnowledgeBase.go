package knowledge

import "nli-go/lib/mentalese"

type SimpleKnowledgeBase interface {
	// goal e.g. father(X, 'john')
	// return subgoalSets e.g. {
	//    { male(X), parent(X, 'john') },
	//    { child('john', X), male(X) }
	// }
	// return bindings e.g. {
	//    { X='Jack' },
	// }
	// Note: bindings are linked to subgoalSets, one on one; but usually just one of the arrays is used
	Bind(goal mentalese.SimpleRelation) ([][]mentalese.SimpleRelation, []mentalese.SimpleBinding)
}
