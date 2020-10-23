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
	predicates := []string{
		mentalese.PredicateNumberOf,
		mentalese.PredicateFirst,
		mentalese.PredicateExists,
		mentalese.PredicateMakeAnd,
		mentalese.PredicateListMake,
	}

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
			newBinding.Set(numberArgumentValue, mentalese.NewTermString(strconv.Itoa(number)))
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
			newBinding.Set(subjectVariable, distinctValues[0])
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

func (base *SystemAggregateBase) makeAnd(input mentalese.Relation, bindings mentalese.Bindings) mentalese.Bindings {

	if !Validate(input, "vv", base.log) {
		return mentalese.Bindings{}
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

	newBindings := mentalese.Bindings{}

	for _, binding := range bindings {
		newBinding := binding.Copy()
		newBinding.Set(andVar, mentalese.NewTermRelationSet(mentalese.RelationSet{ relation }))
		newBindings = append(newBindings, newBinding)
	}

	return newBindings
}

func (base *SystemAggregateBase) listMake(input mentalese.Relation, bindings mentalese.Bindings) mentalese.Bindings {
	if !Validate(input, "V", base.log) {
		return mentalese.Bindings{}
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

	newBindings := mentalese.Bindings{}
	for _, binding := range bindings {
		newBinding := binding.Copy()
		newBinding.Set(listVar, listTerm)
		newBinding = newBinding.FilterOutVariablesByName(variables)
		newBindings = append(newBindings, newBinding)
	}
	return newBindings.UniqueBindings()
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
	case mentalese.PredicateMakeAnd:
		newBindings = base.makeAnd(input, bindings)
	case mentalese.PredicateListMake:
		newBindings = base.listMake(input, bindings)
	default:
		found = false
	}

	return newBindings, found
}
