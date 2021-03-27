package global

import (
	"nli-go/lib/central"
	"nli-go/lib/central/goal"
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
	solverAsync           *central.ProblemSolverAsync
	answerer              *central.Answerer
	generator             *generate.Generator
	surfacer              *generate.SurfaceRepresentation
	processList	          *goal.ProcessList
	processRunner         *central.ProcessRunner
}

// Low-level function to inspect the internal state of the system
func (system *System) Query(relations string) mentalese.BindingSet {
	set := system.internalGrammarParser.CreateRelationSet(relations)
	result := system.processRunner.RunRelationSet(set)

	system.dialogContext.Store()

	return result
}

// API call that executes input, runs all processes, and returns all waiting messages
func (system *System) SendMessage(relations mentalese.RelationSet) []mentalese.RelationSet {

	system.processRunner.RunRelationSet(relations)

	system.run()

	relationSets := system.getWaitingMessages()

	system.dialogContext.Store()

	return relationSets
}

// processes input and return the answer, respond to all waiting relations
func (system *System) Answer(input string) (string, *common.Options) {

	answer := ""
	anAnswer := ""
	options := common.NewOptions()
	someOptions := common.NewOptions()
	done := false

	// find or create a goal
	system.getGoalId(input)

	for !done {

		// execute all goals
		system.run()

		if len(system.log.GetErrors()) > 0 {
			break
		}

		// read answer action
		anAnswer, someOptions, done = system.readAnswer()
		if anAnswer != "" {
			answer = anAnswer
		}
		if someOptions.HasOptions() {
			options = someOptions
		}
	}

	// store the updated dialog context
	system.dialogContext.Store()

	return answer, options
}

func (system *System) ResetSession() {
	system.dialogContext.Initialize()
	system.dialogContext.Store()

	system.solverAsync.ResetSession()
}

// returns the children of all `wait_for` relations
func (system *System) getWaitingMessages() []mentalese.RelationSet {
	kessages := []mentalese.RelationSet{}

	for _, process := range system.processList.GetProcesses() {
		beforeLastFrame := process.GetBeforeLastFrame()
		if beforeLastFrame != nil {
			if beforeLastFrame.Relations[beforeLastFrame.RelationIndex].Predicate == mentalese.PredicateWaitFor {
				lastFrame := process.GetLastFrame()
				binding := lastFrame.InBindings.Get(lastFrame.InBindingIndex)
				boundRelations := lastFrame.Relations.BindSingle(binding)
				kessages = append(kessages, boundRelations)
			}
		}
	}

	return kessages
}

func (system *System) assert(relation mentalese.Relation) {

	// go:assert(go:goal(go:respond(input, Id)))
	set := mentalese.RelationSet{
		mentalese.NewRelation(false, mentalese.PredicateAssert, []mentalese.Term{
			mentalese.NewTermRelationSet(mentalese.RelationSet{relation}),
		}),
	}
	system.processRunner.RunRelationSet(set)
}

func (system *System) createAnswerGoal(input string) string {

	uuid := common.CreateUuid()

	system.assert(mentalese.NewRelation(false, mentalese.PredicateGoal, []mentalese.Term{
		mentalese.NewTermRelationSet(mentalese.RelationSet{
			mentalese.NewRelation(false, mentalese.PredicateRespond, []mentalese.Term{
				mentalese.NewTermString(input),
			}),
		}),
		mentalese.NewTermString(uuid),
	}))

	return uuid
}

func (system *System) deleteGoal(goalId string) {

	system.processList.RemoveProcess(goalId)

	set := mentalese.RelationSet{
		mentalese.NewRelation(false, mentalese.PredicateRetract, []mentalese.Term{
			mentalese.NewTermRelationSet(mentalese.RelationSet{
				mentalese.NewRelation(false, mentalese.PredicateGoal, []mentalese.Term{
					mentalese.NewTermAnonymousVariable(),
					mentalese.NewTermString(goalId),
				})}),
		}),
	}
	system.processRunner.RunRelationSet(set)
}

func (system *System) getGoalId(input string) string {

	goalId := ""

	for _, process := range system.processList.GetProcesses() {
		beforeLastFrame := process.GetBeforeLastFrame()
		if beforeLastFrame != nil {
			if beforeLastFrame.Relations[beforeLastFrame.RelationIndex].Predicate == mentalese.PredicateWaitFor {
				lastFrame := process.GetLastFrame()
				if lastFrame.Relations[0].Predicate == mentalese.PredicateUserSelect {
					binding := lastFrame.InBindings.Get(lastFrame.InBindingIndex)
					boundRelation := lastFrame.Relations[0].BindSingle(binding)
					boundRelation.Arguments[1] = mentalese.NewTermString(input)
					system.assert(boundRelation)
					goalId = process.GoalId
					break
				}
			}
		}
	}

	if goalId == "" {
		goalId = system.createAnswerGoal(input)
	}

	return goalId
}

func (system *System) getAllGoals() ([]string, []mentalese.RelationSet) {
	sets := []mentalese.RelationSet{}
	ids := []string{}
	goalString := ""

	// find all goals
	set := mentalese.RelationSet{
		mentalese.NewRelation(false, mentalese.PredicateGoal, []mentalese.Term{
			mentalese.NewTermVariable("Goal"),
			mentalese.NewTermVariable("Id"),
		}),
	}
	bindings := system.processRunner.RunRelationSet(set)

	for _, binding := range bindings.GetAll() {
		goalId := binding.MustGet("Id").TermValue
		goalSet := binding.MustGet("Goal").TermValueRelationSet
		sets = append(sets, goalSet)
		ids = append(ids, goalId)
		goalString += goalId + ":\n"
		for _, relation := range goalSet {
			goalString += relation.String() + "\n"
		}
		goalString += "\n"
	}

	system.log.AddProduction("Goals", goalString)

	return ids, sets
}

func (system *System) run() {

	goalIds, goalSets := system.getAllGoals()

	// go through all goals
	for i, goalId := range goalIds {
		goalSet := goalSets[i]

		// run the process
		process := system.processList.GetOrCreateProcess(goalId, goalSet)
		system.processRunner.RunProcess(process)

		// delete goal when done or failed
		if process.IsDone() || !system.log.IsOk() {
			system.deleteGoal(goalId)
		}
	}
}

func (system *System) buildOptions(process *goal.Process) *common.Options {
	options := common.NewOptions()

	if !process.IsDone() {
		lastFrame := process.GetLastFrame()
		relation := lastFrame.Relations[0]
		if relation.Predicate == mentalese.PredicateUserSelect {
			for i, value := range relation.Arguments[0].TermValueList.GetValues() {
				options.AddOption(strconv.Itoa(i), value)
			}
		}
	}

	return options
}

func (system *System) readAnswer() (string, *common.Options, bool) {

	relationSets := system.getWaitingMessages()
	answer := ""
	options := common.NewOptions()

	done := len(relationSets) == 0

	for _, relationSet := range relationSets {
		firstRelation := relationSet[0]
		if firstRelation.Predicate == mentalese.PredicatePrint {
			answer = firstRelation.Arguments[1].TermValue
		}
		if firstRelation.Predicate == mentalese.PredicateUserSelect {
			for i, value := range firstRelation.Arguments[0].TermValueList.GetValues() {
				options.AddOption(strconv.Itoa(i), value)
			}
			done = true
			continue
		}

		for _, relation := range relationSet {
			system.assert(relation)
		}
	}

	return answer, options, done
}
