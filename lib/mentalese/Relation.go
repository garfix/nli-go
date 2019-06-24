package mentalese

import "nli-go/lib/common"

type Relation struct {
	Predicate string
	Arguments []Term
}

const PredicateQuantification = "quantification"
const PredicateQuant = "quant"

const PredicateSequence = "sequence"
const PredicateSeq = "seq"

const PredicateName = "name"
const PredicateSense = "sense"

const PredicateAssert = "assert"
const PredicateRetract = "retract"

const PredicateNumber = "number"

const QuantificationRangeVariableIndex = 0
const QuantificationRangeIndex = 1
const QuantificationQuantifierVariableIndex = 2
const QuantificationQuantifierIndex = 3
const QuantificationScopeIndex = 4

const SeqFirstOperandIndex = 0;
const SeqSecondOperandIndex = 2;

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
			arg.TermValueRelationSet = argument.TermValueRelationSet.BindSingle(binding)
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

func (relation Relation) String() string {

	args, sep := "", ""

	for _, Argument := range relation.Arguments {

		args += sep + Argument.String()
		sep = ", "
	}

	return relation.Predicate + "(" + args + ")"
}
