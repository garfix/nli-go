package knowledge

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
	"strings"
	"time"
)

type SystemFunctionBase struct {
	KnowledgeBaseCore
	matcher *mentalese.RelationMatcher
	log *common.SystemLog
}

func NewSystemFunctionBase(name string, log *common.SystemLog) *SystemFunctionBase {
	return &SystemFunctionBase{ log: log, KnowledgeBaseCore: KnowledgeBaseCore{ name }, matcher: mentalese.NewRelationMatcher(log) }
}

func (base *SystemFunctionBase) HandlesPredicate(predicate string) bool {
	predicates := []string{"split", "join", "concat", "greater_than", "less_than", "equals", "not_equals", "unify", "add", "subtract", "date_today", "date_subtract_years"}

	for _, p := range predicates {
		if p == predicate {
			return true
		}
	}
	return false
}

func (base *SystemFunctionBase) split(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !validate(bound, "ssV", base.log) {
		return nil
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

	if !validate(bound, "vsS", base.log) {
		return nil
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

	if !validate(bound, "vS", base.log) {
		return nil
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

	if !validate(bound, "ii", base.log) {
		return nil
	}

	int1, _ := strconv.Atoi(bound.Arguments[0].TermValue)
	int2, _ := strconv.Atoi(bound.Arguments[1].TermValue)

	if int1 > int2 {
		return binding
	} else {
		return nil
	}
}

func (base *SystemFunctionBase) lessThan(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !validate(bound, "ii", base.log) {
		return nil
	}

	int1, _ := strconv.Atoi(bound.Arguments[0].TermValue)
	int2, _ := strconv.Atoi(bound.Arguments[1].TermValue)

	if int1 < int2 {
		return binding
	} else {
		return nil
	}
}

func (base *SystemFunctionBase) add(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !validate(bound, "iiv", base.log) {
		return nil
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

	if !validate(bound, "iiv", base.log) {
		return nil
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

	if !validate(bound, "--", base.log) {
		return nil
	}

	if !bound.Arguments[0].Equals(bound.Arguments[1]) {
		return nil
	} else {
		return binding
	}
}

func (base *SystemFunctionBase) unify(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	if !validate(input, "--", base.log) {
		return nil
	}

	bound := input.BindSingle(binding)

	firstBinding, match1 := base.matcher.MatchTerm(bound.Arguments[0], bound.Arguments[1], mentalese.Binding{})
	secondBinding, match2 := base.matcher.MatchTerm(bound.Arguments[1], bound.Arguments[0], mentalese.Binding{})
	combinedBinding := firstBinding.Merge(secondBinding).RemoveVariables()
	newBinding := binding.Merge(combinedBinding)

	if !match1 || !match2 {
		return nil
	} else {
		return newBinding
	}
}

func (base *SystemFunctionBase) notEquals(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !validate(bound, "--", base.log) {
		return nil
	}

	if bound.Arguments[0].Equals(bound.Arguments[1]) {
		return nil
	} else {
		return binding
	}
}

func (base *SystemFunctionBase) dateToday(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !validate(bound, "v", base.log) {
		return nil
	}

	now := time.Now()
	formatted := now.Format("2006-01-02")

	newBinding := binding.Copy()
	newBinding[input.Arguments[0].TermValue] = mentalese.NewString(formatted)

	return newBinding
}

func (base *SystemFunctionBase) dateSubtractYears(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !validate(bound, "ssv", base.log) {
		return nil
	}

	date1, err1 := time.Parse("2006-01-02", bound.Arguments[0].TermValue)
	date2, err2 := time.Parse("2006-01-02", bound.Arguments[1].TermValue)
	years := 0.0

	newBinding := binding.Copy()
	if err1 != nil || err2 != nil {
		newBinding = nil
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
	case "unify":
		newBinding = base.unify(input, binding)
	case "date_today":
		newBinding = base.dateToday(input, binding)
	case "date_subtract_years":
		newBinding = base.dateSubtractYears(input, binding)
	default:
		found = false
	}

	return newBinding, found
}
