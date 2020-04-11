package global

type relationSetFactBase struct {
	Name     string
	Facts    string
	ReadMap  string
	WriteMap string
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
	Grammars           []string
	Rulebases          []string
	Factbases          factBases
	Solutions          []string
	Generationgrammars []string
	Predicates         string
}

func (firstConfig systemConfig) Merge(secondConfig systemConfig) systemConfig {

	predicates := firstConfig.Predicates

	if secondConfig.Predicates != "" {
		predicates = secondConfig.Predicates
	}

	newConfig := systemConfig{
		ParentConfig: secondConfig.ParentConfig,
		Grammars: append(firstConfig.Grammars, secondConfig.Grammars...),
		Rulebases: append(firstConfig.Rulebases, secondConfig.Rulebases...),
		Factbases: factBases {
			Relation: append(firstConfig.Factbases.Relation, secondConfig.Factbases.Relation...),
			Mysql: append(firstConfig.Factbases.Mysql, secondConfig.Factbases.Mysql...),
			Sparql: append(firstConfig.Factbases.Sparql, secondConfig.Factbases.Sparql...),
		},
		Solutions: append(firstConfig.Solutions, secondConfig.Solutions...),
		Generationgrammars: append(firstConfig.Generationgrammars, secondConfig.Generationgrammars...),
		Predicates: predicates,
	}

	return newConfig
}