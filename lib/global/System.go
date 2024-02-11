package global

import (
	"nli-go/lib/api"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"

	"golang.org/x/net/websocket"
)

type System struct {
	log                   *common.SystemLog
	dialogContext         *central.DialogContext
	internalGrammarParser *importer.InternalGrammarParser
	nameResolver          *central.NameResolver
	grammars              []parse.Grammar
	meta                  *mentalese.Meta
	relationizer          *parse.Relationizer
	matcher               *central.RelationMatcher
	variableGenerator     *mentalese.VariableGenerator
	solver                *central.ProblemSolver
	answerer              *central.Answerer
	generator             *generate.Generator
	surfacer              *generate.SurfaceRepresentation
	processList           *central.ProcessList
	processRunner         *central.ProcessRunner
	clientConnector       api.ClientConnector
}

func (system *System) GetLog() *common.SystemLog {
	return system.log
}

func (system *System) GetClientConnector() api.ClientConnector {
	return system.clientConnector
}

func (system *System) CreatClientConnector(conn *websocket.Conn) *ClientConnector {
	return &ClientConnector{
		conn:   conn,
		system: system,
	}
}

func (system *System) HandleRequest(request mentalese.Request) {
	switch request.MessageType {
	case mentalese.MessageRespond:
		// start a process
		system.processRunner.StartProcess(
			central.RESOURCE_LANGUAGE,
			mentalese.RelationSet{
				mentalese.NewRelation(false, mentalese.PredicateRespond,
					[]mentalese.Term{mentalese.NewTermString(request.Message.(string))}, 0),
			},
			mentalese.NewBinding(),
		)
	default:
		// send a message to a process
		system.processRunner.SendMessage(request)
	}
}

func (system *System) RunRelationSet(resource string, relationSet mentalese.RelationSet) mentalese.BindingSet {
	return system.processRunner.RunRelationSet(resource, relationSet)
}

func (system *System) RunRelationSetString(resource string, relationSet string) mentalese.BindingSet {
	relations := system.internalGrammarParser.CreateRelationSet(relationSet)
	return system.processRunner.RunRelationSet(resource, relations)
}
