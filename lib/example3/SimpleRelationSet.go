package example3

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
