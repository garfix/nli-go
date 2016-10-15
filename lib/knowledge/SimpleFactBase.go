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
func (factBase *SimpleFactBase) Bind(goal mentalese.SimpleRelation) ([]mentalese.SimpleRelationSet, []mentalese.SimpleBinding) {

	common.LogTree("Factbase Bind", goal);

	matcher := mentalese.NewSimpleRelationMatcher()

	subgoalRelationSets := []mentalese.SimpleRelationSet{}
	subgoalBindings := []mentalese.SimpleBinding{}

	for _, ds2db := range factBase.ds2db {

		// gender(14, G), gender(A, male) => externalBinding: G = male
		externalBinding, match := matcher.MatchTwoRelations(goal, ds2db.Goal, mentalese.SimpleBinding{})
		if match {

			// gender(14, G), gender(A, male) => internalBinding: A = 14
			internalBinding, _ := matcher.MatchTwoRelations(ds2db.Goal, goal, mentalese.SimpleBinding{})

			// create a version of the conditions with bound variables
			boundConditions := matcher.BindMultipleRelationsSingleBinding(ds2db.Pattern, internalBinding)
			// match this bound version to the database
			_, internalBinding, match = matcher.MatchSequenceToSet(boundConditions, factBase.facts, mentalese.SimpleBinding{})

			if match {
				subgoalRelationSets = append(subgoalRelationSets, mentalese.SimpleRelationSet{})
				subgoalBindings = append(subgoalBindings, externalBinding.Merge(internalBinding))
			}
		}
	}

	common.LogTree("Factbase Bind", subgoalRelationSets, subgoalBindings);

	return subgoalRelationSets, subgoalBindings
}
