package example3

type SimpleFactBase struct {
	facts []SimpleRelation
}

func NewSimpleFactBase(facts []SimpleRelation) *SimpleFactBase {
	return &SimpleFactBase{facts: facts}
}

func (factBase *SimpleFactBase) Bind(goal SimpleRelation) map[string]SimpleTerm {

}
