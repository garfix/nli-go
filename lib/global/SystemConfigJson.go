package global

type relationSetFactBase struct {
	Name     string
	Facts    string
	ReadMap  string
	WriteMap string
	Stats    string
	Entities string
	SharedIds string
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
	Stats	 string
	Entities string
	SharedIds string
	Tables   []mysqlTable
	Enabled  bool
}

type sparqlFactBase struct {
	Name            string
	Baseurl         string
	Defaultgraphuri string
	Map             string
	Names           string
	Stats           string
	Entities        string
	SharedIds 		string
	DoCache         bool
}

type Entities map[string]EntityInfo

type EntityInfo struct {
	Name string
	Knownby map[string]string
}

type factBases struct {
	Relation []relationSetFactBase
	Mysql    []mysqlFactBase
	Sparql   []sparqlFactBase
}

type systemConfig struct {
	ParentConfig       string
	Lexicons           []string
	Grammars           []string
	Rulebases          []string
	Factbases          factBases
	Solutions          []string
	Generationlexicons []string
	Generationgrammars []string
	Generic2ds         []string
	Ds2generic         []string
	Predicates         string
}

func (firstConfig systemConfig) Merge(secondConfig systemConfig) systemConfig {

	predicates := firstConfig.Predicates

	if secondConfig.Predicates != "" {
		predicates = secondConfig.Predicates
	}

	newConfig := systemConfig{
		ParentConfig: secondConfig.ParentConfig,
		Lexicons: append(firstConfig.Lexicons, secondConfig.Lexicons...),
		Grammars: append(firstConfig.Grammars, secondConfig.Grammars...),
		Rulebases: append(firstConfig.Rulebases, secondConfig.Rulebases...),
		Factbases: factBases {
			Relation: append(firstConfig.Factbases.Relation, secondConfig.Factbases.Relation...),
			Mysql: append(firstConfig.Factbases.Mysql, secondConfig.Factbases.Mysql...),
			Sparql: append(firstConfig.Factbases.Sparql, secondConfig.Factbases.Sparql...),
		},
		Solutions: append(firstConfig.Solutions, secondConfig.Solutions...),
		Generationlexicons: append(firstConfig.Generationlexicons, secondConfig.Generationlexicons...),
		Generationgrammars: append(firstConfig.Generationgrammars, secondConfig.Generationgrammars...),
		Generic2ds: append(firstConfig.Generic2ds, secondConfig.Generic2ds...),
		Ds2generic: append(firstConfig.Ds2generic, secondConfig.Ds2generic...),
		Predicates: predicates,
	}

	return newConfig
}