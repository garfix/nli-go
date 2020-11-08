package mentalese

type Entities map[string]SortInfo

type SortInfo struct {
	Name Relation
	Knownby map[string]Relation
	Entity RelationSet
}

const NameVar = "Name"
const ValueVar = "Value"
const IdVar = "Id"
