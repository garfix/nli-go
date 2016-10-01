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
	goalSet := NewSimpleRelationSet()
	goalSet.AddRelation(goal)
	dbRelations := transformer.Extract(goalSet)

// TODO: deze dbRelations hebben nog vrije variabelen die in goal al gebonden waren!

if len(dbRelations.GetRelations()) == 0 {
	fmt.Printf("DB: %v\n", factBase.ds2db)
}

fmt.Printf("Extracted: %v\n", dbRelations);

	matcher := NewSimpleRelationMatcher()

	for _, dbRelation := range dbRelations.GetRelations() {

fmt.Printf("Relation %v\n", dbRelation);

		simpleBinding := SimpleBinding{}
		success := true

		for _, dbFact := range factBase.facts {

			// reset binding
			simpleBinding = SimpleBinding{}

			fmt.Printf("Match %v %v\n", dbRelation, dbFact);

			simpleBinding, success = matcher.matchSubjectToPattern(dbRelation, dbFact, simpleBinding)

			fmt.Printf("Binding %v\n", simpleBinding);

			if success {
				break
			}
		}

		if !success {
			continue
		}

		subgoalRelationSet := []SimpleRelation{}

		subgoalRelationSets = append(subgoalRelationSets, subgoalRelationSet)
		subgoalBindings = append(subgoalBindings, simpleBinding)
	}

fmt.Printf("Factbase end %v %v\n", subgoalRelationSets, subgoalBindings);

	return subgoalRelationSets, subgoalBindings
}
