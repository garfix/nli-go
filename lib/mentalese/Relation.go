package mentalese

type Relation struct {
	Predicate string
	Arguments []Term
}

const Predicate_Quantification = "quantification"
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

func (relation Relation) String() string {

	args, sep := "", ""

	for _, Argument := range relation.Arguments {

		args += sep + Argument.String()
		sep = ", "
	}

	return relation.Predicate + "(" + args + ")"
}
