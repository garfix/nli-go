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

func (base *SystemAggregateBase) validate(input mentalese.Relation, format string) bool {

	expectedLength := len(format)

	for i, c := range format {
		if i >= len(input.Arguments) {
			base.log.AddError("Function '" + input.Predicate + "' expects at least " + strconv.Itoa(expectedLength) + " arguments")
			return false
		}
		arg := input.Arguments[i]
		if c == 'v' && !arg.IsVariable() {
			base.log.AddError("Function '" + input.Predicate + "' expects argument " + strconv.Itoa(i + 1) + " to be an unbound variable")
			return false
		}
		if c == 's' && !arg.IsString() {
			base.log.AddError("Function '" + input.Predicate + "' expects argument " + strconv.Itoa(i + 1) + " to be a string")
			return false
		}
		if c == 'i' && !arg.IsNumber() {
			//			base.log.AddError("Function '" + input.Predicate + "' expects argument " + strconv.Itoa(i + 1) + " to be a number")
			return false
		}
		if c == 'S' {
			expectedLength = len(input.Arguments)
			for j := i; j < len(input.Arguments); j++ {
				arg = input.Arguments[j]
				if !arg.IsString() {
					base.log.AddError("Function '" + input.Predicate + "' expects argument " + strconv.Itoa(j + 1) + " to be a string")
					return false
				}
			}
			break
		}
		if c == 'V' {
			expectedLength = len(input.Arguments)
			for j := i; j < len(input.Arguments); j++ {
				arg = input.Arguments[j]
				if !arg.IsVariable() {
					base.log.AddError("Function '" + input.Predicate + "' expects argument " + strconv.Itoa(j + 1) + " to be an unbound variable")
					return false
				}
			}
			break
		}
	}

	if expectedLength != len(input.Arguments) {
		base.log.AddError("Function '" + input.Predicate + "' expects " + strconv.Itoa(expectedLength) + " arguments")
		return false
	}

	return true
}

func (base *SystemAggregateBase) numberOf(input mentalese.Relation, bindings mentalese.Bindings) mentalese.Bindings {

	if !base.validate(input, "--") {
		return mentalese.Bindings{}
	}

	subjectVariable := input.Arguments[0].TermValue
	numberArgumentValue := input.Arguments[1].TermValue
	number :=  bindings.GetDistinctValueCount(subjectVariable)

	newBindings := mentalese.Bindings{}

	if input.Arguments[1].IsVariable() {
		for _, binding := range bindings {
			newBinding := binding.Copy()
			newBinding[numberArgumentValue] = mentalese.NewString(strconv.Itoa(number))
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

	if !base.validate(input, "v") {
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
			newBinding[subjectVariable] = mentalese.NewString(distinctValues[0])
			newBindings = append(newBindings, newBinding)
		}
	}

	return newBindings
}

func (base *SystemAggregateBase) exists(input mentalese.Relation, bindings mentalese.Bindings) mentalese.Bindings {

	if !base.validate(input, "") {
		return mentalese.Bindings{}
	}

	return bindings
}

func (base *SystemAggregateBase) Execute(input mentalese.Relation, bindings mentalese.Bindings) (mentalese.Bindings, bool) {

	newBindings := bindings
	found := true

	switch input.Predicate {
	case "number_of":
		newBindings = base.numberOf(input, bindings)
	case "first":
		newBindings = base.first(input, bindings)
	case "exists":
		newBindings = base.exists(input, bindings)
	default:
		found = false
	}

	return newBindings, found
}
