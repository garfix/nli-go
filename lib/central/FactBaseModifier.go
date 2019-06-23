package central

import (
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
)

type FactBaseModifier struct {
	matcher *mentalese.RelationMatcher
}

func NewFactBaseModifier(log *common.SystemLog) *FactBaseModifier {
	return &FactBaseModifier{
		matcher: mentalese.NewRelationMatcher(log),
	}
}

func (modifier FactBaseModifier) Assert(set mentalese.RelationSet, factBase knowledge.FactBase, nameStore *mentalese.ResolvedNameStore) {

	for _, mapping := range factBase.GetWriteMappings() {

		bindings, _, indexesPerNode, match := modifier.matcher.MatchSequenceToSetWithIndexes(mapping.Pattern, set, mentalese.Binding{})

		if match {

			binding := bindings[0]
			indexes := indexesPerNode[0].Indexes

			matchingRelations := mentalese.RelationSet{}
			for _, i := range indexes {
				matchingRelations = append(matchingRelations, set[i])
			}

			boundReplacement := mapping.Replacement.BindSingle(binding)

			keyBoundReplacement := nameStore.BindToRelationSet(boundReplacement, factBase.GetName())

			for _, replacementRelation := range keyBoundReplacement {

				factBase.Assert(replacementRelation)
			}
		}
	}
}

func (modifier FactBaseModifier) Retract(set mentalese.RelationSet, factBase knowledge.FactBase, nameStore *mentalese.ResolvedNameStore) {

	for _, mapping := range factBase.GetWriteMappings() {

		bindings, _, indexesPerNode, match := modifier.matcher.MatchSequenceToSetWithIndexes(mapping.Pattern, set, mentalese.Binding{})

		if match {

			binding := bindings[0]
			indexes := indexesPerNode[0].Indexes

			matchingRelations := mentalese.RelationSet{}
			for _, i := range indexes {
				matchingRelations = append(matchingRelations, set[i])
			}

			boundReplacement := mapping.Replacement.BindSingle(binding)

			keyBoundReplacement := nameStore.BindToRelationSet(boundReplacement, factBase.GetName())

			for _, replacementRelation := range keyBoundReplacement {

				factBase.Retract(replacementRelation)
			}
		}
	}
}
