package example3

import "nli-go/lib/example2"

type simpleRelationTransformer struct {
	transformations []SimpleRelationTransformation
}

func NewSimpleRelationTransformer(transformations[]SimpleRelationTransformation) *simpleRelationTransformer {
	return &simpleRelationTransformer{transformations: transformations}
}

func (transformer *simpleRelationTransformer) Process(relations []example2.SimpleRelation) []example2.SimpleRelation {
	return []example2.SimpleRelation{}
}
