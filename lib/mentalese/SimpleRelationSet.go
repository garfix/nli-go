package mentalese

type SimpleRelationSet []SimpleRelation

func (set SimpleRelationSet) String() string {

	s, sep := "", ""

	for _, relation := range set {
		s += sep + relation.String()
		sep = " "
	}

	return "[" + s + "]";
}