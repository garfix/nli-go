package central

// stores and retrieves database id's for names
type ResolvedNameStore struct {
	data map[string]map[string]string
}

func NewResolvedNameStore() *ResolvedNameStore {
	return &ResolvedNameStore{
		data: map[string]map[string]string{},
	}
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