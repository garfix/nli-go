package global

import (
	"fmt"
	"nli-go/lib/api"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"strconv"

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

	waitingFor *mentalese.Relation
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

// Low-level function to inspect the internal state of the system
func (system *System) Query(relations string) mentalese.BindingSet {
	set := system.internalGrammarParser.CreateRelationSet(relations)
	result := system.processRunner.RunRelationSet(central.SIMPLE_PROCESS, set)

	return result
}

func (system *System) AddListener(l central.MessageListener) []mentalese.RelationSet {
	return system.processList.AddListener(l)
}

func (system *System) RemoveListener(l central.MessageListener) {
	system.processList.RemoveListener(l)
}

func (system *System) SendAndWaitForResponse(clientMessage mentalese.RelationSet) mentalese.RelationSet {

	responseMessage := mentalese.RelationSet{}

	// waiting on another thread
	// https://medium.com/@matryer/golang-advent-calendar-day-two-starting-and-stopping-things-with-a-signal-channel-f5048161018
	done := make(chan struct{})

	callback := func(serverMessage mentalese.RelationSet) {

		if len(serverMessage) > 0 {
			responseMessage = responseMessage.Merge(serverMessage)
			close(done)
		} else {
			if system.processList.IsEmpty() {
				close(done)
			}
		}
	}

	pendingMessages := system.AddListener(callback)
	responseMessage = responseMessage.MergeMultiple(pendingMessages)

	defer system.RemoveListener(callback)

	// enter the message into the system
	system.processRunner.RunRelationSet(central.LANGUAGE_PROCESS, clientMessage)

	// wait until done
	<-done

	fmt.Println("response: " + responseMessage.String() + "\n")
	if responseMessage.IsEmpty() {
		fmt.Println("------------------------------\n")
	}

	return responseMessage
}

func (system *System) Answer(input string) (string, *common.Options) {

	answer := ""
	anAnswer := ""
	options := common.NewOptions()
	someOptions := common.NewOptions()
	responseMessage := mentalese.RelationSet{}

	done := make(chan struct{})

	callback := func(message mentalese.RelationSet) {
		if len(message) > 0 {
			firstRelation := message[0]
			if firstRelation.Predicate == mentalese.PredicatePrint {
				responseMessage = responseMessage.Merge(message)
			}
			if firstRelation.Predicate == mentalese.PredicateUserSelect {
				responseMessage = responseMessage.Merge(message)
				close(done)
				return
			}
			// ping pong
			for _, relation := range message {
				go system.assert(relation)
			}
		} else {
			if system.processList.IsEmpty() {
				close(done)
			}
		}
	}

	pendingMessages := system.AddListener(callback)
	responseMessage = responseMessage.MergeMultiple(pendingMessages)

	defer system.RemoveListener(callback)

	// find or create a goal
	system.createOrUpdateProcess(input)

	// wait until done
	<-done

	// read answer action
	anAnswer, someOptions, _ = system.readAnswer(responseMessage)
	if anAnswer != "" {
		answer = anAnswer
	}
	if someOptions.HasOptions() {
		options = someOptions
	}

	return answer, options
}

func (system *System) assert(relation mentalese.Relation) {

	// go:assert()
	set := mentalese.RelationSet{
		mentalese.NewRelation(false, mentalese.PredicateAssert, []mentalese.Term{
			mentalese.NewTermRelationSet(mentalese.RelationSet{relation}),
		}),
	}
	system.processRunner.RunRelationSet(central.SIMPLE_PROCESS, set)
}

func (system *System) createOrUpdateProcess(input string) {

	// if there are open system-questions, the user input will be regarded as the response
	// and its goal will be made the active goal
	waitingFor := system.waitingFor
	if waitingFor != nil {
		waitingIsOver := waitingFor.Copy()
		waitingIsOver.Arguments[2] = mentalese.NewTermString(input)
		system.assert(waitingIsOver)
		return
	}

	system.processRunner.StartProcess(
		central.LANGUAGE_PROCESS,
		[]mentalese.Relation{
			mentalese.NewRelation(false, mentalese.PredicateRespond, []mentalese.Term{
				mentalese.NewTermString(input),
			}),
		}, mentalese.NewBinding())
}

func (system *System) readAnswer(message mentalese.RelationSet) (string, *common.Options, bool) {

	answer := ""
	options := common.NewOptions()

	if len(message) == 0 {
		return answer, options, true
	}

	firstRelation := message[0]
	if firstRelation.Predicate == mentalese.PredicatePrint {
		answer = firstRelation.Arguments[1].TermValue
	}

	system.waitingFor = nil
	if firstRelation.Predicate == mentalese.PredicateUserSelect {
		system.waitingFor = &firstRelation
		for i, value := range firstRelation.Arguments[1].TermValueList.GetValues() {
			options.AddOption(strconv.Itoa(i), value)
		}
	}

	return answer, options, true
}
