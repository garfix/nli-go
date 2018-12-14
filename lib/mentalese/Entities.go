package mentalese

type Entities map[string]EntityInfo

type EntityInfo struct {
	Name RelationSet
	Knownby map[string]RelationSet
}

const NameField = "name"
const KnownByField = "knownby"

const NameVar = "Name"
const ValueVar = "Value"
const IdVar = "Id"
