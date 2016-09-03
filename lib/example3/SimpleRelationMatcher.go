package example3

type simpleRelationMatcher struct {

}

func NewSimpleRelationMatcher() *simpleRelationMatcher {
	return &simpleRelationMatcher{}
}

func (matcher *simpleRelationMatcher) Match(pattern *SimpleRelationSet, subject *SimpleRelationSet) bool {
	matchedIndexes, _ := matchRelations(subject.relations, pattern.relations)
	return len(matchedIndexes) > 0
}

