package global

type relationSetFactBase struct {
	Facts string
	Map   string
}

type mysqlColumn struct {
	Name string
}

type mysqlTable struct {
	Name    string
	Columns []mysqlColumn
}

type mysqlFactBase struct {
	Domain   string
	Username string
	Password string
	Database string
	Map      string
	Tables   []mysqlTable
	Enabled  bool
}

type factBases struct {
	Relation []relationSetFactBase
	Mysql    []mysqlFactBase
}

type systemConfig struct {
	Lexicons           []string
	Grammars           []string
	Rulebases          []string
	Factbases          factBases
	Solutions          []string
	Generationlexicons []string
	Generationgrammars []string
	Generic2ds         []string
	Ds2generic         []string
}
