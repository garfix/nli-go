package example3

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

	subgoalRelationSets := [][]SimpleRelation{}
	subgoalBindings := []SimpleBinding{}

	transformer := NewSimpleRelationTransformer2(factBase.ds2db)
	dbRelations := transformer.Extract([]SimpleRelation{goal})

	matcher := NewSimpleRelationMatcher()

	for _, dbRelation := range dbRelations.GetRelations() {

		simpleBinding := SimpleBinding{}
		success := true

		for _, dbFact := range factBase.facts {
			simpleBinding, success = matcher.matchRelationToRelation(dbRelation, dbFact, simpleBinding)
			if !success {
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

	return subgoalRelationSets, subgoalBindings
}
