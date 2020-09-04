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
	predicates := []string{
		mentalese.PredicateSplit,
		mentalese.PredicateJoin,
		mentalese.PredicateConcat,
		mentalese.PredicateGreaterThan,
		mentalese.PredicateLessThan,
		mentalese.PredicateGreaterThanEquals,
		mentalese.PredicateLessThanEquals,
		mentalese.PredicateEquals,
		mentalese.PredicateNotEquals,
		mentalese.PredicateCompare,
		mentalese.PredicateUnify,
		mentalese.PredicateAdd,
		mentalese.PredicateSubtract,
		mentalese.PredicateMin,
		mentalese.PredicateDateToday,
		mentalese.PredicateDateSubtractYears,
	}

	for _, p := range predicates {
		if p == predicate {
			return true
		}
	}
	return false
}

func (base *SystemFunctionBase) split(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !Validate(bound, "ssV", base.log) {
		return nil
	}

	newBinding := binding.Copy()
	parts := strings.Split(bound.Arguments[0].TermValue, bound.Arguments[1].TermValue)

	for i, argument := range bound.Arguments[2:] {
		newBinding[argument.TermValue] = mentalese.NewTermString(parts[i])
	}

	return newBinding
}

func (base *SystemFunctionBase) join(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !Validate(bound, "vsS", base.log) {
		return nil
	}

	newBinding := binding.Copy()
	sep := ""
	result := ""
	for _, argument := range bound.Arguments[2:] {
		result += sep + argument.TermValue
		sep = input.Arguments[1].TermValue
	}

	newBinding[input.Arguments[0].TermValue] = mentalese.NewTermString(result)

	return newBinding
}

func (base *SystemFunctionBase) concat(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !Validate(bound, "vS", base.log) {
		return nil
	}

	newBinding := binding.Copy()
	result := ""
	for _, argument := range bound.Arguments[1:] {
		result += argument.TermValue
	}

	newBinding[input.Arguments[0].TermValue] = mentalese.NewTermString(result)

	return newBinding
}

func (base *SystemFunctionBase) greaterThan(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !Validate(bound, "ii", base.log) {
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

	if !Validate(bound, "ii", base.log) {
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

func (base *SystemFunctionBase) greaterThanEquals(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !Validate(bound, "ii", base.log) {
		return nil
	}

	int1, _ := strconv.Atoi(bound.Arguments[0].TermValue)
	int2, _ := strconv.Atoi(bound.Arguments[1].TermValue)

	if int1 >= int2 {
		return binding
	} else {
		return nil
	}
}

func (base *SystemFunctionBase) lessThanEquals(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !Validate(bound, "ii", base.log) {
		return nil
	}

	int1, _ := strconv.Atoi(bound.Arguments[0].TermValue)
	int2, _ := strconv.Atoi(bound.Arguments[1].TermValue)

	if int1 <= int2 {
		return binding
	} else {
		return nil
	}
}

func (base *SystemFunctionBase) add(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !Validate(bound, "iiv", base.log) {
		return nil
	}

	int1, _ := strconv.Atoi(bound.Arguments[0].TermValue)
	int2, _ := strconv.Atoi(bound.Arguments[1].TermValue)

	result := int1 + int2

	newBinding := binding.Copy()
	newBinding[input.Arguments[2].TermValue] = mentalese.NewTermString(strconv.Itoa(result))

	return newBinding
}

func (base *SystemFunctionBase) subtract(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !Validate(bound, "iiv", base.log) {
		return nil
	}

	int1, _ := strconv.Atoi(bound.Arguments[0].TermValue)
	int2, _ := strconv.Atoi(bound.Arguments[1].TermValue)

	result := int1 - int2

	newBinding := binding.Copy()
	newBinding[input.Arguments[2].TermValue] = mentalese.NewTermString(strconv.Itoa(result))

	return newBinding
}

func (base *SystemFunctionBase) min(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !Validate(bound, "iiv", base.log) {
		return nil
	}

	int1, _ := strconv.Atoi(bound.Arguments[0].TermValue)
	int2, _ := strconv.Atoi(bound.Arguments[1].TermValue)

	result := int1
	if int2 < int1 {
		result = int2
	}

	newBinding := binding.Copy()
	newBinding[input.Arguments[2].TermValue] = mentalese.NewTermString(strconv.Itoa(result))

	return newBinding
}

func (base *SystemFunctionBase) compare(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !Validate(bound, "ssv", base.log) {
		return nil
	}

	n1 := bound.Arguments[0].TermValue
	n2 := bound.Arguments[1].TermValue

	result := strings.Compare(n1, n2)

	newBinding := binding.Copy()
	newBinding[input.Arguments[2].TermValue] = mentalese.NewTermString(strconv.Itoa(result))

	return newBinding
}

func (base *SystemFunctionBase) equals(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {
	bound := input.BindSingle(binding)

	if !Validate(bound, "--", base.log) {
		return nil
	}

	if !bound.Arguments[0].Equals(bound.Arguments[1]) {
		return nil
	} else {
		return binding
	}
}

func (base *SystemFunctionBase) unify(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	if !Validate(input, "--", base.log) {
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

	if !Validate(bound, "--", base.log) {
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

	if !Validate(bound, "v", base.log) {
		return nil
	}

	now := time.Now()
	formatted := now.Format("2006-01-02")

	newBinding := binding.Copy()
	newBinding[input.Arguments[0].TermValue] = mentalese.NewTermString(formatted)

	return newBinding
}

func (base *SystemFunctionBase) dateSubtractYears(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	if !Validate(bound, "ssv", base.log) {
		return nil
	}

	date1, err1 := time.Parse("2006-01-02", bound.Arguments[0].TermValue)
	date2, err2 := time.Parse("2006-01-02", bound.Arguments[1].TermValue)

	newBinding := binding.Copy()
	if err1 != nil || err2 != nil {
		newBinding = nil
	} else {
		//years := 0
		//if date1.Year() < date2.Year() {
		//	years = date1.Year() - date2.Year()
		//} else if date1.YearDay() < date2.YearDay() {
		//	years = date1.Year() - date2.Year() - 1
		//} else {
		//	years = date1.Year() - date2.Year()
		//}

		years := 0
		if date1.YearDay() < date2.YearDay() {
			years = date1.Year() - date2.Year() - 1
		} else {
			years = date1.Year() - date2.Year()
		}

		newBinding[input.Arguments[2].TermValue] = mentalese.NewTermString(strconv.Itoa(int(years)))
	}

	return newBinding
}

func (base *SystemFunctionBase) Execute(input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	newBinding := binding
	found := true

	switch input.Predicate {
	case mentalese.PredicateSplit:
		newBinding = base.split(input, binding)
	case mentalese.PredicateJoin:
		newBinding = base.join(input, binding)
	case mentalese.PredicateConcat:
		newBinding = base.concat(input, binding)
	case mentalese.PredicateGreaterThan:
		newBinding = base.greaterThan(input, binding)
	case mentalese.PredicateLessThan:
		newBinding = base.lessThan(input, binding)
	case mentalese.PredicateGreaterThanEquals:
		newBinding = base.greaterThanEquals(input, binding)
	case mentalese.PredicateLessThanEquals:
		newBinding = base.lessThanEquals(input, binding)
	case mentalese.PredicateAdd:
		newBinding = base.add(input, binding)
	case mentalese.PredicateSubtract:
		newBinding = base.subtract(input, binding)
	case mentalese.PredicateMin:
		newBinding = base.min(input, binding)
	case mentalese.PredicateEquals:
		newBinding = base.equals(input, binding)
	case mentalese.PredicateCompare:
		newBinding = base.compare(input, binding)
	case mentalese.PredicateNotEquals:
		newBinding = base.notEquals(input, binding)
	case mentalese.PredicateUnify:
		newBinding = base.unify(input, binding)
	case mentalese.PredicateDateToday:
		newBinding = base.dateToday(input, binding)
	case mentalese.PredicateDateSubtractYears:
		newBinding = base.dateSubtractYears(input, binding)
	default:
		found = false
	}

	return newBinding, found
}
