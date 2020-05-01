package mentalese

import (
	"nli-go/lib/common"
	"strings"
)

type Relation struct {
	Predicate string
	Arguments []Term
}

const PredicateFind = "find"
const PredicateDo = "do"
const PredicateQuant = "quant"
const PredicateSem = "sem"
const PredicateAnd = "and"
const PredicateNot = "not"
const PredicateOr = "or"
const PredicateCall = "call"
const PredicateName = "name"
const PredicateText = "text"

const PredicateAssert = "assert"
const PredicateRetract = "retract"

const PredicateBackReference = "back_reference"
const PredicateDefiniteBackReference = "definite_reference"

const QuantResultCountVariableIndex = 0
const QuantRangeCountVariableIndex = 1
const QuantQuantifierSetIndex = 2
const QuantRangeVariableIndex = 3
const QuantRangeSetIndex = 4

const SeqFirstOperandIndex = 1
const SeqSecondOperandIndex = 2

const NotScopeIndex = 0

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
		} else if argument.IsRelationSet() {
			names = append(names, argument.TermValueRelationSet.GetVariableNames()...)
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
func (relation Relation) BindSingle(binding Binding) Relation {

	boundRelation := Relation{}
	boundRelation.Predicate = relation.Predicate

	for _, argument := range relation.Arguments {
		arg := argument.Bind(binding)
		boundRelation.Arguments = append(boundRelation.Arguments, arg)
	}

	return boundRelation
}

// Returns multiple relations, that has all variables bound to bindings
func (relation Relation) BindMultiple(bindings Bindings) []Relation {

	boundRelations := []Relation{}

	for _, binding := range bindings {
		boundRelations = append(boundRelations, relation.BindSingle(binding))
	}

	return boundRelations
}

// check if relation uses a variable (perhaps in one of its nested arguments)
func (relation Relation) RelationUsesVariable(variable string) bool {

	var found = false

	for _, argument := range relation.Arguments {
		if argument.IsVariable() {
			found = found || argument.TermValue == variable
		} else if argument.IsRelationSet() {
			for _, rel := range argument.TermValueRelationSet {
				found = found || rel.RelationUsesVariable(variable)
			}
		}
	}

	return found
}

func (relation Relation) ConvertVariablesToConstants() Relation {

	newRelation := Relation{ Predicate: relation.Predicate }

	for _, argument := range relation.Arguments {

		newArgument := argument.Copy()

		if argument.IsVariable() {
			newArgument = NewPredicateAtom(strings.ToLower(argument.TermValue))
		} else if argument.IsRelationSet() {
			newArgument = NewRelationSet(argument.TermValueRelationSet.ConvertVariablesToConstants())
		}

		newRelation.Arguments = append(newRelation.Arguments, newArgument)
	}

	return newRelation
}

func (relation Relation) String() string {

	args, sep := "", ""

	for _, Argument := range relation.Arguments {

		args += sep + Argument.String()
		sep = ", "
	}

	return relation.Predicate + "(" + args + ")"
}
