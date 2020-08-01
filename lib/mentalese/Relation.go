package mentalese

import (
	"nli-go/lib/common"
)

type Relation struct {
	Positive  bool
	Predicate string
	Arguments []Term
}

const PredicateNone = "none"
const PredicateFind = "find"
const PredicateDo = "do"
const PredicateQuantOrderedList = "quant_ordered_list"
const PredicateQuant = "quant"
const PredicateQuantifier = "quantifier"
const PredicateQuantifierSome = "some"
const PredicateListOrder = "list_order"
const PredicateListForeach = "list_foreach"
const PredicateSem = "sem"
const PredicateAnd = "and"
const PredicateNot = "not"
const PredicateOr = "or"
const PredicateXor = "xor"
const PredicateIfThenElse = "if_then_else"
const PredicateCall = "call"
const PredicateName = "name"
const PredicateText = "text"
const PredicateAssert = "assert"
const PredicateRetract = "retract"
const PredicateIntent = "intent"
const PredicateBackReference = "back_reference"
const PredicateDefiniteBackReference = "definite_reference"
const PredicateQuantOrder = "quant_order"

const QuantifierResultCountVariableIndex = 0
const QuantifierRangeCountVariableIndex = 1
const QuantifierSetIndex = 2

const QuantQuantifierIndex = 0
const QuantRangeVariableIndex = 1
const QuantRangeSetIndex = 2

const SeqFirstOperandIndex = 1
const SeqSecondOperandIndex = 2

const NotScopeIndex = 0

func NewRelation(positive bool, predicate string, arguments []Term) Relation {
	return Relation{
		Positive:  positive,
		Predicate: predicate,
		Arguments: arguments,
	}
}

func (relation Relation) GetVariableNames() []string {

	var names []string

	for _, argument := range relation.Arguments {
		names = append(names, argument.GetVariableNames()...)
	}

	return common.StringArrayDeduplicate(names)
}

func (relation Relation) Equals(otherRelation Relation) bool {

	equals := relation.Predicate == otherRelation.Predicate

	equals = equals && relation.Positive == otherRelation.Positive

	for i, argument := range relation.Arguments {
		equals = equals && argument.Equals(otherRelation.Arguments[i])
	}

	return equals
}

func (relation Relation) Copy() Relation {

	newRelation := Relation{}
	newRelation.Predicate = relation.Predicate
	newRelation.Positive = relation.Positive
	newRelation.Arguments = []Term{}
	for _, argument := range relation.Arguments {
		newRelation.Arguments = append(newRelation.Arguments, argument.Copy())
	}
	return newRelation
}

// Returns a new relation, that has all variables bound to bindings
func (relation Relation) BindSingle(binding Binding) Relation {

	boundArguments := []Term{}

	for _, argument := range relation.Arguments {
		arg := argument.Bind(binding)
		boundArguments = append(boundArguments, arg)
	}

	return NewRelation(relation.Positive, relation.Predicate, boundArguments)
}

// Returns multiple relations, that has all variables bound to bindings
func (relation Relation) BindMultiple(bindings Bindings) []Relation {

	boundRelations := []Relation{}

	for _, binding := range bindings {
		boundRelations = append(boundRelations, relation.BindSingle(binding))
	}

	return boundRelations
}

func (relation Relation) IsBound() bool {
	for _, arg := range relation.Arguments {
		if arg.IsVariable() || arg.IsAnonymousVariable() {
			return false
		}
	}

	return true
}

// check if relation uses a variable (perhaps in one of its nested arguments)
func (relation Relation) UsesVariable(variable string) bool {

	var found = false

	for _, argument := range relation.Arguments {
		found = found || argument.UsesVariable(variable)
	}

	return found
}

func (relation Relation) ConvertVariablesToConstants() Relation {

	newArguments := []Term{}

	for _, argument := range relation.Arguments {

		newArgument := argument.ConvertVariablesToConstants()
		newArguments = append(newArguments, newArgument)
	}

	return NewRelation(relation.Positive, relation.Predicate, newArguments)
}

func (relation Relation) String() string {

	args, sep := "", ""

	for _, Argument := range relation.Arguments {

		args += sep + Argument.String()
		sep = ", "
	}

	sign := ""
	if !relation.Positive {
		sign = "-"
	}

	return sign + relation.Predicate + "(" + args + ")"
}
