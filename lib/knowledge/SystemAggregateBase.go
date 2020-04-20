package knowledge

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
)

type SystemAggregateBase struct {
	KnowledgeBaseCore
	rules []mentalese.Rule
	log   *common.SystemLog
}

func NewSystemAggregateBase(name string, log *common.SystemLog) *SystemAggregateBase {
	return &SystemAggregateBase{KnowledgeBaseCore: KnowledgeBaseCore{ Name: name }, log: log}
}

func (base *SystemAggregateBase) HandlesPredicate(predicate string) bool {
	predicates := []string{"number_of", "first", "exists"}

	for _, p := range predicates {
		if p == predicate {
			return true
		}
	}
	return false
}

func (base *SystemAggregateBase) Bind(goal mentalese.Relation, bindings mentalese.Bindings) (mentalese.Bindings, bool) {

	base.log.StartDebug("SystemAggregateBase BindSingle", goal, bindings)

	newBindings := mentalese.Bindings{}
	found := true
	aggregate := mentalese.Term{}

	aggregateArgument := mentalese.NewAnonymousVariable()
	aggregateVariable := ""

	if len(goal.Arguments) != 0 {
		aggregateArgument = goal.Arguments[0]
		aggregateVariable = aggregateArgument.TermValue
	}

	if goal.Predicate == "number_of" {

		// exception
		aggregateArgument = goal.Arguments[1]
		aggregateVariable = aggregateArgument.TermValue

		sourceVariable := goal.Arguments[0].TermValue

		differentValues := base.getDifferentValues(bindings, sourceVariable)
		aggregate = mentalese.Term{TermType: mentalese.TermNumber, TermValue: strconv.Itoa(len(differentValues))}

	} else if goal.Predicate == "first" {

		for _, binding := range bindings {

			alreadyPresent := false

			for _, newBinding := range newBindings {

				allFound := true

				for _, argument := range goal.Arguments {

					_, found := newBinding[argument.TermValue]
					if !found {
						allFound = false
					}
				}

				if allFound {
					alreadyPresent = true
				}

			}

			if !alreadyPresent {
				newBindings = append(newBindings, binding)
			}

		}

// todo the first values must be applied to all bindings; do not just throw them away!

		return newBindings, true


	// check if there are still any bindings
	} else if goal.Predicate == "exists" {

		if len(bindings) == 0 {
			found = false
		}

	} else {
		found = false
	}

	if found {
		newBindings = mentalese.Bindings{}

		if aggregateVariable == "" {

			newBindings = bindings

			// number_of(E1, 4)
		} else if aggregateArgument.IsNumber() {

			if aggregateArgument.TermValue == aggregate.TermValue {
				newBindings = bindings
			}

			// number_of(E1, N)
		} else {

			if len(bindings) > 0 {

				for _, binding := range bindings {
					newBinding := binding.Copy()
					newBinding[aggregateVariable] = aggregate
					newBindings = append(newBindings, newBinding)
				}
			} else {

				newBinding := mentalese.Binding{}
				newBinding[aggregateVariable] = aggregate
				newBindings = append(newBindings, newBinding)

			}
		}
	}

	base.log.EndDebug("SystemAggregateBase BindSingle", newBindings, found)

	return newBindings, found
}

func (base *SystemAggregateBase) getDifferentValues(bindings mentalese.Bindings, subjectVariable string) []mentalese.Term {

	differentValues := []mentalese.Term{}

	for _, binding := range bindings {
		value, found := binding[subjectVariable]
		if found {
			preExists := false
			for _, differentValue := range differentValues {
				if differentValue.Equals(value) {
					preExists = true
				}
			}
			if !preExists {
				differentValues = append(differentValues, value)
			}
		}
	}

	return differentValues
}
