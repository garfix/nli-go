package knowledge

import (
	"nli-go/lib/mentalese"
	"strconv"
)

type SystemPredicateBase struct {
	rules []mentalese.Rule
}

func NewSystemPredicateBase() *SystemPredicateBase {
	return &SystemPredicateBase{}
}

func (base *SystemPredicateBase) Bind(goal mentalese.Relation, bindings []mentalese.Binding) ([]mentalese.Binding, bool) {

	newBindings := []mentalese.Binding{}
	ok := true

	if goal.Predicate == "numberOf" {

		subjectVariable := goal.Arguments[0].TermValue
		numberOfVariable := goal.Arguments[1].TermValue
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

		numberOf := len(differentValues)

		newBindings = []mentalese.Binding{}
		for _, binding := range bindings {
			newBinding := binding.Copy()
			newBinding[numberOfVariable] = mentalese.Term{ TermType:mentalese.Term_number, TermValue: strconv.Itoa(numberOf) }
			newBindings = append(newBindings, newBinding)
		}

	} else {

		ok = false

	}

	return newBindings, ok
}