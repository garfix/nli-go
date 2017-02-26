package knowledge

// https://dinosaurscode.xyz/go/2016/06/19/golang-mysql-authentication/
// go get github.com/go-sql-driver/mysql

import "nli-go/lib/mentalese"
import "database/sql"
import (
	_ "github.com/go-sql-driver/mysql"
	"nli-go/lib/common"
	"fmt"
	"strings"
	"log"
)

type MySqlFactBase struct{
	db *sql.DB
	tableDescriptions map[string] []string
	ds2db []mentalese.Rule
	matcher *mentalese.RelationMatcher
}

func NewMySqlFactBase(domain string, username string, password string, database string, ds2db []mentalese.Rule) *MySqlFactBase {

	db, err := sql.Open("mysql", username + ":" + password + "@/" + database)
	if err != nil {
		panic(err.Error())
	}

	return &MySqlFactBase{ db: db, tableDescriptions: map[string] []string{}, ds2db: ds2db, matcher: mentalese.NewRelationMatcher() }
}

func (factBase MySqlFactBase) AddTableDescription(tableName string, columns []string) {
	factBase.tableDescriptions[tableName] = columns
}

// todo: remove code duplication
func (factBase MySqlFactBase) Bind(goal mentalese.Relation) []mentalese.Binding {

common.LoggerActive=true

	common.LogTree("MySqlFactBase.Bind", goal)

	bindings := []mentalese.Binding{}

	for _, ds2db := range factBase.ds2db {

		// gender(14, G), gender(A, male) => externalBinding: G = male
		externalBinding, match := factBase.matcher.MatchTwoRelations(goal, ds2db.Goal, mentalese.Binding{})
		if match {

			// gender(14, G), gender(A, male) => internalBinding: A = 14
			internalBinding, _ := factBase.matcher.MatchTwoRelations(ds2db.Goal, goal, mentalese.Binding{})

			// create a version of the conditions with bound variables
			boundConditions := factBase.matcher.BindRelationSetSingleBinding(ds2db.Pattern, internalBinding)
			// match this bound version to the database
			internalBindings, match := factBase.MatchSequenceToDatabase(boundConditions)

			if match {
				for _, binding := range internalBindings {
					bindings = append(bindings, externalBinding.Intersection(binding))
				}
			}
		}
	}

	common.LogTree("MySqlFactBase.Bind", bindings)

	common.LoggerActive=false

	return bindings
}

// Matches a sequence of relations to the relations of the MySql database
// sequence: [ marriages(A, C) person(A, 'John', _, _) ]
func (factBase MySqlFactBase) MatchSequenceToDatabase(sequence mentalese.RelationSet) ([]mentalese.Binding, bool){

	common.LogTree("MatchSequenceToDatabase", sequence)

	// bindings using database level variables
	sequenceBindings := []mentalese.Binding{}
	match := true

	for _, relation := range sequence {

		relationBindings := []mentalese.Binding{}

		if (len(relationBindings) == 0) {

			resultBindings := factBase.matchRelationToDatabase(relation)
			relationBindings = append(relationBindings, resultBindings...)

		} else {

			for _, binding := range sequenceBindings {
				boundRelation := factBase.matcher.BindSingleRelationSingleBinding(relation, binding)
				resultBindings := factBase.matchRelationToDatabase(boundRelation)
				relationBindings = append(relationBindings, resultBindings...)
			}
		}

		if len(relationBindings) == 0 {
			match = false
			break
		}

		sequenceBindings = relationBindings
	}

	common.LogTree("MatchSequenceToDatabase", sequenceBindings, match)

	return sequenceBindings, match
}

func (factBase MySqlFactBase) matchRelationToDatabase(needleRelation mentalese.Relation) []mentalese.Binding {

	common.LogTree("matchRelationToDatabase", needleRelation)

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
	fmt.Printf("        %s, %v\n", query, values)

	rows, err := factBase.db.Query(query, values...)
	if err != nil {
		panic(err.Error())
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
		if err := rows.Scan(resultValueRefs...); err != nil {
			log.Fatal(err)
		}

		for i, argument := range needleRelation.Arguments {
			if argument.IsVariable() {
				variable := argument.TermValue
				binding[variable] = mentalese.Term{ TermType: mentalese.Term_stringConstant, TermValue: resultValues[i] }
			}
		}

		dbBindings = append(dbBindings, binding)
	}

	common.LogTree("matchRelationToDatabase", dbBindings)

	return dbBindings
}