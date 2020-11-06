package knowledge

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"nli-go/lib/common"
	"os"
	"strconv"
)

func (factBase *SparqlFactBase) doCachedQuery(query string) sparqlResponse {

	sparqlResponse := sparqlResponse{}
	queryHash := fmt.Sprintf("%x", md5.Sum([]byte(query)))
	cacheDir := factBase.cacheDir
	queryCachePath := cacheDir + "/" + queryHash + ".json"

	_, err := os.Stat(cacheDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(cacheDir, 0777)
		if err != nil {
			factBase.log.AddError("Error creating cache dir " + cacheDir + " (" + err.Error() + ")")
			return sparqlResponse
		}
	}

	// check if cache exists
	_, err = os.Stat(queryCachePath)
	if os.IsNotExist(err) {
		sparqlResponse = factBase.populateCache(query, queryCachePath)
	} else {
		sparqlResponse = factBase.readFromCache(queryCachePath)
		if factBase.log.Active() { factBase.log.AddDebug("SPARQL", query + " (from cache) " + strconv.Itoa(len(sparqlResponse.Results.Bindings)) + " results)") }
	}

	return sparqlResponse
}

func (factBase *SparqlFactBase) populateCache(query string, queryCachePath string) sparqlResponse {

	sparqlResponse := factBase.callSparql(query)

	// put them in the cache
	jsonBytes, _ := json.Marshal(sparqlResponse)
	err := common.WriteFile(queryCachePath, string(jsonBytes))
	if err != nil {
		factBase.log.AddError("Error creating cache file: " + err.Error())
		return sparqlResponse
	}

	return sparqlResponse
}

func (factBase *SparqlFactBase) readFromCache(queryCachePath string) sparqlResponse {

	sparqlResponse := sparqlResponse{}
	jsonString, err := common.ReadFile(queryCachePath)

	if err != nil {
		factBase.log.AddError("Error reading cache file: " + err.Error())
		return sparqlResponse
	}

	err = json.Unmarshal([]byte(jsonString), &sparqlResponse)
	if err != nil {
		factBase.log.AddError("Error parsing SPARQL response cache: " + err.Error())
		return sparqlResponse
	}

	return sparqlResponse
}
