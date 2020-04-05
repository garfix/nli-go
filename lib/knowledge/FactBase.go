package knowledge

import (
	"nli-go/lib/mentalese"
)

type FactBase interface {
	KnowledgeBase
	MatchRelationToDatabase(needleRelation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings
	Assert(relation mentalese.Relation)
	Retract(relation mentalese.Relation)
	GetMappings() []mentalese.RelationTransformation
	GetWriteMappings() []mentalese.RelationTransformation
	GetEntities() mentalese.Entities
	GetLocalId(sharedId string, entityType string) string
	GetSharedId(localId string, entityType string) string
}

func getFactBaseMatchingGroups(matcher *mentalese.RelationMatcher, set mentalese.RelationSet, factBase FactBase) []RelationGroup {

	matchingGroups := []RelationGroup{}

	matchingGroups = append(matchingGroups, getFactBaseReadGroups(matcher, set, factBase)...)

	matchingGroups = append(matchingGroups, getFactBaseWriteGroups(matcher, set, factBase)...)
	matchingGroups = append(matchingGroups, getFactBaseWriteGroups(matcher, set, factBase)...)

	return matchingGroups
}

func getFactBaseReadGroups(matcher *mentalese.RelationMatcher, set mentalese.RelationSet, factBase FactBase) []RelationGroup {

	matchingGroups := []RelationGroup{}

	for _, mapping := range factBase.GetMappings() {

		bindings, _, indexesPerNode, match := matcher.MatchSequenceToSetWithIndexes(mapping.Pattern, set, mentalese.Binding{})

		if match {

			for i := range bindings {

				indexes := indexesPerNode[i].Indexes

				matchingRelations := mentalese.RelationSet{}
				for _, i := range indexes {
					matchingRelations = append(matchingRelations, set[i])
				}

				matchingGroups = append(matchingGroups, RelationGroup{matchingRelations, factBase.GetName()})
			}
		}
	}

	return matchingGroups
}

func getFactBaseWriteGroups(matcher *mentalese.RelationMatcher, set mentalese.RelationSet, factBase FactBase) []RelationGroup {

	matchingGroups := []RelationGroup{}

	for _, relation := range set {
		if relation.Predicate == mentalese.PredicateAssert || relation.Predicate == mentalese.PredicateRetract {
			content := relation.Arguments[0].TermValueRelationSet

			for _, mapping := range factBase.GetWriteMappings() {

				_, _, indexesPerNode, match := matcher.MatchSequenceToSetWithIndexes(mapping.Pattern, content, mentalese.Binding{})

				if match {

					indexes := indexesPerNode[0].Indexes

					matchingRelations := mentalese.RelationSet{}
					for _, i := range indexes {
						matchingRelations = append(matchingRelations, set[i])
					}

					matchingGroups = append(matchingGroups, RelationGroup{matchingRelations, factBase.GetName()})
				}
			}
		}
	}

	return matchingGroups
}
