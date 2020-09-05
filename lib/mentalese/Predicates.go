package mentalese

// maps a predicate to information about a relation
type Predicates map[string]PredicateInfo

// for each argument an entity type
type PredicateInfo struct {
	EntityTypes []string
}

func (predicates Predicates) AddPredicates(p Predicates) {
	for predicate, info := range p {
		predicates[predicate] = info
	}
}

func (predicates Predicates) GetEntityType(predicate string, argumentIndex int) string {

	pred, found := predicates[predicate]
	if found {
		return pred.EntityTypes[argumentIndex]
	}
	return ""
}