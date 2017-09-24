package knowledge

// https://dinosaurscode.xyz/go/2016/06/19/golang-mysql-authentication/
// go get github.com/go-sql-driver/mysql

import "nli-go/lib/mentalese"
import "database/sql"
import (
	_ "github.com/go-sql-driver/mysql"
	"nli-go/lib/common"
	"strings"
)

type MySqlFactBase struct {
	db                *sql.DB
	tableDescriptions map[string][]string
	ds2db             []mentalese.RelationTransformation
	matcher           *mentalese.RelationMatcher
	log               *common.SystemLog
}

func NewMySqlFactBase(domain string, username string, password string, database string, matcher *mentalese.RelationMatcher, ds2db []mentalese.RelationTransformation, log *common.SystemLog) *MySqlFactBase {

	db, err := sql.Open("mysql", username+":"+password+"@/"+database)
	if err != nil {
		log.AddError("Error opening MySQL: " + err.Error())
	}

	return &MySqlFactBase{db: db, tableDescriptions: map[string][]string{}, ds2db: ds2db, matcher: matcher, log: log}
}

func (factBase MySqlFactBase) GetMatchingGroups(set mentalese.RelationSet, knowledgeBaseIndex int) RelationGroups {
	return getFactBaseMatchingGroups(factBase.matcher, set, factBase, knowledgeBaseIndex)
}

func (factBase MySqlFactBase) GetMappings() []mentalese.RelationTransformation {
	return factBase.ds2db
}

func (factBase MySqlFactBase) GetStatistics() mentalese.DbStats {
	// TODO
	return mentalese.DbStats{}
}

func (factBase MySqlFactBase) AddTableDescription(tableName string, columns []string) {
	factBase.tableDescriptions[tableName] = columns
}

func (factBase MySqlFactBase) Bind(goal []mentalese.Relation) ([]mentalese.Binding, bool) {

	factBase.log.StartDebug("MySqlFactBase.Bind", goal)

	internalBindings, match := factBase.MatchSequenceToDatabase(goal)

	factBase.log.EndDebug("MySqlFactBase.Bind", internalBindings, match)

	return internalBindings, match
}

// Matches a sequence of relations to the relations of the MySql database
// sequence: [ marriages(A, C) person(A, 'John', _, _) ]
// return: [ { C: 1, A: 5 } ]
func (factBase MySqlFactBase) MatchSequenceToDatabase(sequence mentalese.RelationSet) ([]mentalese.Binding, bool) {

	factBase.log.StartDebug("MatchSequenceToDatabase", sequence)

	// bindings using database level variables
	sequenceBindings := []mentalese.Binding{}
	match := true

	for _, relation := range sequence {

		relationBindings := []mentalese.Binding{}

		if len(relationBindings) == 0 {

			resultBindings := factBase.matchRelationToDatabase(relation)
			relationBindings = resultBindings

		} else {

			// go through the bindings resulting from previous relation
			for _, binding := range sequenceBindings {

				boundRelation := factBase.matcher.BindSingleRelationSingleBinding(relation, binding)
				resultBindings := factBase.matchRelationToDatabase(boundRelation)

				// found bindings must be extended with the bindings already present
				for _, resultBinding := range resultBindings {
					newRelationBinding := binding.Merge(resultBinding)
					relationBindings = append(relationBindings, newRelationBinding)
				}
			}
		}

		sequenceBindings = relationBindings

		if len(sequenceBindings) == 0 {
			match = false
			break
		}
	}

	factBase.log.EndDebug("MatchSequenceToDatabase", sequenceBindings, match)

	return sequenceBindings, match
}

// Matches needleRelation to all relations in the database
// Returns a set of bindings
func (factBase MySqlFactBase) matchRelationToDatabase(needleRelation mentalese.Relation) []mentalese.Binding {

	factBase.log.StartDebug("matchRelationToDatabase", needleRelation)

	dbBindings := []mentalese.Binding{}

	table := needleRelation.Predicate
	columns := factBase.tableDescriptions[table]

	whereClause := ""
	var values = []interface{}{}

	for i, argument := range needleRelation.Arguments {
		column := columns[i]
		if argument.TermType != mentalese.Term_anonymousVariable && argument.TermType != mentalese.Term_variable {
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
					binding[variable] = mentalese.Term{TermType: mentalese.Term_stringConstant, TermValue: resultValues[i]}
				}
			}

			dbBindings = append(dbBindings, binding)
		}
	}

	factBase.log.EndDebug("matchRelationToDatabase", dbBindings)

	return dbBindings
}
