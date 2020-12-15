package global

type index struct {
	Type     string
	Version         string
	Uses            map[string]string
	Read            []string
	Write           []string
	Solution        []string
	Rules           []string
	Facts           []string
	Shared          []string
	BaseUrl         string
	DefaultGraphUri string
	Names           string
	Cache           bool
	Database        string
	Username        string
	Password        string
	Tables          []table
	Sorts           string
	Predicates      string
	Subsorts        string
	TokenExpression string
	Morphology		map[string]string
}

type table struct {
	Name string
	Columns []column
}

type column struct {
	Name string
}

type Entities map[string]EntityInfo

type EntityInfo struct {
	Name string
	Knownby map[string]string
	Entity string
}
