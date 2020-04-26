package nested

import (
	"nli-go/lib/mentalese"
)

func (base *SystemNestedStructureBase) Call(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	child := relation.Arguments[0].TermValueRelationSet

	newBindings := base.solver.SolveRelationSet(child, mentalese.Bindings{ binding })

	return newBindings
}
