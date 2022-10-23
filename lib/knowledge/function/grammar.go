package function

import (
	"nli-go/lib/api"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
)

func (base *SystemSolverFunctionBase) eventReference(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !knowledge.Validate(bound, "v", base.log) {
		return mentalese.NewBindingSet()
	}

	return mentalese.InitBindingSet(binding)
}
