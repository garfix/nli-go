package knowledge

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"strings"
)

const max_sparql_results = 100

type SparqlFactBase struct {
	baseUrl           string
	defaultGraphUri   string
	ds2db             []mentalese.RelationTransformation
	names 			  mentalese.ConfigMap
	stats			  mentalese.DbStats
	matcher           *mentalese.RelationMatcher
	log               *common.SystemLog
}

func NewSparqlFactBase(baseUrl string, defaultGraphUri string, matcher *mentalese.RelationMatcher, ds2db []mentalese.RelationTransformation, names mentalese.ConfigMap, stats mentalese.DbStats, log *common.SystemLog) *SparqlFactBase {

	return &SparqlFactBase{baseUrl: baseUrl, defaultGraphUri: defaultGraphUri, ds2db: ds2db, names: names, stats: stats, matcher: matcher, log: log}
}

func (factBase SparqlFactBase) Bind(goal []mentalese.Relation) ([]mentalese.Binding, bool) {

	factBase.log.StartDebug("SparqlFactBase Bind", goal)

	// bindings using database level variables
	sequenceBindings := []mentalese.Binding{}
	match := true

	for _, relation := range goal {

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

	factBase.log.EndDebug("SparqlFactBase Bind", sequenceBindings)

	return sequenceBindings, match
}

func (factBase SparqlFactBase) GetMappings() []mentalese.RelationTransformation {
	return factBase.ds2db
}

func (factBase SparqlFactBase) GetMatchingGroups(set mentalese.RelationSet, knowledgeBaseIndex int) []RelationGroup {
	return getFactBaseMatchingGroups(factBase.matcher, set, factBase, knowledgeBaseIndex)
}

func (factBase SparqlFactBase) GetStatistics() mentalese.DbStats {
	return factBase.stats
}

// Matches needleRelation to all relations in the database
// Returns a set of bindings
func (factBase SparqlFactBase) matchRelationToDatabase(relation mentalese.Relation) []mentalese.Binding {

	factBase.log.StartDebug("matchRelationToDatabase", relation)

	bindings := []mentalese.Binding{}

	if len(relation.Arguments) != 2 {
		factBase.log.AddError("Relation does not have exactly two arguments: " + relation.String())
		return bindings
	}

	var1 := ""
	var2 := ""
	variables := []string{}

	if relation.Arguments[0].TermType == mentalese.Term_anonymousVariable || relation.Arguments[0].TermType == mentalese.Term_variable {
		var1 = "?variable1"
		variables = append(variables, var1)
	} else {
		var1 = relation.Arguments[0].String()
		if relation.Arguments[0].TermType == mentalese.Term_stringConstant {
			var1 += "@en"
		} else if relation.Arguments[0].TermType == mentalese.Term_id {
			var1 = "<" + var1 + ">"
		}
	}

	if relation.Arguments[1].TermType == mentalese.Term_anonymousVariable || relation.Arguments[1].TermType == mentalese.Term_variable {
		var2 = "?variable2"
		variables = append(variables, var2)
	} else {
		var2 = relation.Arguments[1].String()
		if relation.Arguments[1].TermType == mentalese.Term_stringConstant {
			var2 += "@en"
		} else if relation.Arguments[1].TermType == mentalese.Term_id {
			var2 = "<" + var2 + ">"
		}
	}

	if len(variables) == 0 {
		variables = append(variables, "1")
	}

	relationUri, ok := factBase.names[relation.Predicate]
	if !ok {
		factBase.log.AddError("Relation uri not found in names: " + relation.Predicate)
		return bindings
	}

	query := "select " + strings.Join(variables, ", ") + " where { " + var1 + " <" + relationUri + "> " + var2  + "} limit " + strconv.Itoa(max_sparql_results)

	resp, err := http.PostForm(factBase.baseUrl,
		url.Values{
			"default-graph-uri": {factBase.defaultGraphUri},
			"query": {query},
			"format": {"application/json"},
		})

	if err != nil {
		factBase.log.AddError("Error posting SPARQL request: " + err.Error())
		return bindings
	}

	defer resp.Body.Close()
	bodyJson, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		factBase.log.AddError("Error reading body of SPARQL response: " + err.Error())
		return bindings
	}

	response := sparqlResponse{}
	err = json.Unmarshal([]byte(bodyJson), &response)
	if err != nil {
		factBase.log.AddError("Error parsing SPARQL response JSON: " + err.Error() + "\nResponse body: " + string(bodyJson))
		return bindings
	}

	for _, resultBinding := range response.Results.Bindings  {

		binding := mentalese.Binding{}

		if relation.Arguments[0].IsVariable() {

			termType := mentalese.Term_stringConstant
			if resultBinding.Variable1.Type == "uri" {
				termType = mentalese.Term_id
			}
			binding[relation.Arguments[0].TermValue] = mentalese.Term{ TermType: termType, TermValue: resultBinding.Variable1.Value }
		}

		if relation.Arguments[1].IsVariable() {

			termType := mentalese.Term_stringConstant
			if resultBinding.Variable2.Type == "uri" {
				termType = mentalese.Term_id
			}

			binding[relation.Arguments[1].TermValue] = mentalese.Term{ TermType: termType, TermValue: resultBinding.Variable2.Value }
		}

		bindings = append(bindings, binding)
	}

	factBase.log.EndDebug("matchRelationToDatabase", bindings)

	return bindings
}
