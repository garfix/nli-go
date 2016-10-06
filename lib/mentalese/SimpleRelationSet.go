package mentalese

// TODO type mentalese.SimpleRelationSet relations []mentalese.SimpleRelation
type SimpleRelationSet struct {
	Relations []SimpleRelation
}

func NewSimpleRelationSet() *SimpleRelationSet {
	return &SimpleRelationSet{}
}

func NewSimpleRelationSet2(relations []SimpleRelation) *SimpleRelationSet {
	return &SimpleRelationSet{Relations: relations}
}

func (set *SimpleRelationSet) AddRelation(relation SimpleRelation) {
	set.Relations = append(set.Relations, relation)
}

func (set *SimpleRelationSet) AddRelations(relations []SimpleRelation) {
	set.Relations = append(set.Relations, relations...)
}

func (set *SimpleRelationSet) GetRelations() []SimpleRelation {
	return set.Relations
}

func (set *SimpleRelationSet) String() string {

	s, sep := "", ""

	for _, relation := range set.Relations {
		s += sep + relation.String()
		sep = " "
	}

	return "[" + s + "]";
}