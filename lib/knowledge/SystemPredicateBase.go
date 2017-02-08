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
	aggregate := mentalese.Term{}

	subjectVariable := goal.Arguments[0].TermValue
	numberOfVariable := goal.Arguments[1].TermValue

	if goal.Predicate == "numberOf" {

		differentValues := base.getDifferentValues(bindings, subjectVariable)
		aggregate = mentalese.Term{ TermType:mentalese.Term_number, TermValue: strconv.Itoa(len(differentValues)) }

	} else {
		ok = false
	}

	if ok {
		newBindings = []mentalese.Binding{}
		for _, binding := range bindings {
			newBinding := binding.Copy()
			newBinding[numberOfVariable] = aggregate
			newBindings = append(newBindings, newBinding)
		}
	}

	return newBindings, ok
}

func (base *SystemPredicateBase) getDifferentValues(bindings []mentalese.Binding, subjectVariable string) []mentalese.Term {

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