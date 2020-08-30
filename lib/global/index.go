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
	Entities        []string
	Shared          []string
	BaseUrl         string
	DefaultGraphUri string
	Names           string
	Cache           bool
	Database		string
	Domain			string
	Username		string
	Password		string
	Tables 			[]table
}

type table struct {
	Name string
	Columns []column
}

type column struct {
	Name string
}