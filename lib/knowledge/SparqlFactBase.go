package knowledge

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
	"strings"
	"time"
)

const max_sparql_results = 100

const MAX_QUERIES = 100


type SparqlFactBase struct {
	KnowledgeBaseCore
	baseUrl           string
	defaultGraphUri   string
	ds2db             []mentalese.RelationTransformation
	names 			  mentalese.ConfigMap
	stats			  mentalese.DbStats
	entities 		  mentalese.Entities
	matcher           *mentalese.RelationMatcher
	queryCount 		  int
	log               *common.SystemLog
}

func NewSparqlFactBase(name string, baseUrl string, defaultGraphUri string, matcher *mentalese.RelationMatcher, ds2db []mentalese.RelationTransformation, names mentalese.ConfigMap, stats mentalese.DbStats, entities mentalese.Entities, log *common.SystemLog) *SparqlFactBase {

	return &SparqlFactBase{
		KnowledgeBaseCore: KnowledgeBaseCore{ Name: name},
		baseUrl: baseUrl,
		defaultGraphUri: defaultGraphUri,
		ds2db: ds2db,
		names: names,
		stats: stats,
		entities: entities,
		matcher: matcher,
		queryCount: 0,
		log: log,
	}
}

func (factBase *SparqlFactBase) GetMappings() []mentalese.RelationTransformation {
	return factBase.ds2db
}

func (factBase *SparqlFactBase) GetMatchingGroups(set mentalese.RelationSet, knowledgeBaseIndex int) []RelationGroup {
	return getFactBaseMatchingGroups(factBase.matcher, set, factBase, knowledgeBaseIndex)
}

func (factBase *SparqlFactBase) GetStatistics() mentalese.DbStats {
	return factBase.stats
}

func (factBase *SparqlFactBase) GetEntities() mentalese.Entities {
	return factBase.entities
}

// Matches needleRelation to all relations in the database
// Returns a set of bindings
func (factBase *SparqlFactBase) MatchRelationToDatabase(relation mentalese.Relation) []mentalese.Binding {

	factBase.log.StartDebug("MatchRelationToDatabase", relation)

	bindings := []mentalese.Binding{}

	factBase.queryCount++

	if factBase.queryCount > MAX_QUERIES {
		factBase.log.AddError("Too many SPARQL queries")
		return bindings
	}


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

	start := time.Now()

	resp, err := http.PostForm(factBase.baseUrl,
		url.Values{
			"default-graph-uri": {factBase.defaultGraphUri},
			"query": {query},
			"format": {"application/json"},
		})

	t := time.Now()
	elapsed := t.Sub(start)

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

	factBase.log.AddProduction("SPARQL Query", query + " (" + elapsed.String() + ", " + strconv.Itoa(len(response.Results.Bindings)) + " results)")

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

	factBase.log.EndDebug("MatchRelationToDatabase", bindings)

	return bindings
}
