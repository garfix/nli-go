package knowledge

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
)

type FactBase struct {
	facts []mentalese.Relation
	ds2db []mentalese.Rule
}

func NewFactBase(facts []mentalese.Relation, ds2db []mentalese.Rule) *FactBase {
	return &FactBase{facts: facts, ds2db: ds2db}
}

// Note! An internal fact base would use the same predicates as the domain language;
// This is an simulation of an external database
func (factBase *FactBase) Bind(goal mentalese.Relation) ([]mentalese.RelationSet, []mentalese.Binding) {

	common.LogTree("Factbase Bind", goal);

	matcher := mentalese.NewRelationMatcher()

	subgoalRelationSets := []mentalese.RelationSet{}
	subgoalBindings := []mentalese.Binding{}

	for _, ds2db := range factBase.ds2db {

		// gender(14, G), gender(A, male) => externalBinding: G = male
		externalBinding, match := matcher.MatchTwoRelations(goal, ds2db.Goal, mentalese.Binding{})
		if match {

			// gender(14, G), gender(A, male) => internalBinding: A = 14
			internalBinding, _ := matcher.MatchTwoRelations(ds2db.Goal, goal, mentalese.Binding{})

			// create a version of the conditions with bound variables
			boundConditions := matcher.BindMultipleRelationsSingleBinding(ds2db.Pattern, internalBinding)
			// match this bound version to the database
			_, internalBinding, match = matcher.MatchSequenceToSet(boundConditions, factBase.facts, mentalese.Binding{})

			if match {
				subgoalRelationSets = append(subgoalRelationSets, mentalese.RelationSet{})
				subgoalBindings = append(subgoalBindings, externalBinding.Merge(internalBinding))
			}
		}
	}

	common.LogTree("Factbase Bind", subgoalRelationSets, subgoalBindings);

	return subgoalRelationSets, subgoalBindings
}
