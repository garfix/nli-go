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

func (ruleBase *SystemAggregateBase) GetMatchingGroups(set mentalese.RelationSet, keyCabinet *mentalese.KeyCabinet) []RelationGroup {

	matchingGroups := []RelationGroup{}
	predicates := []string{"number_of", "exists", "first"}

	for _, setRelation := range set {
		for _, predicate:= range predicates {
			if predicate == setRelation.Predicate {
// TODO calculate real cost
				matchingGroups = append(matchingGroups, RelationGroup{mentalese.RelationSet{setRelation}, ruleBase.Name, worst_cost})
				break
			}
		}
	}

	return matchingGroups
}

func (base *SystemAggregateBase) Bind(goal mentalese.Relation, bindings mentalese.Bindings) (mentalese.Bindings, bool) {

	base.log.StartDebug("SystemAggregateBase BindSingle", goal, bindings)

	newBindings := mentalese.Bindings{}
	found := true
	aggregate := mentalese.Term{}

	if len(goal.Arguments) == 0 {
		return newBindings, false
	}

	aggregateArgument := goal.Arguments[0]
	aggregateVariable := aggregateArgument.TermValue

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

	} else if goal.Predicate == "exists" {

		subjectVariable := goal.Arguments[1].TermValue

		differentValues := base.getDifferentValues(bindings, subjectVariable)
		val := "false"
		if len(differentValues) > 0 {
			val = "true"
		}
		aggregate = mentalese.Term{TermType: mentalese.TermPredicateAtom, TermValue: val}

	} else {
		found = false
	}

	if found {
		newBindings = mentalese.Bindings{}

		// number_of(E1, 4)
		if aggregateArgument.IsNumber() {

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
