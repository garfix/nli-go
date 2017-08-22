package knowledge

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
)

type SparqlFactBase struct {
	baseUrl           string
	defaultGraphUri   string
	ds2db             []mentalese.DbMapping
	matcher           *mentalese.RelationMatcher
	log               *common.SystemLog
}

func NewSparqlFactBase(baseUrl string, defaultGraphUri string, ds2db []mentalese.DbMapping, log *common.SystemLog) *SparqlFactBase {

	return &SparqlFactBase{baseUrl: baseUrl, defaultGraphUri: defaultGraphUri, ds2db: ds2db, matcher: mentalese.NewRelationMatcher(log), log: log}
}

func (factBase SparqlFactBase) Bind(goal []mentalese.Relation) ([]mentalese.Binding, bool) {

	bindings := []mentalese.Binding{}

	return bindings, true
}

func (factBase SparqlFactBase) GetMappings() []mentalese.DbMapping {
	return factBase.ds2db
}