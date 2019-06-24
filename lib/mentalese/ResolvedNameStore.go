package mentalese

// stores and retrieves database id's for names
type KeyCabinet struct {
	data map[string]map[string]string
}

func NewKeyCabinet() *KeyCabinet {
	return &KeyCabinet{
		data: map[string]map[string]string{},
	}
}

func (store *KeyCabinet) IsEmpty() bool {
	return len(store.data) == 0
}

func (store *KeyCabinet) AddName(variable string, databaseName string, entityId string) {
	_, found := store.data[databaseName]

	if !found {
		store.data[databaseName] = map[string]string{}
	}

	store.data[databaseName][variable] = entityId
}

func (store *KeyCabinet) GetValues(databaseName string) map[string]string {

	values := map[string]string{}

	_, found := store.data[databaseName]

	if found {
		values = store.data[databaseName]
	}

	return values
}

func (store *KeyCabinet) ReplaceVariables(oldStore *KeyCabinet, binding Binding) *KeyCabinet {

	newStore := NewKeyCabinet()

	for knowledgeBaseName, values := range oldStore.data {
		for oldVariable, value := range values {

			found := false

			for newVariable, term := range binding {
				if term.TermValue == oldVariable {
					found = true
					newStore.AddName(newVariable, knowledgeBaseName, value)
					break
				}
			}

			if !found {
				newStore.AddName(oldVariable, knowledgeBaseName, value)
			}
		}
	}

	return newStore
}

func (store *KeyCabinet) BindToRelationSet(set RelationSet, knowledgeBaseName string) RelationSet {

	newSet := RelationSet{}

	databaseValues := store.GetValues(knowledgeBaseName)

	for _, relation := range set {

		newRelation := relation.Copy()

		for i, argument := range relation.Arguments {
			if argument.IsVariable() {
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

	for databaseName, ids := range store.data {

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