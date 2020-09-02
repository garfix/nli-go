package mentalese

import (
	"nli-go/lib/common"
)

type Relation struct {
	Positive  bool
	Predicate string
	Arguments []Term
}

const PredicateCanned = "go_canned"
const PredicateQuantCheck = "go_quant_check"
const PredicateQuantForeach = "go_quant_foreach"
const PredicateQuantOrderedList = "go_quant_ordered_list"
const PredicateQuant = "go_quant"
const PredicateQuantifier = "go_quantifier"
const PredicateListOrder = "go_list_order"
const PredicateListForeach = "go_list_foreach"
const PredicateAnd = "go_and"
const PredicateNot = "go_not"
const PredicateOr = "go_or"
const PredicateXor = "go_xor"
const PredicateIfThenElse = "go_if_then_else"
const PredicateCall = "go_call"
const PredicateAssert = "go_assert"
const PredicateRetract = "go_retract"
const PredicateIntent = "go_intent"
const PredicateBackReference = "go_back_reference"
const PredicateDefiniteBackReference = "go_definite_reference"
const PredicateNumberOf = "go_number_of"
const PredicateFirst = "go_first"
const PredicateExists = "go_exists"
const PredicateExec = "go_exec"
const PredicateExecResponse = "go_exec_response"
const PredicateMakeAnd = "go_make_and"

const PREDICATE_SPLIT = "go_split"
const PREDICATE_JOIN = "go_join"
const PREDICATE_CONCAT = "go_concat"
const PREDICATE_GREATER_THAN = "go_greater_than"
const PREDICATE_LESS_THAN = "go_less_than"
const PREDICATE_GREATER_THAN_EQUALS = "go_greater_than_equals"
const PREDICATE_LESS_THAN_EQUALS = "go_less_than_equals"
const PREDICATE_EQUALS = "go_equals"
const PREDICATE_NOT_EQUALS = "go_not_equals"
const PREDICATE_COMPARE = "go_compare"
const PREDICATE_UNIFY = "go_unify"
const PREDICATE_ADD = "go_add"
const PREDICATE_SUBTRACT = "go_subtract"
const PREDICATE_MIN = "go_min"
const PREDICATE_DATE_TODAY = "go_date_today"
const PREDICATE_DATE_SUBTRACT_YEARS = "go_date_subtract_years"

const PredicateNone = "none"
const PredicateQuantifierSome = "some"
const PredicateName = "name"
const PredicateSem = "sem"
const PredicateText = "text"
const PredicateProperNoun = "proper_noun"

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
