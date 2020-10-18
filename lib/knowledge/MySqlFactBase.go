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
	readMap           []mentalese.Rule
	writeMap           []mentalese.Rule
	sharedIds         SharedIds
	matcher           *mentalese.RelationMatcher
	log               *common.SystemLog
}

type tableDescription struct {
	tableName string
	columns []string
}

func NewMySqlFactBase(name string, username string, password string, database string, matcher *mentalese.RelationMatcher, readMap []mentalese.Rule, writeMap []mentalese.Rule, log *common.SystemLog) *MySqlFactBase {

	db, _ := sql.Open("mysql", username+":"+password+"@/"+database)
	err := db.Ping()
	if err != nil {
		log.AddError("Error opening MySQL: " + err.Error())
	}

	return &MySqlFactBase{
		KnowledgeBaseCore: KnowledgeBaseCore{ Name: name},
		db:                db,
		tableDescriptions: map[string]tableDescription{},
		readMap:           readMap,
		writeMap: 		   writeMap,
		sharedIds:         SharedIds{},
		matcher:           matcher,
		log:               log,
	}
}

func (factBase *MySqlFactBase) HandlesPredicate(predicate string) bool {
	for _, rule := range factBase.readMap {
		if rule.Goal.Predicate == predicate {
			return true
		}
	}
	if len(factBase.writeMap) > 0 && (predicate == mentalese.PredicateAssert || predicate == mentalese.PredicateRetract) {
		return true
	}
	return false
}

func (factBase *MySqlFactBase) GetReadMappings() []mentalese.Rule {
	return factBase.readMap
}

func (factBase *MySqlFactBase) GetWriteMappings() []mentalese.Rule {
	return factBase.writeMap
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
func (factBase *MySqlFactBase) MatchRelationToDatabase(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	factBase.log.StartDebug("MatchRelationToDatabase", relation)

	relation = relation.BindSingle(binding)

	dbBindings := mentalese.Bindings{}
	description := factBase.tableDescriptions[relation.Predicate]
	columns := description.columns
	tableName := description.tableName

	whereClause := ""
	var values = []interface{}{}

	for i, argument := range relation.Arguments {
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

		binding := mentalese.NewBinding()

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

			for i, argument := range relation.Arguments {
				if argument.IsVariable() {
					variable := argument.TermValue
					binding.Set(variable, mentalese.Term{TermType: mentalese.TermTypeStringConstant, TermValue: resultValues[i]})
				}
			}

			dbBindings = append(dbBindings, binding)
		}
	}

	factBase.log.EndDebug("MatchRelationToDatabase", dbBindings)

	return dbBindings
}

func (factBase *MySqlFactBase) Assert(relation mentalese.Relation) {

	factBase.log.StartDebug("Assert", relation)

	argCount := len(relation.Arguments)

	if argCount == 0 { return }

	// check if relation already present; do not duplicate!
	existingBindings := factBase.MatchRelationToDatabase(relation, mentalese.NewBinding())
	if len(existingBindings) > 0 { return }

	description := factBase.tableDescriptions[relation.Predicate]
	columns := description.columns
	tableName := description.tableName

	var values = []interface{}{}
	valueClause := ""
	sep := ""

	for i, argument := range relation.Arguments {
		column := columns[i]
		if !argument.IsAnonymousVariable() && !argument.IsVariable() {
			valueClause += sep + column + " = ?"
			sep = ", "
			values = append(values, argument.TermValue)
		}
	}

	query := "INSERT INTO " + tableName + " SET " + valueClause
	_, err := factBase.db.Exec(query, values...)

	if err != nil {
		factBase.log.AddError(err.Error())
	}

	factBase.log.EndDebug("Assert", relation)
}

func (factBase *MySqlFactBase) Retract(relation mentalese.Relation) {

	factBase.log.StartDebug("Retract", relation)

	argCount := len(relation.Arguments)

	if argCount == 0 { return }

	description := factBase.tableDescriptions[relation.Predicate]
	columns := description.columns
	tableName := description.tableName

	whereClause := ""
	sep := ""
	var values = []interface{}{}

	for i, argument := range relation.Arguments {
		column := columns[i]
		if !argument.IsAnonymousVariable() && !argument.IsVariable() {
			whereClause += sep + column + " = ?"
			sep = " AND "
			values = append(values, argument.TermValue)
		}
	}

	query := "DELETE FROM " + tableName + " WHERE " + whereClause

	_, err := factBase.db.Exec(query, values...)
	if err != nil {
		factBase.log.AddError("Error on querying MySQL: " + err.Error())
	}

	factBase.log.EndDebug("Retract", relation)
}
