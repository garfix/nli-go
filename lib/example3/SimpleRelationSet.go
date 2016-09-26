package example3

// TODO type SimpleRelationSet relations []SimpleRelation
type SimpleRelationSet struct {
	relations []SimpleRelation
}

func NewSimpleRelationSet() *SimpleRelationSet {
	return &SimpleRelationSet{}
}

func NewSimpleRelationSet2(relations []SimpleRelation) *SimpleRelationSet {
	return &SimpleRelationSet{relations: relations}
}

func (set *SimpleRelationSet) AddRelation(relation SimpleRelation) {
	set.relations = append(set.relations, relation)
}

func (set *SimpleRelationSet) AddRelations(relations []SimpleRelation) {
	set.relations = append(set.relations, relations...)
}

func (set *SimpleRelationSet) GetRelations() []SimpleRelation {
	return set.relations
}

func (set *SimpleRelationSet) String() string {

	s, sep := "", ""

	for _, relation := range set.relations {
		s += sep + relation.String()
		sep = " "
	}

	return "[" + s + "]";
}