package nested

import (
	"nli-go/lib/mentalese"
)

func (base *SystemNestedStructureBase) SolveNot(notRelation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	scope := notRelation.Arguments[mentalese.NotScopeIndex].TermValueRelationSet

	newBindings := base.solver.SolveRelationSet(scope, mentalese.Bindings{ binding })
	resultBindings := mentalese.Bindings{}

	if len(newBindings) > 0 {
		resultBindings = mentalese.Bindings{}
	} else {
		resultBindings = mentalese.Bindings{ binding }
	}

	return resultBindings
}
