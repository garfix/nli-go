package knowledge

import (
	"nli-go/lib/mentalese"
	"strconv"
	"nli-go/lib/common"
)

type SystemPredicateBase struct {
	rules []mentalese.Rule
}

func NewSystemPredicateBase() *SystemPredicateBase {
	return &SystemPredicateBase{}
}

func (base *SystemPredicateBase) Bind(goal mentalese.Relation, bindings []mentalese.Binding) ([]mentalese.Binding, bool) {

	common.LogTree("SystemPredicateBase Bind", goal, bindings)

	newBindings := []mentalese.Binding{}
	ok := true
	aggregate := mentalese.Term{}

	resultVariable := goal.Arguments[0].TermValue

	if goal.Predicate == "numberOf" {

		subjectVariable := goal.Arguments[1].TermValue

		differentValues := base.getDifferentValues(bindings, subjectVariable)
		aggregate = mentalese.Term{ TermType:mentalese.Term_number, TermValue: strconv.Itoa(len(differentValues)) }

	} else if goal.Predicate == "exists" {

		subjectVariable := goal.Arguments[1].TermValue

		differentValues := base.getDifferentValues(bindings, subjectVariable)
		val := "false"
		if len(differentValues) > 0 {
			val = "true"
		}
		aggregate = mentalese.Term{ TermType:mentalese.Term_predicateAtom, TermValue: val }

	} else {
		ok = false
	}

	if ok {
		newBindings = []mentalese.Binding{}

		if len(newBindings) > 0 {

			for _, binding := range bindings {
				newBinding := binding.Copy()
				newBinding[resultVariable] = aggregate
				newBindings = append(newBindings, newBinding)
			}
		} else {

			newBinding := mentalese.Binding{}
			newBinding[resultVariable] = aggregate
			newBindings = append(newBindings, newBinding)

		}
	}

	common.LogTree("SystemPredicateBase Bind", newBindings, ok)

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