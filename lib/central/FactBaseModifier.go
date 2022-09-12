package central

import (
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type FactBaseModifier struct {
	matcher           *RelationMatcher
	variableGenerator *mentalese.VariableGenerator
}

func NewFactBaseModifier(log *common.SystemLog, variableGenerator *mentalese.VariableGenerator) *FactBaseModifier {
	return &FactBaseModifier{
		matcher:           NewRelationMatcher(log),
		variableGenerator: variableGenerator,
	}
}

func (modifier FactBaseModifier) Assert(relation mentalese.Relation, factBase api.FactBase) bool {

	found := false

	for _, mapping := range factBase.GetWriteMappings() {

		activeBinding2, match2 := modifier.matcher.MatchTwoRelations(mapping.Goal, relation, mentalese.NewBinding())
		if !match2 {
			continue
		}

		dbRelations := mapping.Pattern.ConvertVariables(activeBinding2, modifier.variableGenerator)

		for _, replacementRelation := range dbRelations {

			factBase.Assert(replacementRelation)
			found = true
		}
	}

	return found
}

func (modifier FactBaseModifier) Retract(relation mentalese.Relation, factBase api.FactBase) bool {

	found := false

	for _, mapping := range factBase.GetWriteMappings() {

		activeBinding2, match2 := modifier.matcher.MatchTwoRelations(mapping.Goal, relation, mentalese.NewBinding())
		if !match2 {
			continue
		}

		dbRelations := mapping.Pattern.ConvertVariables(activeBinding2, modifier.variableGenerator)

		for _, replacementRelation := range dbRelations {

			factBase.Retract(replacementRelation)
			found = true
		}
	}

	return found
}
