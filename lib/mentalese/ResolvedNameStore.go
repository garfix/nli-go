package mentalese

// stores and retrieves database id's for names
type ResolvedNameStore struct {
	data map[string]map[string]string
}

func NewResolvedNameStore() *ResolvedNameStore {
	return &ResolvedNameStore{
		data: map[string]map[string]string{},
	}
}

func (store *ResolvedNameStore) IsEmpty() bool {
	return len(store.data) == 0
}

func (store *ResolvedNameStore) AddName(variable string, databaseName string, entityId string) {
	_, found := store.data[databaseName]

	if !found {
		store.data[databaseName] = map[string]string{}
	}

	store.data[databaseName][variable] = entityId
}

func (store *ResolvedNameStore) GetValues(databaseName string) map[string]string {

	values := map[string]string{}

	_, found := store.data[databaseName]

	if found {
		values = store.data[databaseName]
	}

	return values
}

func (store *ResolvedNameStore) ReplaceVariables(oldStore *ResolvedNameStore, binding Binding) *ResolvedNameStore {

	newStore := NewResolvedNameStore()

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

func (store *ResolvedNameStore) BindToRelationSet(set RelationSet, knowledgeBaseName string) RelationSet {

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

func (store *ResolvedNameStore) String() string {

	string := ""
	sep := ""

	for variable, ids := range store.data {

		string += sep + variable + ": "

		sep2 := ""
		for databaseName, entityId := range ids {
			string += sep2 + databaseName + " = " + entityId
			sep2 = ", "
		}

		sep = "; "
	}

	return string
}