package mentalese

// A Key Cabinet a variables to a shared id (which defaults to just a database id)
//
// for instance:
//
// constant ce1 is mapped to shared id `311` which maps to `p100022`  in database 1, and to `urn:red-phone-90712` in database 2
type KeyCabinet struct {
	Data map[string]string
}

func NewKeyCabinet() *KeyCabinet {
	return &KeyCabinet{
		Data: map[string]string{},
	}
}

func (store *KeyCabinet) AddMapping(variable string, entityId string) {
	store.Data[variable] = entityId
}

//func (store *KeyCabinet) getValues(databaseName string) map[string]string {
//
//	values := map[string]string{}
//
//	_, found := store.Data[databaseName]
//
//	if found {
//		values = store.Data[databaseName]
//	}
//
//	return values
//}

//func (store *KeyCabinet) BindToRelationSet(set RelationSet, knowledgeBaseName string) RelationSet {
//
//	newSet := RelationSet{}
//
//	databaseValues := store.getValues(knowledgeBaseName)
//
//	for _, relation := range set {
//
//		newRelation := relation.Copy()
//
//		for i, argument := range relation.Arguments {
//			if argument.IsId() {
//				entityId, found := databaseValues[argument.TermValue]
//				if found {
//					newArgument := NewId(entityId)
//					newRelation.Arguments[i] = newArgument
//				}
//			}
//		}
//
//		newSet = append(newSet, newRelation)
//	}
//
//	return newSet
//}

func (store *KeyCabinet) String() string {

	string := ""
	sep := ""

	for variable, id := range store.Data {

		string += sep + variable + ": " + id
		sep = "; "
	}

	return string
}