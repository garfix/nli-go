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

const maxSparqlResults = 5000
const MaxQueries = 200

type SparqlFactBase struct {
	KnowledgeBaseCore
	baseUrl         string
	defaultGraphUri string
	ds2db           []mentalese.RelationTransformation
	names           mentalese.ConfigMap
	stats           mentalese.DbStats
	entities        mentalese.Entities
	predicates 		mentalese.Predicates
	sharedIds 		SharedIds
	matcher         *mentalese.RelationMatcher
	queryCount      int
	doCache         bool
	log             *common.SystemLog
}

func NewSparqlFactBase(name string, baseUrl string, defaultGraphUri string, matcher *mentalese.RelationMatcher, ds2db []mentalese.RelationTransformation, names mentalese.ConfigMap, stats mentalese.DbStats, entities mentalese.Entities, predicates mentalese.Predicates, doCache bool, log *common.SystemLog) *SparqlFactBase {

	return &SparqlFactBase{
		KnowledgeBaseCore: KnowledgeBaseCore{ Name: name},
		baseUrl:           baseUrl,
		defaultGraphUri:   defaultGraphUri,
		ds2db:             ds2db,
		names:             names,
		stats:             stats,
		entities:          entities,
		predicates:		   predicates,
		sharedIds: 		   SharedIds{},
		matcher:           matcher,
		queryCount:        0,
		doCache:           doCache,
		log:               log,
	}
}

func (factBase *SparqlFactBase) GetMappings() []mentalese.RelationTransformation {
	return factBase.ds2db
}

func (factBase *SparqlFactBase) GetWriteMappings() []mentalese.RelationTransformation {
	return []mentalese.RelationTransformation{}
}

func (factBase *SparqlFactBase) GetMatchingGroups(set mentalese.RelationSet) []RelationGroup {
	return getFactBaseMatchingGroups(factBase.matcher, set, factBase)
}

func (factBase *SparqlFactBase) GetStatistics() mentalese.DbStats {
	return factBase.stats
}

func (factBase *SparqlFactBase) GetEntities() mentalese.Entities {
	return factBase.entities
}

func (factBase *SparqlFactBase) SetSharedIds(sharedIds SharedIds) {
	factBase.sharedIds = sharedIds
}

func (factBase *SparqlFactBase) GetLocalId(inId string, entityType string) string {
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

func (factBase *SparqlFactBase) GetSharedId(inId string, entityType string) string {
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

// Matches needleRelation to all relations in the database
// Returns a set of bindings
func (factBase *SparqlFactBase) MatchRelationToDatabase(relation mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	factBase.log.StartDebug("MatchRelationToDatabase", relation)

	relation = relation.BindSingleRelationSingleBinding(binding)

	bindings := mentalese.Bindings{}

	if len(relation.Arguments) != 2 {
		factBase.log.AddError("Relation does not have exactly two arguments: " + relation.String())
		return bindings
	}

	if factBase.queryCount > MaxQueries {
		factBase.log.AddError("Too many SPARQL queries")
		return bindings
	}

	bindings = factBase.doQuery(relation)

	factBase.log.EndDebug("MatchRelationToDatabase", bindings)

	return bindings
}

func (factBase *SparqlFactBase) doQuery(relation mentalese.Relation) mentalese.Bindings {

	bindings := mentalese.Bindings{}
	sparqlResponse := sparqlResponse{}

	query := factBase.createQuery(relation)
	if query == "" {
		return bindings
	}

	if factBase.doCache {
		sparqlResponse = factBase.doCachedQuery(query)
	} else {
		sparqlResponse = factBase.callSparql(query)
	}

	bindings = factBase.processSparqlResponse(relation, sparqlResponse)

	return bindings
}

func (factBase *SparqlFactBase) createQuery(relation mentalese.Relation) string {
	var1 := ""
	var2 := ""
	extra := ""
	variables := []string{}

	if relation.Arguments[0].TermType == mentalese.TermAnonymousVariable || relation.Arguments[0].TermType == mentalese.TermVariable {
		var1 = "?variable1"
		variables = append(variables, var1)
	} else {
		var1 = relation.Arguments[0].String()
		if relation.Arguments[0].TermType == mentalese.TermStringConstant {
			var1 += "@en"
		} else if relation.Arguments[0].IsId() {
			var1 = "<" + relation.Arguments[0].TermValue + ">"
		}
	}

	if relation.Arguments[1].TermType == mentalese.TermAnonymousVariable || relation.Arguments[1].TermType == mentalese.TermVariable {
		var2 = "?variable2"
		variables = append(variables, var2)
	} else {
		if relation.Arguments[1].TermType == mentalese.TermStringConstant {
//todo make this into a config value
			arg1 := relation.Arguments[1].TermValue
			if false {
				var2 = relation.Arguments[1].String() + "@en"
			// searching for punctuation marks leads to errors
			} else if arg1 == "'?'" {
				return ""
			} else {
				// case insensitive search
				var2 = "?name"
				// http://docs.openlinksw.com/virtuoso/rdfpredicatessparql/
				escA := strings.Replace(strings.Replace(arg1, "'", "\\\\'", -1), "\"", "\\\"", -1)
				extra = " . ?name bif:contains \"'" + escA + "'\""
				// since bif:contains is inexact, it yields false positives; correct for those
				escB := strings.Replace(arg1, "'", "\\'", -1)
				extra += " . FILTER (LCASE(STR(?name)) = '" + strings.ToLower(escB) + "')"
			}
		} else if relation.Arguments[1].IsId() {
			var2 = "<" + relation.Arguments[1].TermValue + ">"
		}
	}

	if len(variables) == 0 {
		variables = append(variables, "1")
	}

	relationUri, ok := factBase.names[relation.Predicate]
	if !ok {
		factBase.log.AddError("Relation uri not found in names: " + relation.Predicate)
		return ""
	}

	query := "select " + strings.Join(variables, ", ") + " where { " + var1 + " <" + relationUri + "> " + var2  + extra + "} limit " + strconv.Itoa(maxSparqlResults)

	return query
}

func (factBase *SparqlFactBase) callSparql(query string) sparqlResponse {

	sparqlResponse := sparqlResponse{}

	start := time.Now()

	factBase.queryCount++

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
		return sparqlResponse
	} else if resp.StatusCode == 405 {
		factBase.log.AddError("Error posting SPARQL request: 405 Not Allowed. Probably too many queries (" + strconv.Itoa(factBase.queryCount) + ")")
		return sparqlResponse
	} else if resp.StatusCode != 200 {
		factBase.log.AddError("Error posting SPARQL request: " + http.StatusText(resp.StatusCode))
		return sparqlResponse
	}

	defer resp.Body.Close()

	bodyJson, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		factBase.log.AddError("Error reading body of SPARQL sparqlResponse: " + err.Error())
		return sparqlResponse
	}

	err = json.Unmarshal([]byte(bodyJson), &sparqlResponse)
	if err != nil {
		factBase.log.AddError("Error parsing SPARQL sparqlResponse JSON: " + err.Error() + "\nResponse body: " + string(bodyJson))
		return sparqlResponse
	}

	factBase.log.AddProduction("SPARQL", query + " (" + elapsed.String() + ", " + strconv.Itoa(len(sparqlResponse.Results.Bindings)) + " results)")

	return sparqlResponse
}

func (factBase *SparqlFactBase) processSparqlResponse(relation mentalese.Relation,  sparqlResponse sparqlResponse) mentalese.Bindings {

	bindings := mentalese.Bindings{}

	for _, resultBinding := range sparqlResponse.Results.Bindings {

		binding := mentalese.Binding{}

		if relation.Arguments[0].IsVariable() {

			if resultBinding.Variable1.Type == "uri" {
// todo predicates does not contain database relations
				entityType := factBase.predicates.GetEntityType(relation.Predicate, 0)
				binding[relation.Arguments[0].TermValue] = mentalese.NewId(resultBinding.Variable1.Value, entityType)
			} else {
				// skip non-english results
				if resultBinding.Variable1.Lang != "" && resultBinding.Variable1.Lang != "en" {
					continue
				}
				binding[relation.Arguments[0].TermValue] = mentalese.NewString(resultBinding.Variable1.Value)
			}
		}

		if relation.Arguments[1].IsVariable() {

			if resultBinding.Variable2.Type == "uri" {
// todo predicates does not contain database relations
				entityType := factBase.predicates.GetEntityType(relation.Predicate, 1)
				binding[relation.Arguments[1].TermValue] = mentalese.NewId(resultBinding.Variable2.Value, entityType)
			} else {
				// skip non-english results
				if resultBinding.Variable2.Lang != "" && resultBinding.Variable2.Lang != "en" {
					continue
				}
				binding[relation.Arguments[1].TermValue] = mentalese.NewString(resultBinding.Variable2.Value)
			}
		}

		bindings = append(bindings, binding)
	}

	return bindings
}

func (factBase *SparqlFactBase) Assert(relation mentalese.Relation) {

}

func (factBase *SparqlFactBase) Retract(relation mentalese.Relation) {

}
