package mentalese

// A Key Cabinet maps variables to database ids for multiple databases
//
// for instance:
//
// entity variable E1 is identified by `311` in database 1, and by `urn:red-phone-90712` in database 2
type KeyCabinet struct {
	Data map[string]map[string]string
}

func NewKeyCabinet() *KeyCabinet {
	return &KeyCabinet{
		Data: map[string]map[string]string{},
	}
}

func (store *KeyCabinet) AddName(variable string, databaseName string, entityId string) {
	_, found := store.Data[databaseName]

	if !found {
		store.Data[databaseName] = map[string]string{}
	}

	store.Data[databaseName][variable] = entityId
}

func (store *KeyCabinet) GetValues(databaseName string) map[string]string {

	values := map[string]string{}

	_, found := store.Data[databaseName]

	if found {
		values = store.Data[databaseName]
	}

	return values
}

func (store *KeyCabinet) BindToRelationSet(set RelationSet, knowledgeBaseName string) RelationSet {

	newSet := RelationSet{}

	databaseValues := store.GetValues(knowledgeBaseName)

	for _, relation := range set {

		newRelation := relation.Copy()

		for i, argument := range relation.Arguments {
			if argument.IsId() {
				entityId, found := databaseValues[argument.TermValue]
				if found {
					newArgument := NewId(entityId)
					newRelation.Arguments[i] = newArgument
				}
			}
		}

		newSet = append(newSet, newRelation)
	}

	return newSet
}

func (store *KeyCabinet) String() string {

	string := ""
	sep := ""

	for databaseName, ids := range store.Data {

		string += sep + "{"

		sep2 := ""
		for variable, entityId := range ids {
			string += sep2 + variable + ": " + entityId
			sep2 = ", "
		}

		string += "}@" + databaseName
		sep = "; "
	}

	return string
}