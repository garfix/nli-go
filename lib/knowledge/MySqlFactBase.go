package knowledge

// https://dinosaurscode.xyz/go/2016/06/19/golang-mysql-authentication/
// go get github.com/go-sql-driver/mysql

import (
	"nli-go/lib/mentalese"
)
import "database/sql"
import (
	_ "github.com/go-sql-driver/mysql"
	"nli-go/lib/common"
	"strings"
)

type MySqlFactBase struct {
	KnowledgeBaseCore
	db                *sql.DB
	tableDescriptions map[string][]string
	ds2db             []mentalese.RelationTransformation
	stats			  mentalese.DbStats
	entities 		  mentalese.Entities
	matcher           *mentalese.RelationMatcher
	log               *common.SystemLog
}

func NewMySqlFactBase(name string, domain string, username string, password string, database string, matcher *mentalese.RelationMatcher, ds2db []mentalese.RelationTransformation, stats mentalese.DbStats, entities mentalese.Entities, log *common.SystemLog) *MySqlFactBase {

	db, err := sql.Open("mysql", username+":"+password+"@/"+database)
	if err != nil {
		log.AddError("Error opening MySQL: " + err.Error())
	}

	return &MySqlFactBase{
		KnowledgeBaseCore: KnowledgeBaseCore{ Name: name},
		db: db,
		tableDescriptions: map[string][]string{},
		ds2db: ds2db,
		stats: stats,
		entities: entities,
		matcher: matcher,
		log: log,
	}
}

func (factBase *MySqlFactBase) GetMatchingGroups(set mentalese.RelationSet, keyCabinet *mentalese.KeyCabinet) []RelationGroup {
	return getFactBaseMatchingGroups(factBase.matcher, set, factBase, keyCabinet)
}

func (factBase *MySqlFactBase) GetMappings() []mentalese.RelationTransformation {
	return factBase.ds2db
}

func (factBase *MySqlFactBase) GetWriteMappings() []mentalese.RelationTransformation {
	return []mentalese.RelationTransformation{}
}

func (factBase *MySqlFactBase) GetStatistics() mentalese.DbStats {
	return factBase.stats
}

func (factBase *MySqlFactBase) GetEntities() mentalese.Entities {
	return factBase.entities
}

func (factBase *MySqlFactBase) AddTableDescription(tableName string, columns []string) {
	factBase.tableDescriptions[tableName] = columns
}

// Matches needleRelation to all relations in the database
// Returns a set of bindings
func (factBase *MySqlFactBase) MatchRelationToDatabase(needleRelation mentalese.Relation) []mentalese.Binding {

	factBase.log.StartDebug("MatchRelationToDatabase", needleRelation)

	dbBindings := []mentalese.Binding{}

	table := needleRelation.Predicate
	columns := factBase.tableDescriptions[table]

	whereClause := ""
	var values = []interface{}{}

	for i, argument := range needleRelation.Arguments {
		column := columns[i]
		if argument.TermType != mentalese.TermAnonymousVariable && argument.TermType != mentalese.TermVariable {
			whereClause += " AND " + column + " = ?"
			values = append(values, argument.TermValue)
		}
	}

	columnClause := strings.Join(columns, ", ")

	query := "SELECT " + columnClause + " FROM " + table + " WHERE TRUE" + whereClause

	rows, err := factBase.db.Query(query, values...)
	if err != nil {
		factBase.log.AddError("Error on querying MySQL: " + err.Error())
	}

	defer rows.Close()
	for rows.Next() {

		binding := mentalese.Binding{}

		// prepare an array of result value references
		resultValues := []string{}

		for range columns {
			resultValues = append(resultValues, "")
		}
		resultValueRefs := []interface{}{}
		for i := range columns {
			resultValueRefs = append(resultValueRefs, &resultValues[i])
		}

		// query all rows
		err := rows.Scan(resultValueRefs...)
		if err != nil {

			factBase.log.AddError("Error on querying MySQL: " + err.Error())

		} else {

			for i, argument := range needleRelation.Arguments {
				if argument.IsVariable() {
					variable := argument.TermValue
					binding[variable] = mentalese.Term{TermType: mentalese.TermStringConstant, TermValue: resultValues[i]}
				}
			}

			dbBindings = append(dbBindings, binding)
		}
	}

	factBase.log.EndDebug("MatchRelationToDatabase", dbBindings)

	return dbBindings
}

func (factBase *MySqlFactBase) Assert(relation mentalese.Relation) {

}

func (factBase *MySqlFactBase) Retract(relation mentalese.Relation) {

}
