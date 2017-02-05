package mentalese

type Relation struct {
	Predicate string
	Arguments []Term
}

func (relation Relation) Equals(otherRelation Relation) bool {

	equals := relation.Predicate == otherRelation.Predicate

	for i, argument := range relation.Arguments {
		equals = equals && argument	== otherRelation.Arguments[i]
	}

	return equals
}

func (relation Relation) String() string {

	args, sep := "", ""

	for _, Argument := range relation.Arguments {

		term := Argument.TermValue
		if Argument.TermType == Term_stringConstant {
			term = "'" + term + "'"
		}

		args += sep + term
		sep = ", "
	}

	return relation.Predicate + "(" + args + ")"
}
