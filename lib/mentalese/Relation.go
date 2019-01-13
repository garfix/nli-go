package mentalese

import "nli-go/lib/common"

type Relation struct {
	Predicate string
	Arguments []Term
}

const Predicate_Quantification = "quantification"
const Predicate_Temp_Quantification = "temp_quantification"
const Predicate_Quant = "quant"

const PredicateName = "name"
const PredicateSense = "sense"

const Quantification_RangeVariableIndex = 0
const Quantification_RangeIndex = 1
const Quantification_QuantifierVariableIndex = 2
const Quantification_QuantifierIndex = 3
const Quantification_ScopeIndex = 4

func NewRelation(predicate string, arguments []Term) Relation {
	return Relation{
		Predicate: predicate,
		Arguments: arguments,
	}
}

func (relation Relation) GetVariableNames() []string {

	var names []string

	for _, argument := range relation.Arguments {
		if argument.IsVariable() {
			names = append(names, argument.TermValue)
		}
	}

	return common.StringArrayDeduplicate(names)
}

func (relation Relation) Equals(otherRelation Relation) bool {

	equals := relation.Predicate == otherRelation.Predicate

	for i, argument := range relation.Arguments {
		equals = equals && argument.Equals(otherRelation.Arguments[i])
	}

	return equals
}

func (relation Relation) Copy() Relation {

	newRelation := Relation{}
	newRelation.Predicate = relation.Predicate
	newRelation.Arguments = []Term{}
	for _, argument := range relation.Arguments {
		newRelation.Arguments = append(newRelation.Arguments, argument.Copy())
	}
	return newRelation
}

// Returns a new relation, that has all variables bound to bindings
func (relation Relation) BindSingleRelationSingleBinding(binding Binding) Relation {

	boundRelation := Relation{}
	boundRelation.Predicate = relation.Predicate

	for _, argument := range relation.Arguments {

		arg := argument
		if argument.IsVariable() {
			newValue, found := binding[argument.TermValue]
			if found {
				arg = newValue
			}
		} else if argument.IsRelationSet() {
			arg.TermValueRelationSet = argument.TermValueRelationSet.BindRelationSetSingleBinding(binding)
		}

		boundRelation.Arguments = append(boundRelation.Arguments, arg)
	}

	return boundRelation
}

// Returns multiple relations, that has all variables bound to bindings
func (relation Relation) BindSingleRelationMultipleBindings(bindings []Binding) []Relation {

	boundRelations := []Relation{}

	for _, binding := range bindings {
		boundRelations = append(boundRelations, relation.BindSingleRelationSingleBinding(binding))
	}

	return boundRelations
}

func (relation Relation) String() string {

	args, sep := "", ""

	for _, Argument := range relation.Arguments {

		args += sep + Argument.String()
		sep = ", "
	}

	return relation.Predicate + "(" + args + ")"
}
