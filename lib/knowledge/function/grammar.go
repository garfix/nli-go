package function

import (
	"nli-go/lib/api"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
)

func (base *SystemSolverFunctionBase) intent(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !knowledge.Validate(bound, "a*", base.log) {
		return mentalese.NewBindingSet()
	}

	return mentalese.InitBindingSet(binding)
}

func (base *SystemSolverFunctionBase) getSense(messenger api.ProcessMessenger) mentalese.RelationSet {
	term, found := messenger.GetProcessSlot(mentalese.SlotSense)
	if !found {
		base.log.AddError("Slot 'sense' not found")
	}

	return term.TermValueRelationSet
}

// ask the user which of the specified entities he/she means
func (base *SystemSolverFunctionBase) rangeIndexClarification(messenger api.ProcessMessenger) {

	messenger.SetProcessSlot(mentalese.SlotSolutionOutput, mentalese.NewTermString("I don't understand which one you mean"))
}
