package central

import (
	"nli-go/lib/mentalese"
)

func (solver ProblemSolver) SolveNot(notRelation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	//variable := notRelation.Arguments[mentalese.NotVariableIndex].TermValueRelationSet
	scope := notRelation.Arguments[mentalese.NotScopeIndex].TermValueRelationSet

	newBindings := solver.SolveRelationSet(scope, mentalese.Bindings{ binding })
	resultBindings := mentalese.Bindings{}

	if len(newBindings) > 0 {
		resultBindings = mentalese.Bindings{}
	} else {
		resultBindings = mentalese.Bindings{ binding }
	}

	return resultBindings
}
