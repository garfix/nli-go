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
	tableDescriptions map[string]tableDescription
	ds2db             []mentalese.Rule
	entities 		  mentalese.Entities
	sharedIds 		  SharedIds
	matcher           *mentalese.RelationMatcher
	log               *common.SystemLog
}

type tableDescription struct {
	tableName string
	columns []string
}

func NewMySqlFactBase(name string, domain string, username string, password string, database string, matcher *mentalese.RelationMatcher, ds2db []mentalese.Rule, entities mentalese.Entities, log *common.SystemLog) *MySqlFactBase {

	db, err := sql.Open("mysql", username+":"+password+"@/"+database)
	if err != nil {
		log.AddError("Error opening MySQL: " + err.Error())
	}

	return &MySqlFactBase{
		KnowledgeBaseCore: KnowledgeBaseCore{ Name: name},
		db: db,
		tableDescriptions: map[string]tableDescription{},
		ds2db: ds2db,
		entities: entities,
		sharedIds: SharedIds{},
		matcher: matcher,
		log: log,
	}
}

func (factBase *MySqlFactBase) HandlesPredicate(predicate string) bool {
	for _, rule := range factBase.ds2db {
		if rule.Goal.Predicate == predicate {
			return true
		}
	}
	return false
}

func (factBase *MySqlFactBase) GetMappings() []mentalese.Rule {
	return factBase.ds2db
}

func (factBase *MySqlFactBase) GetWriteMappings() []mentalese.Rule {
	return []mentalese.Rule{}
}

func (factBase *MySqlFactBase) GetEntities() mentalese.Entities {
	return factBase.entities
}

func (factBase *MySqlFactBase) SetSharedIds(sharedIds SharedIds) {
	factBase.sharedIds = sharedIds
}

func (factBase *MySqlFactBase) GetLocalId(inId string, entityType string) string {
	outId := ""

	_, found := factBase.sharedIds[entityType]
	if !found { return inId }

	for localId, sharedId := range factBase.sharedIds[entityType] {
		if inId == sharedId {
			outId = localId
			break
		}
	}

	return outId
}

func (factBase *MySqlFactBase) GetSharedId(inId string, entityType string) string {
	outId := ""

	_, found := factBase.sharedIds[entityType]
	if !found { return inId }

	for localId, sharedId := range factBase.sharedIds[entityType] {
		if inId == localId {
			outId = sharedId
			break
		}
	}

	return outId
}

func (factBase *MySqlFactBase) AddTableDescription(predicate string, tableName string, columns []string) {
	factBase.tableDescriptions[predicate] = tableDescription{ tableName: tableName, columns: columns }
}

// Matches needleRelation to all relations in the database
// Returns a set of bindings
func (factBase *MySqlFactBase) MatchRelationToDatabase(needleRelation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	factBase.log.StartDebug("MatchRelationToDatabase", needleRelation)

	needleRelation = needleRelation.BindSingle(binding)

	dbBindings := mentalese.Bindings{}

	description := factBase.tableDescriptions[needleRelation.Predicate]
	columns := description.columns
	tableName := description.tableName

	whereClause := ""
	var values = []interface{}{}

	for i, argument := range needleRelation.Arguments {
		column := columns[i]
		if !argument.IsAnonymousVariable() && !argument.IsVariable() {
			whereClause += " AND " + column + " = ?"
			values = append(values, argument.TermValue)
		}
	}

	columnClause := strings.Join(columns, ", ")

	query := "SELECT " + columnClause + " FROM " + tableName + " WHERE TRUE" + whereClause

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
					binding[variable] = mentalese.Term{TermType: mentalese.TermTypeStringConstant, TermValue: resultValues[i]}
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
