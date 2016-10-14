package knowledge

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
)

type SimpleFactBase struct {
	facts []mentalese.SimpleRelation
	ds2db []mentalese.SimpleRule
}

func NewSimpleFactBase(facts []mentalese.SimpleRelation, ds2db []mentalese.SimpleRule) *SimpleFactBase {
	return &SimpleFactBase{facts: facts, ds2db: ds2db}
}

// Note! An internal fact base would use the same predicates as the domain language;
// This is an simulation of an external database
func (factBase *SimpleFactBase) Bind(goal mentalese.SimpleRelation) ([][]mentalese.SimpleRelation, []mentalese.SimpleBinding) {

	common.LogTree("Factbase Bind", goal);

	matcher := mentalese.NewSimpleRelationMatcher()

	subgoalRelationSets := []mentalese.SimpleRelationSet{}
	subgoalBindings := []mentalese.SimpleBinding{}

	for _, ds2db := range factBase.ds2db {

		externalBinding := mentalese.SimpleBinding{}
		match := false

		externalBinding, match = matcher.MatchTwoRelations(goal, ds2db.Goal, externalBinding)
		if match {

			transBinding := mentalese.SimpleBinding{}
			transBinding, match = matcher.MatchTwoRelations(goal, ds2db.Goal, transBinding)

			_, internalBinding, match := matcher.MatchSequenceToSet(ds2db.Pattern, factBase.facts, mentalese.SimpleBinding{})

			if match {
				subgoalRelationSets = append(subgoalRelationSets, []mentalese.SimpleRelation{})
				subgoalBindings = append(subgoalBindings, externalBinding.Bind(internalBinding))
			}
		}
	}

	subgoalRelationSets3 := [][]mentalese.SimpleRelation{}
// TODO: ieuw!
	for _, r := range subgoalRelationSets {
		subgoalRelationSets3 = append(subgoalRelationSets3, r)
	}

	common.LogTree("Factbase Bind", subgoalRelationSets3, subgoalBindings);

	return subgoalRelationSets3, subgoalBindings
}
