package central

// information about a name in a database, for disambiguation by the user
type NameInformation struct {
	Name         string
	Gender       string
	Number       string
	DatabaseName string
	EntityType   string
	SharedId     string
	Information  string
}
