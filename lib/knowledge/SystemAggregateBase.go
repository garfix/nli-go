package knowledge

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
)

type SystemAggregateFunctionBase struct {
	KnowledgeBaseCore
	rules []mentalese.Rule
	log   *common.SystemLog
}

func NewSystemAggregateBase(name string, log *common.SystemLog) *SystemAggregateFunctionBase {
	return &SystemAggregateFunctionBase{KnowledgeBaseCore: KnowledgeBaseCore{ Name: name }, log: log}
}

func (base *SystemAggregateFunctionBase) HandlesPredicate(predicate string) bool {
	predicates := []string{
		mentalese.PredicateNumberOf,
		mentalese.PredicateFirst,
		mentalese.PredicateExists,
		mentalese.PredicateMakeAnd,
		mentalese.PredicateMakeList,
	}

	for _, p := range predicates {
		if p == predicate {
			return true
		}
	}
	return false
}

func (base *SystemAggregateFunctionBase) Execute(input mentalese.Relation, bindings mentalese.BindingSet) (mentalese.BindingSet, bool) {

	newBindings := bindings
	found := true

	switch input.Predicate {
	case mentalese.PredicateNumberOf:
		newBindings = base.numberOf(input, bindings)
	case mentalese.PredicateFirst:
		newBindings = base.first(input, bindings)
	case mentalese.PredicateExists:
		newBindings = base.exists(input, bindings)
	case mentalese.PredicateMakeAnd:
		newBindings = base.makeAnd(input, bindings)
	case mentalese.PredicateMakeList:
		newBindings = base.makeList(input, bindings)
	default:
		found = false
	}

	return newBindings, found
}

func (base *SystemAggregateFunctionBase) numberOf(input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

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

func (base *SystemAggregateFunctionBase) first(input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

	if !Validate(input, "v", base.log) {
		return mentalese.NewBindingSet()
	}

	subjectVariable := input.Arguments[0].TermValue
	distinctValues := bindings.GetDistinctValues(subjectVariable)

	newBindings := mentalese.NewBindingSet()
	if len(distinctValues) == 0 {
		newBindings = bindings
	} else {
		for _, binding := range bindings.GetAll() {
			newBinding := binding.Copy()
			newBinding.Set(subjectVariable, distinctValues[0])
			newBindings.Add(newBinding)
		}
	}

	return newBindings
}

func (base *SystemAggregateFunctionBase) exists(input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

	if !Validate(input, "", base.log) {
		return mentalese.NewBindingSet()
	}

	return bindings
}

func (base *SystemAggregateFunctionBase) makeAnd(input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

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
		relation = mentalese.NewRelation(true, mentalese.PredicateAnd, []mentalese.Term{
			uniqueValues[0],
			uniqueValues[0],
		})
	} else {
		relation = mentalese.NewRelation(true, mentalese.PredicateAnd, []mentalese.Term{
			uniqueValues[count - 2],
			uniqueValues[count - 1],
		})
		for i := len(uniqueValues)-3; i >= 0 ; i-- {
			relation = mentalese.NewRelation(true, mentalese.PredicateAnd, []mentalese.Term{
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

func (base *SystemAggregateFunctionBase) makeList(input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {
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
