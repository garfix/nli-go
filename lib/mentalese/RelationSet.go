package mentalese

type RelationSet []Relation

func (set RelationSet) String() string {

	s, sep := "", ""

	for _, relation := range set {
		s += sep + relation.String()
		sep = " "
	}

	return "[" + s + "]";
}