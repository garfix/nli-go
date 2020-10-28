package central

import (
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type FactBaseModifier struct {
	matcher *RelationMatcher
}

func NewFactBaseModifier(log *common.SystemLog) *FactBaseModifier {
	return &FactBaseModifier{
		matcher: NewRelationMatcher(log),
	}
}

func (modifier FactBaseModifier) Assert(relation mentalese.Relation, factBase api.FactBase) {

	for _, mapping := range factBase.GetWriteMappings() {

		activeBinding2, match2 := modifier.matcher.MatchTwoRelations(mapping.Goal, relation, mentalese.NewBinding())
		if !match2 { continue }

		dbRelations := mapping.Pattern.ImportBinding(activeBinding2)

		for _, replacementRelation := range dbRelations {

			factBase.Assert(replacementRelation)
		}
	}
}

func (modifier FactBaseModifier) Retract(relation mentalese.Relation, factBase api.FactBase) {

	for _, mapping := range factBase.GetWriteMappings() {

		activeBinding2, match2 := modifier.matcher.MatchTwoRelations(mapping.Goal, relation, mentalese.NewBinding())
		if !match2 { continue }

		dbRelations := mapping.Pattern.ImportBinding(activeBinding2)

		for _, replacementRelation := range dbRelations {

			factBase.Retract(replacementRelation)
		}
	}
}
