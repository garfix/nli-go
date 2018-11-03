package central

type NameInformation struct {
	Name         string
	DatabaseName string
	EntityId     string
	Information  string
}

func (nameInformation NameInformation) GetIdentifier() string {
	return nameInformation.DatabaseName + "/" + nameInformation.EntityId
}