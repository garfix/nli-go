package knowledge

import (
	"nli-go/lib/mentalese"
	"strconv"
	"strings"
	"time"
)

type SystemFunctionBase struct {
	KnowledgeBaseCore
}

func NewSystemFunctionBase(name string) *SystemFunctionBase {
	return &SystemFunctionBase{ KnowledgeBaseCore{ Name: name } }
}

func (base *SystemFunctionBase) HandlesPredicate(predicate string) bool {
	predicates := []string{"split", "join", "concat", "greater_than", "less_than", "equals", "add", "date_today", "date_subtract_years"}

	for _, p := range predicates {
		if p == predicate {
			return true
		}
	}
	return false
}

func (base *SystemFunctionBase) Execute(input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	newBinding := binding
	found := false

	if input.Predicate == "split" {

		sourceValue, sourceFound := binding[input.Arguments[0].TermValue]
		if sourceFound {

			parts := strings.Split(sourceValue.TermValue, input.Arguments[1].TermValue)
			newBinding = binding.Copy()

			for i, argument := range input.Arguments[2:] {
				variable := argument.TermValue
				newTerm := mentalese.Term{}
				newTerm.TermType = mentalese.TermStringConstant
				newTerm.TermValue = parts[i]
				newBinding[variable] = newTerm
			}

			found = true
		}
	}

	// join(result, separator, arg1, arg2, ...)
	if input.Predicate == "join" {

		sep := ""
		result := mentalese.Term{}

		for _, argument := range input.Arguments[2:] {
			variable := argument.TermValue
			variableValue, variableFound := binding[variable]
			if variableFound {
				result.TermType = mentalese.TermStringConstant
				result.TermValue += sep + variableValue.TermValue
				sep = input.Arguments[1].TermValue
			}
		}

		newBinding = binding.Copy()
		newBinding[input.Arguments[0].TermValue] = result

		found = true

	}

	// concat(result, arg1, arg2, ...)
	if input.Predicate == "concat" {

		result := mentalese.Term{}
		result.TermType = mentalese.TermStringConstant

		for _, argument := range input.Arguments[1:] {
			if (argument.IsVariable()) {
				variable := argument.TermValue
				variableValue, variableFound := binding[variable]
				if variableFound {
					result.TermValue += variableValue.TermValue
				}
			} else {
				result.TermValue += argument.TermValue
			}
		}

		newBinding = binding.Copy()
		newBinding[input.Arguments[0].TermValue] = result

		found = true

	}

	// greater_than(arg1, arg2)
	if input.Predicate == "greater_than" {

		arg1 := input.Arguments[0]
		arg2 := input.Arguments[1]

		int1, _ := strconv.Atoi(arg1.TermValue)
		int2, _ := strconv.Atoi(arg2.TermValue)

		value, foundInBinding := binding[input.Arguments[0].TermValue]
		if foundInBinding {
			int1, _ = strconv.Atoi(value.TermValue)
		}

		value, foundInBinding = binding[input.Arguments[1].TermValue]
		if foundInBinding {
			int2, _ = strconv.Atoi(value.TermValue)
		}

		if int1 > int2 {
			found = true
		} else {
			found = false
		}
	}

	// less_than(arg1, arg2)
	if input.Predicate == "less_than" {

		arg1 := input.Arguments[0]
		arg2 := input.Arguments[1]

		int1, _ := strconv.Atoi(arg1.TermValue)
		int2, _ := strconv.Atoi(arg2.TermValue)

		value, foundInBinding := binding[input.Arguments[0].TermValue]
		if foundInBinding {
			int1, _ = strconv.Atoi(value.TermValue)
		}

		value, foundInBinding = binding[input.Arguments[1].TermValue]
		if foundInBinding {
			int2, _ = strconv.Atoi(value.TermValue)
		}

		if int1 < int2 {
			found = true
		} else {
			found = false
		}
	}

	// equals(arg1, arg2)
	if input.Predicate == "equals" {

		arg1 := input.Arguments[0]
		arg2 := input.Arguments[1]

		int1, _ := strconv.Atoi(arg1.TermValue)
		int2, _ := strconv.Atoi(arg2.TermValue)

		value, foundInBinding := binding[input.Arguments[0].TermValue]
		if foundInBinding {
			int1, _ = strconv.Atoi(value.TermValue)
		}

		value, foundInBinding = binding[input.Arguments[1].TermValue]
		if foundInBinding {
			int2, _ = strconv.Atoi(value.TermValue)
		}

		if int1 == int2 {
			found = true
		} else {
			found = false
		}
	}

	// add(arg1, arg2, sum)
	if input.Predicate == "add" {

		arg1 := input.Arguments[0]
		arg2 := input.Arguments[1]

		int1, _ := strconv.Atoi(arg1.TermValue)
		int2, _ := strconv.Atoi(arg2.TermValue)

		value, foundInBinding := binding[input.Arguments[0].TermValue]
		if foundInBinding {
			int1, _ = strconv.Atoi(value.TermValue)
		}

		value, foundInBinding = binding[input.Arguments[1].TermValue]
		if foundInBinding {
			int2, _ = strconv.Atoi(value.TermValue)
		}

		sum := int1 + int2

		newBinding[input.Arguments[2].TermValue] = mentalese.NewString(strconv.Itoa(sum))

		found = true
	}

	if input.Predicate == "date_today" {

		now := time.Now()
		formatted := now.Format("2006-01-02")

		newBinding[input.Arguments[0].TermValue] = mentalese.NewString(formatted)

		found = true
	}

	if input.Predicate == "date_subtract_years" {

		value1 := input.Arguments[0].Resolve(binding)
		value2 := input.Arguments[1].Resolve(binding)

		date1, err := time.Parse("2006-01-02", value1.TermValue)
		date2, err := time.Parse("2006-01-02", value2.TermValue)
		years := 0.0

		if err == nil {
			duration := date1.Sub(date2)
			hours := duration.Hours()
			totalDays := hours / 24
			years = totalDays / 365
		}

		newBinding[input.Arguments[2].TermValue] = mentalese.NewString(strconv.Itoa(int(years)))

		found = true
	}

	return newBinding, found
}
