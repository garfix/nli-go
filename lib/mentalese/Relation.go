package mentalese

type Relation struct {
	Predicate string
	Arguments []Term
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
