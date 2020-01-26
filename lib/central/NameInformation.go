package central

// information about a name in a database, for disambiguation by the user
type NameInformation struct {
	Name         string
	DatabaseName string
	EntityId     string
	Information  string
}

func (nameInformation NameInformation) GetIdentifier() string {
	return nameInformation.DatabaseName + "/" + nameInformation.EntityId
}