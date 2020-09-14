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

const PredicateSplit = "go_split"
const PredicateJoin = "go_join"
const PredicateConcat = "go_concat"
const PredicateGreaterThan = "go_greater_than"
const PredicateLessThan = "go_less_than"
const PredicateGreaterThanEquals = "go_greater_than_equals"
const PredicateLessThanEquals = "go_less_than_equals"
const PredicateEquals = "go_equals"
const PredicateNotEquals = "go_not_equals"
const PredicateCompare = "go_compare"
const PredicateUnify = "go_unify"
const PredicateAdd = "go_add"
const PredicateSubtract = "go_subtract"
const PredicateMultiply = "go_multiply"
const PredicateMin = "go_min"
const PredicateDateToday = "go_date_today"
const PredicateDateSubtractYears = "go_date_subtract_years"
const PredicateSem = "go_sem"

const CategoryText = "text"
const CategoryProperNoun = "proper_noun"

const AtomNone = "none"
const AtomSome = "some"

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
