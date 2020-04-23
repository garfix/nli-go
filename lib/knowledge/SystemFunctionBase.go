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
	predicates := []string{"split", "join", "concat", "greater_than", "less_than", "equals", "not_equals", "add", "subtract", "date_today", "date_subtract_years"}

	for _, p := range predicates {
		if p == predicate {
			return true
		}
	}
	return false
}

func (base *SystemFunctionBase) validate(input mentalese.Relation, format string) bool {

	for i, c := range format {
		if i >= len(input.Arguments) {
			return false
		}
		arg := input.Arguments[i]
		if c == 'v' && !arg.IsVariable() {
			return false
		}
		if c == 's' && !arg.IsString() {
			return false
		}
		if c == 'i' && !arg.IsNumber() {
			return false
		}
		if c == 'S' {
			for j := i; j < len(input.Arguments); j++ {
				arg = input.Arguments[j]
				if !arg.IsString() {
					return false
				}
			}
			break
		}
		if c == 'V' {
			for j := i; j < len(input.Arguments); j++ {
				arg = input.Arguments[j]
				if !arg.IsVariable() {
					return false
				}
			}
			break
		}
	}

	return true
}

func (base *SystemFunctionBase) split(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !base.validate(bound, "ssV") {
		return mentalese.Binding{}
	}

	newBinding := binding.Copy()
	parts := strings.Split(bound.Arguments[0].TermValue, bound.Arguments[1].TermValue)

	for i, argument := range bound.Arguments[2:] {
		newBinding[argument.TermValue] = mentalese.NewString(parts[i])
	}

	return newBinding
}

func (base *SystemFunctionBase) join(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !base.validate(bound, "vsS") {
		return mentalese.Binding{}
	}

	newBinding := binding.Copy()
	sep := ""
	result := ""
	for _, argument := range bound.Arguments[2:] {
		result += sep + argument.TermValue
		sep = input.Arguments[1].TermValue
	}

	newBinding[input.Arguments[0].TermValue] = mentalese.NewString(result)

	return newBinding
}

func (base *SystemFunctionBase) concat(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !base.validate(bound, "vS") {
		return mentalese.Binding{}
	}

	newBinding := binding.Copy()
	result := ""
	for _, argument := range bound.Arguments[1:] {
		result += argument.TermValue
	}

	newBinding[input.Arguments[0].TermValue] = mentalese.NewString(result)

	return newBinding
}

func (base *SystemFunctionBase) greaterThan(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !base.validate(bound, "ii") {
		return mentalese.Binding{}
	}

	int1, _ := strconv.Atoi(bound.Arguments[0].TermValue)
	int2, _ := strconv.Atoi(bound.Arguments[1].TermValue)

	if int1 > int2 {
		return binding
	} else {
		return mentalese.Binding{}
	}
}

func (base *SystemFunctionBase) lessThan(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !base.validate(bound, "ii") {
		return mentalese.Binding{}
	}

	int1, _ := strconv.Atoi(bound.Arguments[0].TermValue)
	int2, _ := strconv.Atoi(bound.Arguments[1].TermValue)

	if int1 < int2 {
		return binding
	} else {
		return mentalese.Binding{}
	}
}

func (base *SystemFunctionBase) add(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !base.validate(bound, "iiv") {
		return mentalese.Binding{}
	}

	int1, _ := strconv.Atoi(bound.Arguments[0].TermValue)
	int2, _ := strconv.Atoi(bound.Arguments[1].TermValue)

	result := int1 + int2

	newBinding := binding.Copy()
	newBinding[input.Arguments[2].TermValue] = mentalese.NewString(strconv.Itoa(result))

	return newBinding
}

func (base *SystemFunctionBase) subtract(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !base.validate(bound, "iiv") {
		return mentalese.Binding{}
	}

	int1, _ := strconv.Atoi(bound.Arguments[0].TermValue)
	int2, _ := strconv.Atoi(bound.Arguments[1].TermValue)

	result := int1 - int2

	newBinding := binding.Copy()
	newBinding[input.Arguments[2].TermValue] = mentalese.NewString(strconv.Itoa(result))

	return newBinding
}

func (base *SystemFunctionBase) equals(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {
	bound := input.BindSingle(binding)

	if !base.validate(bound, "--") {
		return mentalese.Binding{}
	}

	if !bound.Arguments[0].Equals(bound.Arguments[1]) {
		return mentalese.Binding{}
	} else {
		return binding
	}
}

func (base *SystemFunctionBase) notEquals(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !base.validate(bound, "--") {
		return mentalese.Binding{}
	}

	if bound.Arguments[0].Equals(bound.Arguments[1]) {
		return mentalese.Binding{}
	} else {
		return binding
	}
}

func (base *SystemFunctionBase) dateToday(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !base.validate(bound, "v") {
		return mentalese.Binding{}
	}

	now := time.Now()
	formatted := now.Format("2006-01-02")

	newBinding := binding.Copy()
	newBinding[input.Arguments[0].TermValue] = mentalese.NewString(formatted)

	return newBinding
}

func (base *SystemFunctionBase) dateSubtractYears(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !base.validate(bound, "ssv") {
		return mentalese.Binding{}
	}

	date1, err1 := time.Parse("2006-01-02", bound.Arguments[0].TermValue)
	date2, err2 := time.Parse("2006-01-02", bound.Arguments[1].TermValue)
	years := 0.0

	newBinding := binding.Copy()
	if err1 != nil || err2 != nil {
		newBinding = mentalese.Binding{}
	} else {
		duration := date1.Sub(date2)
		hours := duration.Hours()
		totalDays := hours / 24
		years = totalDays / 365
		newBinding[input.Arguments[2].TermValue] = mentalese.NewString(strconv.Itoa(int(years)))
	}

	return newBinding
}

func (base *SystemFunctionBase) Execute(input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	newBinding := binding
	found := true

	switch input.Predicate {
	case "split":
		newBinding = base.split(input, binding)
	case "join":
		newBinding = base.join(input, binding)
	case "concat":
		newBinding = base.concat(input, binding)
	case "greater_than":
		newBinding = base.greaterThan(input, binding)
	case "less_than":
		newBinding = base.lessThan(input, binding)
	case "add":
		newBinding = base.add(input, binding)
	case "subtract":
		newBinding = base.subtract(input, binding)
	case "equals":
		newBinding = base.equals(input, binding)
	case "not_equals":
		newBinding = base.notEquals(input, binding)
	case "date_today":
		newBinding = base.dateToday(input, binding)
	case "date_subtract_years":
		newBinding = base.dateSubtractYears(input, binding)
	default:
		found = false
	}

	return newBinding, found
}
