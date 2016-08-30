package example3

type SimpleRelation struct {
	Predicate string
	Arguments []SimpleTerm
}

func (relation *SimpleRelation) ToString() string {

	args, sep := "", ""

	for _, Argument := range relation.Arguments {
		args += sep + Argument.TermValue
		sep = ", "
	}

	return relation.Predicate + "(" + args + ")"
}
