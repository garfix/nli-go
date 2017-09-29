package knowledge

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
)

type SystemAggregateBase struct {
	rules []mentalese.Rule
	log   *common.SystemLog
}

func NewSystemAggregateBase(log *common.SystemLog) *SystemAggregateBase {
	return &SystemAggregateBase{log: log}
}

func (ruleBase *SystemAggregateBase) GetMatchingGroups(set mentalese.RelationSet, knowledgeBaseIndex int) []RelationGroup {

	matchingGroups := []RelationGroup{}
	predicates := []string{"number_of", "exists"}

	for _, setRelation := range set {
		for _, predicate:= range predicates {
			if predicate == setRelation.Predicate {
// TODO calculate real cost
				matchingGroups = append(matchingGroups, RelationGroup{mentalese.RelationSet{setRelation}, knowledgeBaseIndex, worst_cost})
				break
			}
		}
	}

	return matchingGroups
}

func (base *SystemAggregateBase) Bind(goal mentalese.Relation, bindings []mentalese.Binding) ([]mentalese.Binding, bool) {

	base.log.StartDebug("SystemAggregateBase Bind", goal, bindings)

	newBindings := []mentalese.Binding{}
	found := true
	aggregate := mentalese.Term{}

	resultArgument := goal.Arguments[0]
	resultVariable := resultArgument.TermValue

	if goal.Predicate == "number_of" {

		subjectVariable := goal.Arguments[1].TermValue

		differentValues := base.getDifferentValues(bindings, subjectVariable)
		aggregate = mentalese.Term{TermType: mentalese.Term_number, TermValue: strconv.Itoa(len(differentValues))}

	} else if goal.Predicate == "exists" {

		subjectVariable := goal.Arguments[1].TermValue

		differentValues := base.getDifferentValues(bindings, subjectVariable)
		val := "false"
		if len(differentValues) > 0 {
			val = "true"
		}
		aggregate = mentalese.Term{TermType: mentalese.Term_predicateAtom, TermValue: val}

	} else {
		found = false
	}

	if found {
		newBindings = []mentalese.Binding{}

		// number_of(4, E1)
		if resultArgument.IsNumber() {

			if resultArgument.TermValue == aggregate.TermValue {
				newBindings = bindings
			}

			// number_of(N, E1)
		} else {

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
	}

	base.log.EndDebug("SystemAggregateBase Bind", newBindings, found)

	return newBindings, found
}

func (base *SystemAggregateBase) getDifferentValues(bindings []mentalese.Binding, subjectVariable string) []mentalese.Term {

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
