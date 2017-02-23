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

	} else if goal.Predicate == "and" {

		entityVar := goal.Arguments[0].TermValue
		conjVar := goal.Arguments[1].TermValue
		leftVar := goal.Arguments[2].TermValue
		rightVar := goal.Arguments[3].TermValue

		if len(bindings) > 1 {

			newBindings = []mentalese.Binding{}

			binding1 := bindings[0].Copy()
			left := binding1[entityVar]

			binding2 := bindings[1].Copy()
			right := binding2[entityVar]
			conj := mentalese.Term{ TermType:mentalese.Term_predicateAtom, TermValue: "q1" }

			binding1[conjVar] = conj
			binding1[leftVar] = left
			binding1[rightVar] = right

			binding2[conjVar] = conj
			binding2[leftVar] = left
			binding2[rightVar] = right

			newBindings = append(newBindings, binding1)
			newBindings = append(newBindings, binding2)

			left = conj

			for i := 2; i < len(bindings); i++ {

				binding := bindings[i].Copy()
				right = binding[entityVar]

				conj = mentalese.Term{ TermType:mentalese.Term_predicateAtom, TermValue: "q" + strconv.Itoa(i)}
				binding[conjVar] = conj
				binding[leftVar] = left
				binding[rightVar] = right
				newBindings = append(newBindings, binding)

				left = conj
			}

			// reverse the relations to have the topmost conjunction first
			newBindings2 := []mentalese.Binding{}
			for i := len(newBindings) - 1; i >= 0; i-- {
				newBindings2 = append(newBindings2, newBindings[i])
			}
			newBindings = newBindings2

			ok = false
		}

	} else {
		ok = false
	}

	if ok {
		newBindings = []mentalese.Binding{}

		if len(bindings) > 0 {

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