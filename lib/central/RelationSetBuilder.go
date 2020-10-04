package central

import "nli-go/lib/mentalese"

type RelationSetBuilder struct {
	matcher    *mentalese.RelationMatcher
}

func NewRelationSetBuilder() *RelationSetBuilder {
	return &RelationSetBuilder{}
}

func (builder *RelationSetBuilder) Build(template mentalese.RelationSet, bindings mentalese.Bindings) mentalese.RelationSet {

	newSet := mentalese.RelationSet{}

	if len(bindings) == 0 {
		newSet = template
	} else {

		sets := template.BindRelationSetMultipleBindings(bindings)

		newSet = mentalese.RelationSet{}
		for _, set := range sets {
			newSet = newSet.Merge(set)
		}
	}

	return newSet
}
