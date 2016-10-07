package process

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

	common.Logf("Factbase start %v\n", goal);

	subgoalRelationSets := [][]mentalese.SimpleRelation{}
	subgoalBindings := []mentalese.SimpleBinding{}

	transformer := NewSimpleRelationTransformer2(factBase.ds2db)

	common.Logf("DB: %v\n", factBase.ds2db)

	dbRelationSets, dbBindings := transformer.Extract([]mentalese.SimpleRelation{goal})

// bij het extracten moet je bijhouden aan hoe de oorspronkelijke variabelen gebonden worden
// gebruik deze bindings ook als defaults hieronder

	common.Logf("Extracted: %v %v\n", dbRelationSets, dbBindings);

	matcher := NewSimpleRelationMatcher()
	newSimpleBinding := mentalese.SimpleBinding{}

	for i, dbRelationSet := range dbRelationSets {

		simpleBinding := mentalese.SimpleBinding{}
		relationsFound := true

		for _, dbRelation := range dbRelationSet {

			common.Logf("Relation %v\n", dbRelation);

			factFound := false

			for _, dbFact := range factBase.facts {

				common.Logf("Match %v %v %s\n", dbRelation, dbFact, simpleBinding);

				newSimpleBinding, factFound = matcher.MatchSubjectToPattern(dbRelation, dbFact, simpleBinding)

				common.Logf("Binding %v %b\n", newSimpleBinding, factFound);

				if factFound {
					simpleBinding = newSimpleBinding
					break
				}
			}

			if !factFound {
				relationsFound = false
				break
			}
		}

		common.Logf("Relations found %b\n", relationsFound);

		if relationsFound {
			subgoalRelationSet := []mentalese.SimpleRelation{}

			subgoalRelationSets = append(subgoalRelationSets, subgoalRelationSet)
			subgoalBindings = append(subgoalBindings, dbBindings[i].Merge(simpleBinding))
		}
	}

	common.Logf("Factbase end %v %v\n", subgoalRelationSets, subgoalBindings);

	return subgoalRelationSets, subgoalBindings
}
