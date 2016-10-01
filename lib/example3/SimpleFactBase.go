package example3

import "fmt"

type SimpleFactBase struct {
	facts []SimpleRelation
	ds2db []SimpleRule
}

func NewSimpleFactBase(facts []SimpleRelation, ds2db []SimpleRule) *SimpleFactBase {
	return &SimpleFactBase{facts: facts, ds2db: ds2db}
}

// Note! An internal fact base would use the same predicates as the domain language;
// This is an simulation of an external database
func (factBase *SimpleFactBase) Bind(goal SimpleRelation) ([][]SimpleRelation, []SimpleBinding) {

fmt.Printf("Factbase start %v\n", goal);

	subgoalRelationSets := [][]SimpleRelation{}
	subgoalBindings := []SimpleBinding{}

	transformer := NewSimpleRelationTransformer2(factBase.ds2db)

fmt.Printf("DB: %v\n", factBase.ds2db)

	dbRelationSets, dbBindings := transformer.Extract([]SimpleRelation{goal})

// bij het extracten moet je bijhouden aan hoe de oorspronkelijke variabelen gebonden worden
// gebruik deze bindings ook als defaults hieronder

fmt.Printf("Extracted: %v %v\n", dbRelationSets, dbBindings);

	matcher := NewSimpleRelationMatcher()
	newSimpleBinding := SimpleBinding{}

	for i, dbRelationSet := range dbRelationSets {

		simpleBinding := SimpleBinding{}
		relationsFound := true

		for _, dbRelation := range dbRelationSet {

			fmt.Printf("Relation %v\n", dbRelation);

			factFound := false

			for _, dbFact := range factBase.facts {

				fmt.Printf("Match %v %v %s\n", dbRelation, dbFact, simpleBinding);

				newSimpleBinding, factFound = matcher.matchSubjectToPattern(dbRelation, dbFact, simpleBinding)

				fmt.Printf("Binding %v %b\n", newSimpleBinding, factFound);

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

fmt.Printf("Relations found %b\n", relationsFound);

		if relationsFound {
			subgoalRelationSet := []SimpleRelation{}

			subgoalRelationSets = append(subgoalRelationSets, subgoalRelationSet)
			subgoalBindings = append(subgoalBindings, dbBindings[i].Merge(simpleBinding))
		}
	}

fmt.Printf("Factbase end %v %v\n", subgoalRelationSets, subgoalBindings);

	return subgoalRelationSets, subgoalBindings
}
