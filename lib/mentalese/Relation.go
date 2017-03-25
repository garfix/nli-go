package mentalese

type Relation struct {
	Predicate string
	Arguments []Term
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
