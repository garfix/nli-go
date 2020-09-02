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
	predicates := []string{mentalese.PredicateNumberOf, mentalese.PredicateFirst, mentalese.PredicateExists}

	for _, p := range predicates {
		if p == predicate {
			return true
		}
	}
	return false
}

func (base *SystemAggregateBase) numberOf(input mentalese.Relation, bindings mentalese.Bindings) mentalese.Bindings {

	if !Validate(input, "--", base.log) {
		return mentalese.Bindings{}
	}

	subjectVariable := input.Arguments[0].TermValue
	numberArgumentValue := input.Arguments[1].TermValue
	number :=  bindings.GetDistinctValueCount(subjectVariable)

	newBindings := mentalese.Bindings{}

	if input.Arguments[1].IsVariable() {
		for _, binding := range bindings {
			newBinding := binding.Copy()
			newBinding[numberArgumentValue] = mentalese.NewTermString(strconv.Itoa(number))
			newBindings = append(newBindings, newBinding)
		}
	} else {
		assertedNumber, err := strconv.Atoi(numberArgumentValue)
		if err != nil {
			base.log.AddError("The second argument of number_of() needs to be an integer")
			newBindings = mentalese.Bindings{}
		} else {
			if number == assertedNumber {
				newBindings = bindings
			} else {
				newBindings = mentalese.Bindings{}
			}
		}
	}

	return newBindings
}

func (base *SystemAggregateBase) first(input mentalese.Relation, bindings mentalese.Bindings) mentalese.Bindings {

	if !Validate(input, "v", base.log) {
		return mentalese.Bindings{}
	}

	subjectVariable := input.Arguments[0].TermValue
	distinctValues := bindings.GetDistinctValues(subjectVariable)

	newBindings := mentalese.Bindings{}
	if len(distinctValues) == 0 {
		newBindings = bindings
	} else {
		for _, binding := range bindings {
			newBinding := binding.Copy()
			newBinding[subjectVariable] = mentalese.NewTermString(distinctValues[0])
			newBindings = append(newBindings, newBinding)
		}
	}

	return newBindings
}

func (base *SystemAggregateBase) exists(input mentalese.Relation, bindings mentalese.Bindings) mentalese.Bindings {

	if !Validate(input, "", base.log) {
		return mentalese.Bindings{}
	}

	return bindings
}

func (base *SystemAggregateBase) Execute(input mentalese.Relation, bindings mentalese.Bindings) (mentalese.Bindings, bool) {

	newBindings := bindings
	found := true

	switch input.Predicate {
	case mentalese.PredicateNumberOf:
		newBindings = base.numberOf(input, bindings)
	case mentalese.PredicateFirst:
		newBindings = base.first(input, bindings)
	case mentalese.PredicateExists:
		newBindings = base.exists(input, bindings)
	default:
		found = false
	}

	return newBindings, found
}
