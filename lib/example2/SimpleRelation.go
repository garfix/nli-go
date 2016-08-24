package example2

type SimpleRelation struct {
	Predicate string
	Arguments []string
}

func (relation *SimpleRelation) ToString() string {

	args, sep := "", ""

	for _, Argument := range relation.Arguments {
		args += sep + Argument
		sep = ", "
	}

	return relation.Predicate + "(" + args + ")"
}
