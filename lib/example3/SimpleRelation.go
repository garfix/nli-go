package example3

type SimpleRelation struct {
	Predicate string
	Arguments []SimpleTerm
}

func (relation SimpleRelation) String() string {

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
