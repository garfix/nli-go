package central

import "nli-go/lib/mentalese"

type RelationSetBuilder struct {
	generators []RelationGenerator
	matcher *mentalese.RelationMatcher
}

func NewRelationSetBuilder() *RelationSetBuilder {
	return &RelationSetBuilder{}
}

func (builder *RelationSetBuilder) addGenerator(gen RelationGenerator) {
	builder.generators = append(builder.generators, gen)
}

func (builder *RelationSetBuilder) Build(template mentalese.RelationSet, bindings []mentalese.Binding) mentalese.RelationSet {

	newSet := mentalese.RelationSet{}

	for _, templateRelation := range template {

		generatorUsed := false

		for _, gen := range builder.generators {
			aSet := mentalese.RelationSet{}
			aSet, generatorUsed = gen.generate(templateRelation, bindings)
			if generatorUsed {
				newSet = newSet.Merge(aSet)
				break
			}
		}

		if !generatorUsed {
			relations := builder.matcher.BindSingleRelationMultipleBindings(templateRelation, bindings)
			newSet = newSet.Merge(relations)
		}
	}

	return newSet
}