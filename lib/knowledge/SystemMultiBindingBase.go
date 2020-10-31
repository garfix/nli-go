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
		mentalese.PredicateFirst: base.first,
		mentalese.PredicateExists: base.exists,
		mentalese.PredicateMakeAnd: base.makeAnd,
		mentalese.PredicateMakeList: base.makeList,
	}
}

func (base *SystemMultiBindingFunctionBase) numberOf(input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

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

func (base *SystemMultiBindingFunctionBase) first(input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

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

func (base *SystemMultiBindingFunctionBase) exists(input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

	if !Validate(input, "", base.log) {
		return mentalese.NewBindingSet()
	}

	return bindings
}

func (base *SystemMultiBindingFunctionBase) makeAnd(input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

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

func (base *SystemMultiBindingFunctionBase) makeList(input mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {
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
