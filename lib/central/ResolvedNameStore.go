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

func (store *ResolvedNameStore) AddName(variable string, databaseName string , entityId string) {
	_, found := store.data[variable]

	if !found {
		store.data[variable] = map[string]string{}
	}

	store.data[variable][databaseName] = entityId
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