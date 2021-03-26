package knowledge

import (
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
)

type SystemMultiBindingFunctionBase struct {
	KnowledgeBaseCore
	rules []mentalese.Rule
	log   *common.SystemLog
}

func NewSystemMultiBindingBase(name string, log *common.SystemLog) *SystemMultiBindingFunctionBase {
	return &SystemMultiBindingFunctionBase{KnowledgeBaseCore: KnowledgeBaseCore{ Name: name }, log: log}
}

func (base *SystemMultiBindingFunctionBase) GetFunctions() map[string]api.MultiBindingFunction {
	return map[string]api.MultiBindingFunction{
		mentalese.PredicateNumberOf: base.numberOf,
		mentalese.PredicateFirst:    base.first,
		mentalese.PredicateLast:     base.last,
		mentalese.PredicateGet:      base.get,
		mentalese.PredicateOrder:    base.order,
		mentalese.PredicateLargest:  base.largest,
		mentalese.PredicateSmallest: base.smallest,
		mentalese.PredicateExists:   base.exists,
		mentalese.PredicateCut: 	 base.cut,
		mentalese.PredicateMakeAnd:  base.makeAnd,
		mentalese.PredicateMakeList: base.makeList,
	}
}

func (base *SystemMultiBindingFunctionBase) numberOf(messenger api.ProcessMessenger, input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

	if !Validate(input, "--", base.log) {
		return mentalese.NewBindingSet()
	}

	subjectVariable := input.Arguments[0].TermValue
	numberArgumentValue := input.Arguments[1].TermValue
	number :=  bindings.GetDistinctValueCount(subjectVariable)

	newBindings := mentalese.NewBindingSet()

	if input.Arguments[1].IsVariable() {
		for _, binding := range bindings.GetAll() {
			newBinding := binding.Copy()
			newBinding.Set(numberArgumentValue, mentalese.NewTermString(strconv.Itoa(number)))
			newBindings.Add(newBinding)
		}
	} else {
		assertedNumber, err := strconv.Atoi(numberArgumentValue)
		if err != nil {
			base.log.AddError("The second argument of number_of() needs to be an integer")
			newBindings = mentalese.NewBindingSet()
		} else {
			if number == assertedNumber {
				newBindings = bindings
			} else {
				newBindings = mentalese.NewBindingSet()
			}
		}
	}

	return newBindings
}

func (base *SystemMultiBindingFunctionBase) first(messenger api.ProcessMessenger, input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

	length := 0

	if bindings.IsEmpty() {
		return bindings
	}

	bound := input.BindSingle(bindings.Get(0))

	if len(input.Arguments) == 0 {
		length = 1
	} else if len(input.Arguments) == 1 {

		intValue, ok := bound.Arguments[0].GetIntValue()
		if !ok {
			base.log.AddError("First argument of `first` must be an integer")
			return mentalese.NewBindingSet()
		}
		length = intValue
	} else {
		base.log.AddError("`first` takes at most one argument")
		return mentalese.NewBindingSet()
	}

	newBindings := mentalese.NewBindingSet()
	if bindings.IsEmpty() {
		newBindings = bindings
	} else {
		i := 0
		for _, binding := range bindings.GetAll() {
			newBinding := binding.Copy()
			newBindings.Add(newBinding)
			i++
			if i == length {
				break
			}
			if i == bindings.GetLength() {
				break
			}
		}
	}

	return newBindings
}

func (base *SystemMultiBindingFunctionBase) last(messenger api.ProcessMessenger, input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

	length := 0

	if bindings.IsEmpty() {
		return bindings
	}

	bound := input.BindSingle(bindings.Get(0))

	if len(input.Arguments) == 0 {
		length = 1
	} else if len(input.Arguments) == 1 {

		intValue, ok := bound.Arguments[0].GetIntValue()
		if !ok {
			base.log.AddError("First argument of `last` must be an integer")
			return mentalese.NewBindingSet()
		}
		length = intValue
	} else {
		base.log.AddError("`last` takes at most one argument")
		return mentalese.NewBindingSet()
	}

	newBindings := mentalese.NewBindingSet()
	if bindings.IsEmpty() {
		newBindings = bindings
	} else {
		all := bindings.GetAll()
		for i := len(all) - length; i < len(all); i++ {
			if i < 0 {
				continue
			}
			newBinding := all[i].Copy()
			newBindings.Add(newBinding)
		}
	}

	return newBindings
}

func (base *SystemMultiBindingFunctionBase) get(messenger api.ProcessMessenger, input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

	start := 0
	length := 1

	if bindings.IsEmpty() {
		return bindings
	}

	bound := input.BindSingle(bindings.Get(0))

	if len(input.Arguments) < 1 || len(input.Arguments) > 2 {
		base.log.AddError("`get` takes one or two arguments")
		return mentalese.NewBindingSet()
	}

	{
		intValue, ok := bound.Arguments[0].GetIntValue()
		if !ok {
			base.log.AddError("First argument of `get` must be an integer")
			return mentalese.NewBindingSet()
		}

		start = intValue
	}

	if len(input.Arguments) == 2 {

		intValue, ok := bound.Arguments[1].GetIntValue()
		if !ok {
			base.log.AddError("Second argument of `get` must be an integer")
			return mentalese.NewBindingSet()
		}

		length = intValue
	}

	newBindings := mentalese.NewBindingSet()
	all := bindings.GetAll()
	for i := start; i < start + length; i++ {
		if i < 0 {
			continue
		}
		if i >= len(all) {
			break
		}
		newBinding := all[i].Copy()
		newBindings.Add(newBinding)
	}

	return newBindings
}

func (base *SystemMultiBindingFunctionBase) order(messenger api.ProcessMessenger, input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

	if !Validate(input, "v*", base.log) {
		return mentalese.NewBindingSet()
	}

	if bindings.IsEmpty() {
		return mentalese.NewBindingSet()
	}

	bound := input.BindSingle(bindings.Get(0))

	subjectVariable := input.Arguments[0].TermValue
	direction := bound.Arguments[1].TermValue

	if direction != mentalese.AtomAsc && direction != mentalese.AtomDesc {
		base.log.AddError("Direction of `order` must be either asc or desc")
		return mentalese.NewBindingSet()
	}

	newBindings, ok := bindings.Sort(subjectVariable)
	if !ok {
		base.log.AddError("`order` variable should contain only integers or strings")
		return mentalese.NewBindingSet()
	}

	if direction == mentalese.AtomDesc {
		newBindings = newBindings.Reverse()
	}

	return newBindings
}

func (base *SystemMultiBindingFunctionBase) largest(messenger api.ProcessMessenger, input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

	if !Validate(input, "v", base.log) {
		return mentalese.NewBindingSet()
	}

	if bindings.IsEmpty() {
		return mentalese.NewBindingSet()
	}

	subjectVariable := input.Arguments[0].TermValue
	distinctValues := bindings.GetDistinctValues(subjectVariable)

	largest := 0.0

	for i, value := range distinctValues {
		value, err := strconv.ParseFloat(value.TermValue, 64)
		if err != nil {
			base.log.AddError("Largest takes all numbers")
			return mentalese.NewBindingSet()
		}
		if i == 0 {
			largest = value
		} else if value > largest {
			largest = value
		}
	}

	newBindings := mentalese.NewBindingSet()
	for _, binding := range bindings.GetAll() {
		value, found := binding.Get(subjectVariable)
		if found {
			value, _ := strconv.ParseFloat(value.TermValue, 64)
			if value < largest {
				continue
			}
		}
		newBindings.Add(binding)
	}

	return newBindings
}

func (base *SystemMultiBindingFunctionBase) smallest(messenger api.ProcessMessenger, input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

	if !Validate(input, "v", base.log) {
		return mentalese.NewBindingSet()
	}

	if bindings.IsEmpty() {
		return mentalese.NewBindingSet()
	}

	subjectVariable := input.Arguments[0].TermValue
	distinctValues := bindings.GetDistinctValues(subjectVariable)

	smallest := 0.0

	for i, value := range distinctValues {
		value, err := strconv.ParseFloat(value.TermValue, 64)
		if err != nil {
			base.log.AddError("Smallest takes all numbers")
			return mentalese.NewBindingSet()
		}
		if i == 0 {
			smallest = value
		} else if value < smallest {
			smallest = value
		}
	}

	newBindings := mentalese.NewBindingSet()
	for _, binding := range bindings.GetAll() {
		value, found := binding.Get(subjectVariable)
		if found {
			value, _ := strconv.ParseFloat(value.TermValue, 64)
			if value > smallest {
				continue
			}
		}
		newBindings.Add(binding)
	}

	return newBindings
}

func (base *SystemMultiBindingFunctionBase) exists(messenger api.ProcessMessenger, input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

	if !Validate(input, "", base.log) {
		return mentalese.NewBindingSet()
	}

	return bindings
}

func (base *SystemMultiBindingFunctionBase) cut(messenger api.ProcessMessenger, input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

	if !Validate(input, "ir", base.log) {
		return mentalese.NewBindingSet()
	}

	count, _ := input.Arguments[0].GetIntValue()
	childRelations := input.Arguments[1].TermValueRelationSet

	cursor := messenger.GetCursor()

	index := cursor.GetState("index", 0)
	passed := cursor.GetState("passed", 0)

	if index > 0 {
		childBindings := cursor.GetChildFrameResultBindings()

		if !childBindings.IsEmpty() {
			cursor.AddStepBindings(childBindings)
			passed += 1
			cursor.SetState("passed", passed)

			if passed == count {
				resultBindings := mentalese.NewBindingSet()
				for _, bindingSet := range cursor.GetAllStepBindings() {
					resultBindings.AddMultiple(bindingSet)
				}
				return resultBindings
			}
		}
	}

	if index < bindings.GetLength() {
		binding := bindings.Get(index)
		messenger.CreateChildStackFrame(childRelations, mentalese.InitBindingSet(binding))

		cursor.SetState("index", index + 1)
	}

	return mentalese.NewBindingSet()
}

func (base *SystemMultiBindingFunctionBase) makeAnd(messenger api.ProcessMessenger, input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

	if !Validate(input, "vv", base.log) {
		return mentalese.NewBindingSet()
	}

	entityVar := input.Arguments[0].TermValue
	andVar := input.Arguments[1].TermValue

	uniqueValues := bindings.GetDistinctValues(entityVar)
	relation := mentalese.Relation{}
	count := len(uniqueValues)

	if count == 0 {
		return bindings
	} else if count == 1 {
		relation = mentalese.NewRelation(false, mentalese.PredicateAnd, []mentalese.Term{
			uniqueValues[0],
			uniqueValues[0],
		})
	} else {
		relation = mentalese.NewRelation(false, mentalese.PredicateAnd, []mentalese.Term{
			uniqueValues[count - 2],
			uniqueValues[count - 1],
		})
		for i := len(uniqueValues)-3; i >= 0 ; i-- {
			relation = mentalese.NewRelation(false, mentalese.PredicateAnd, []mentalese.Term{
				uniqueValues[i],
				mentalese.NewTermRelationSet(mentalese.RelationSet{ relation }),
			})
		}
	}

	newBindings := mentalese.NewBindingSet()

	for _, binding := range bindings.GetAll() {
		newBinding := binding.Copy()
		newBinding.Set(andVar, mentalese.NewTermRelationSet(mentalese.RelationSet{ relation }))
		newBindings.Add(newBinding)
	}

	return newBindings
}

func (base *SystemMultiBindingFunctionBase) makeList(messenger api.ProcessMessenger, input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {
	if !Validate(input, "V", base.log) {
		return mentalese.NewBindingSet()
	}

	listVar := input.Arguments[0].TermValue
	list := mentalese.TermList{}

	variables := []string{}
	for i, argument := range input.Arguments {
		if i == 0 { continue }
		variable := argument.TermValue
		variables = append(variables, variable)
		for _, value := range bindings.GetDistinctValues(variable) {
			list = append(list, value)
		}
	}

	listTerm := mentalese.NewTermList(list)

	newBindings := mentalese.NewBindingSet()
	for _, binding := range bindings.GetAll() {
		newBinding := binding.Copy()
		newBinding.Set(listVar, listTerm)
		newBinding = newBinding.FilterOutVariablesByName(variables)
		newBindings.Add(newBinding)
	}
	return newBindings
}
