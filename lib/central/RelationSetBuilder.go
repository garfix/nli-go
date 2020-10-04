package central

import "nli-go/lib/mentalese"

type RelationSetBuilder struct {
	generators []RelationGenerator
	matcher    *mentalese.RelationMatcher
}

func NewRelationSetBuilder() *RelationSetBuilder {
	return &RelationSetBuilder{}
}

func (builder *RelationSetBuilder) addGenerator(gen RelationGenerator) {
	builder.generators = append(builder.generators, gen)
}

func (builder *RelationSetBuilder) Build(template mentalese.RelationSet, bindings mentalese.Bindings) mentalese.RelationSet {

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

			//if len(bindings) == 0 {
				newSet = append(newSet, templateRelation)
			//} else {
				//relations := templateRelation.BindMultiple(bindings)
				//newSet = newSet.Merge(relations)
			//}

		}
	}

// todo remove ^ that by moving make_and

	template = newSet

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
