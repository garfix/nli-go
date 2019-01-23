package mentalese

type Predicates map[string]PredicateInfo

type PredicateInfo struct {
	EntityTypes []string
}

const EntityTypesField = "EntityTypes"