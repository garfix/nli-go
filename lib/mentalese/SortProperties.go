package mentalese

type SortProperties map[string]SortProperty

type SortProperty struct {
	Name    Relation
	Gender  Relation
	Number  Relation
	Knownby map[string]Relation
	Entity  RelationSet
}

const NameVar = "Name"
const ValueVar = "Value"
const IdVar = "Id"
