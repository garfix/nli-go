package mentalese

type Entities map[string]EntityInfo

type EntityInfo struct {
	Name Relation
	Knownby map[string]Relation
	Entity RelationSet
}

const NameVar = "Name"
const ValueVar = "Value"
const IdVar = "Id"
