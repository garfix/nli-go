package knowledge

import (
	"fmt"
	"nli-go/lib/api"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
	"strings"
	"time"
)

type SystemFunctionBase struct {
	KnowledgeBaseCore
	matcher *central.RelationMatcher
	meta    *mentalese.Meta
	log     *common.SystemLog
}

func NewSystemFunctionBase(name string, meta *mentalese.Meta, log *common.SystemLog) *SystemFunctionBase {
	return &SystemFunctionBase{
		log:               log,
		KnowledgeBaseCore: KnowledgeBaseCore{name},
		meta:              meta,
		matcher:           central.NewRelationMatcher(log),
	}
}

func (base *SystemFunctionBase) GetFunctions() map[string]api.SimpleFunction {
	return map[string]api.SimpleFunction{
		mentalese.PredicateBound:             base.bound,
		mentalese.PredicateFree:              base.free,
		mentalese.PredicateAtom:              base.atom,
		mentalese.PredicateSplit:             base.split,
		mentalese.PredicateJoin:              base.join,
		mentalese.PredicateConcat:            base.concat,
		mentalese.PredicateGreaterThan:       base.greaterThan,
		mentalese.PredicateLessThan:          base.lessThan,
		mentalese.PredicateGreaterThanEquals: base.greaterThanEquals,
		mentalese.PredicateLessThanEquals:    base.lessThanEquals,
		mentalese.PredicateEquals:            base.equals,
		mentalese.PredicateNotEquals:         base.notEquals,
		mentalese.PredicateCompare:           base.compare,
		mentalese.PredicateUnify:             base.unify,
		mentalese.PredicateAdd:               base.add,
		mentalese.PredicateSubtract:          base.subtract,
		mentalese.PredicateMultiply:          base.multiply,
		mentalese.PredicateDivide:            base.divide,
		mentalese.PredicateMin:               base.min,
		mentalese.PredicateDateToday:         base.dateToday,
		mentalese.PredicateDateSubtractYears: base.dateSubtractYears,
		mentalese.PredicateLog:               base.debug,
		mentalese.PredicateUuid:              base.uuid,
		mentalese.PredicateHasSort:           base.hasSort,
		mentalese.PredicateListLength:        base.listLength,
		mentalese.PredicateListGet:           base.listGet,
		mentalese.PredicateListHead:          base.listHead,
	}
}

func (base *SystemFunctionBase) typeFunction(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {
	// dummy function that ensures that a call to go:type() does not result in "predicate not supported"
	return binding, false
}

func (base *SystemFunctionBase) bound(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {
	if input.Arguments[0].IsVariable() {
		return mentalese.NewBinding(), false
	} else {
		return binding, true
	}
}

func (base *SystemFunctionBase) free(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {
	if input.Arguments[0].IsVariable() {
		return binding, true
	} else {
		return mentalese.NewBinding(), false
	}
}

func (base *SystemFunctionBase) atom(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {
	inVar := input.Arguments[0]
	outVar := input.Arguments[1].TermValue

	binding.Set(outVar, inVar.ConvertVariablesToConstants())
	return binding, true
}

func (base *SystemFunctionBase) split(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "ssV", base.log) {
		return mentalese.NewBinding(), false
	}

	newBinding := binding.Copy()
	parts := strings.Split(bound.Arguments[0].TermValue, bound.Arguments[1].TermValue)

	for i, argument := range bound.Arguments[2:] {
		newBinding.Set(argument.TermValue, mentalese.NewTermString(parts[i]))
	}

	return newBinding, true
}

func (base *SystemFunctionBase) join(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "vsS", base.log) {
		return mentalese.NewBinding(), false
	}

	newBinding := binding.Copy()
	sep := ""
	result := ""
	for _, argument := range bound.Arguments[2:] {
		result += sep + argument.TermValue
		sep = input.Arguments[1].TermValue
	}

	newBinding.Set(input.Arguments[0].TermValue, mentalese.NewTermString(result))

	return newBinding, true
}

func (base *SystemFunctionBase) concat(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "vS", base.log) {
		return mentalese.NewBinding(), false
	}

	newBinding := binding.Copy()
	result := ""
	for _, argument := range bound.Arguments[1:] {
		result += argument.TermValue
	}

	newBinding.Set(input.Arguments[0].TermValue, mentalese.NewTermString(result))

	return newBinding, true
}

func (base *SystemFunctionBase) greaterThan(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "nn", base.log) {
		return mentalese.NewBinding(), false
	}

	num1, _ := bound.Arguments[0].GetNumber()
	num2, _ := bound.Arguments[1].GetNumber()

	if num1 > num2 {
		return binding, true
	} else {
		return mentalese.NewBinding(), false
	}
}

func (base *SystemFunctionBase) lessThan(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "nn", base.log) {
		return mentalese.NewBinding(), false
	}

	num1, _ := bound.Arguments[0].GetNumber()
	num2, _ := bound.Arguments[1].GetNumber()

	if num1 < num2 {
		return binding, true
	} else {
		return mentalese.NewBinding(), false
	}
}

func (base *SystemFunctionBase) greaterThanEquals(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "nn", base.log) {
		return mentalese.NewBinding(), false
	}

	num1, _ := bound.Arguments[0].GetNumber()
	num2, _ := bound.Arguments[1].GetNumber()

	if num1 >= num2 {
		return binding, true
	} else {
		return mentalese.NewBinding(), false
	}
}

func (base *SystemFunctionBase) lessThanEquals(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "nn", base.log) {
		return mentalese.NewBinding(), false
	}

	num1, _ := bound.Arguments[0].GetNumber()
	num2, _ := bound.Arguments[1].GetNumber()

	if num1 <= num2 {
		return binding, true
	} else {
		return mentalese.NewBinding(), false
	}
}

func (base *SystemFunctionBase) add(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "***", base.log) {
		return mentalese.NewBinding(), false
	}

	variable := input.Arguments[2].TermValue
	newBinding := binding.Copy()
	var value mentalese.Term

	if bound.Arguments[0].IsNumber() {

		num1, _ := bound.Arguments[0].GetNumber()
		num2, _ := bound.Arguments[1].GetNumber()

		result := num1 + num2
		resultString := strconv.FormatFloat(result, 'f', -1, 64)
		value = mentalese.NewTermString(resultString)

	} else if bound.Arguments[0].IsRelationSet() || bound.Arguments[1].IsRelationSet() {

		set1 := bound.Arguments[0].TermValueRelationSet
		set2 := bound.Arguments[1].TermValueRelationSet

		result := set1.Copy()
		result = append(result, set2...)
		value = mentalese.NewTermRelationSet(result)

	} else {
		return mentalese.NewBinding(), false
	}

	newBinding.Set(variable, value)

	return newBinding, true
}

func (base *SystemFunctionBase) subtract(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "nn*", base.log) {
		return mentalese.NewBinding(), false
	}

	num1, _ := bound.Arguments[0].GetNumber()
	num2, _ := bound.Arguments[1].GetNumber()

	result := num1 - num2

	resultString := strconv.FormatFloat(result, 'f', -1, 64)

	newBinding := binding.Copy()
	newBinding.Set(input.Arguments[2].TermValue, mentalese.NewTermString(resultString))

	return newBinding, true
}

func (base *SystemFunctionBase) multiply(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "nn*", base.log) {
		return mentalese.NewBinding(), false
	}

	num1, _ := bound.Arguments[0].GetNumber()
	num2, _ := bound.Arguments[1].GetNumber()

	result := num1 * num2

	resultString := strconv.FormatFloat(result, 'f', -1, 64)

	newBinding := binding.Copy()
	newBinding.Set(input.Arguments[2].TermValue, mentalese.NewTermString(resultString))

	return newBinding, true
}

func (base *SystemFunctionBase) divide(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "nn*", base.log) {
		return mentalese.NewBinding(), false
	}

	num1, _ := bound.Arguments[0].GetNumber()
	num2, _ := bound.Arguments[1].GetNumber()

	result := num1 / num2

	resultString := strconv.FormatFloat(result, 'f', -1, 64)

	newBinding := binding.Copy()
	newBinding.Set(input.Arguments[2].TermValue, mentalese.NewTermString(resultString))

	return newBinding, true
}

func (base *SystemFunctionBase) min(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "iiv", base.log) {
		return mentalese.NewBinding(), false
	}

	int1, _ := strconv.Atoi(bound.Arguments[0].TermValue)
	int2, _ := strconv.Atoi(bound.Arguments[1].TermValue)

	result := int1
	if int2 < int1 {
		result = int2
	}

	newBinding := binding.Copy()
	newBinding.Set(input.Arguments[2].TermValue, mentalese.NewTermString(strconv.Itoa(result)))

	return newBinding, true
}

func (base *SystemFunctionBase) compare(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "ssv", base.log) {
		return mentalese.NewBinding(), false
	}

	n1 := bound.Arguments[0].TermValue
	n2 := bound.Arguments[1].TermValue

	result := strings.Compare(n1, n2)

	newBinding := binding.Copy()
	newBinding.Set(input.Arguments[2].TermValue, mentalese.NewTermString(strconv.Itoa(result)))

	return newBinding, true
}

func (base *SystemFunctionBase) equals(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {
	bound := input.BindSingle(binding)

	if !Validate(bound, "--", base.log) {
		return mentalese.NewBinding(), false
	}

	if !bound.Arguments[0].Equals(bound.Arguments[1]) {
		return mentalese.NewBinding(), false
	} else {
		return binding, true
	}
}

func (base *SystemFunctionBase) unify(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	if !Validate(input, "--", base.log) {
		return mentalese.NewBinding(), false
	}

	bound := input.BindSingle(binding)

	firstBinding, match1 := base.matcher.MatchTerm(bound.Arguments[0], bound.Arguments[1], mentalese.NewBinding())
	secondBinding, match2 := base.matcher.MatchTerm(bound.Arguments[1], bound.Arguments[0], mentalese.NewBinding())
	combinedBinding := firstBinding.Merge(secondBinding).RemoveVariables()
	newBinding := binding.Merge(combinedBinding)

	if !match1 || !match2 {
		return mentalese.NewBinding(), false
	} else {
		return newBinding, true
	}
}

func (base *SystemFunctionBase) notEquals(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "--", base.log) {
		return mentalese.NewBinding(), false
	}

	if bound.Arguments[0].Equals(bound.Arguments[1]) {
		return mentalese.NewBinding(), false
	} else {
		return binding, true
	}
}

func (base *SystemFunctionBase) dateToday(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "v", base.log) {
		return mentalese.NewBinding(), false
	}

	now := time.Now()
	formatted := now.Format("2006-01-02")

	newBinding := binding.Copy()
	newBinding.Set(input.Arguments[0].TermValue, mentalese.NewTermString(formatted))

	return newBinding, true
}

func (base *SystemFunctionBase) dateSubtractYears(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "ssv", base.log) {
		return mentalese.NewBinding(), false
	}

	date1, err1 := time.Parse("2006-01-02", bound.Arguments[0].TermValue)
	date2, err2 := time.Parse("2006-01-02", bound.Arguments[1].TermValue)

	newBinding := binding.Copy()
	if err1 != nil || err2 != nil {
		return mentalese.NewBinding(), false
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

		newBinding.Set(input.Arguments[2].TermValue, mentalese.NewTermString(strconv.Itoa(int(years))))
	}

	return newBinding, true
}

func (base *SystemFunctionBase) debug(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	log := ""
	sep := ""

	for i, argument := range input.Arguments {
		if argument.IsVariable() {
			value, found := binding.Get(argument.TermValue)
			if found {
				log += sep + argument.TermValue + ": " + value.String()
			} else {
				log += sep + argument.TermValue + ": <not bound>"
			}
		} else {
			log += sep + bound.Arguments[i].String()
		}
		sep = "\t"
	}

	fmt.Println(log)

	return binding, true
}

func (base *SystemFunctionBase) uuid(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "va", base.log) {
		return mentalese.NewBinding(), false
	}

	u := bound.Arguments[0].TermValue
	sort := bound.Arguments[1].TermValue

	newBinding := binding.Copy()
	newBinding.Set(u, mentalese.NewTermId(common.CreateUuid(), sort))

	return newBinding, true
}

func (base *SystemFunctionBase) hasSort(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {
	// this function is used internally, so it needs to exist
	return binding, false
}

func (base *SystemFunctionBase) listLength(messenger api.SimpleMessenger, relation mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := relation.BindSingle(binding)

	if !Validate(bound, "lv", base.log) {
		return mentalese.NewBinding(), false
	}

	list := bound.Arguments[0].TermValueList
	lengthVar := bound.Arguments[1].TermValue

	length := len(list)

	newBinding := binding.Copy()
	newBinding.Set(lengthVar, mentalese.NewTermString(strconv.Itoa(length)))
	return newBinding, true
}

func (base *SystemFunctionBase) listGet(messenger api.SimpleMessenger, relation mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := relation.BindSingle(binding)

	if !Validate(bound, "li*", base.log) {
		return mentalese.NewBinding(), false
	}

	list := bound.Arguments[0].TermValueList
	index := bound.Arguments[1].TermValue
	termVar := relation.Arguments[2].TermValue

	i, err := strconv.Atoi(index)
	if err != nil {
		base.log.AddError("Index should be an integer: " + index)
		return mentalese.NewBinding(), false
	}

	if i < 0 || i >= len(list) {
		return mentalese.NewBinding(), false
	}

	term := list[i]

	newBinding := binding.Copy()
	newBinding.Set(termVar, term)

	return newBinding, true
}

func (base *SystemFunctionBase) listHead(messenger api.SimpleMessenger, relation mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := relation.BindSingle(binding)

	if !Validate(bound, "lvv", base.log) {
		return mentalese.NewBinding(), false
	}

	list := bound.Arguments[0].TermValueList
	headVar := bound.Arguments[1].TermValue
	tailVar := relation.Arguments[2].TermValue

	newBinding := mentalese.NewBinding()
	newBinding.Set(headVar, list[0])
	newBinding.Set(tailVar, mentalese.NewTermList(list[1:]))
	return newBinding, true
}
