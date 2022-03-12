package global

// waiting on another thread
// https://medium.com/@matryer/golang-advent-calendar-day-two-starting-and-stopping-things-with-a-signal-channel-f5048161018

import (
	"fmt"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"strconv"
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
	discourseEntities     *mentalese.Binding
	solverAsync           *central.ProblemSolverAsync
	answerer              *central.Answerer
	generator             *generate.Generator
	surfacer              *generate.SurfaceRepresentation
	processList           *central.ProcessList
	processRunner         *central.ProcessRunner
}

func (system *System) GetLog() *common.SystemLog {
	return system.log
}

// Low-level function to inspect the internal state of the system
func (system *System) Query(relations string) mentalese.BindingSet {
	set := system.internalGrammarParser.CreateRelationSet(relations)
	result := system.processRunner.RunRelationSet(set)

	system.persistSession()

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

	done := make(chan struct{})

	println("")
	println("in: " + clientMessage.String())

	callback := func(serverMessage mentalese.RelationSet) {

		fmt.Println("server: " + serverMessage.String())

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
	system.processRunner.RunRelationSet(clientMessage)

	// wait until done
	<-done

	fmt.Println("done!")

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

func (system *System) ResetSession() {
	system.dialogContext.Initialize()
	system.solverAsync.ResetSession()

	system.persistSession()
}

func (system *System) persistSession() {
	system.dialogContext.Store()
	system.solverAsync.PersistSessionBases()
}

func (system *System) assert(relation mentalese.Relation) {

	// go:assert()
	set := mentalese.RelationSet{
		mentalese.NewRelation(false, mentalese.PredicateAssert, []mentalese.Term{
			mentalese.NewTermRelationSet(mentalese.RelationSet{relation}),
		}),
	}
	system.processRunner.RunRelationSet(set)
}

func (system *System) createOrUpdateProcess(input string) {

	// if there are open system-questions, the user input will be regarded as the response
	// and its goal will be made the active goal
	for _, process := range system.processList.GetProcesses() {
		waitingFor := process.GetWaitingFor()
		if waitingFor != nil {
			waitingIsOver := waitingFor[0].Copy()
			waitingIsOver.Arguments[2] = mentalese.NewTermString(input)
			system.assert(waitingIsOver)
			return
		}
	}

	//system.createAnswerGoal(input)
	system.processRunner.StartProcess([]mentalese.Relation{
		mentalese.NewRelation(false, mentalese.PredicateRespond, []mentalese.Term{
			mentalese.NewTermString(input),
		}),
	}, mentalese.NewBinding())
}

func (system *System) readAnswer(message mentalese.RelationSet) (string, *common.Options, bool) {

	answer := ""
	options := common.NewOptions()

	firstRelation := message[0]
	if firstRelation.Predicate == mentalese.PredicatePrint {
		answer = firstRelation.Arguments[1].TermValue
	}
	if firstRelation.Predicate == mentalese.PredicateUserSelect {
		for i, value := range firstRelation.Arguments[1].TermValueList.GetValues() {
			options.AddOption(strconv.Itoa(i), value)
		}
	}

	return answer, options, true
}
