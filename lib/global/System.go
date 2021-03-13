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

func (system *System) CreateAnswerGoal(input string) string {

	uuid := common.CreateUuid()

	// go:assert(go:goal(go:respond(input, Id)))
	set := mentalese.RelationSet{
		mentalese.NewRelation(true, mentalese.PredicateAssert, []mentalese.Term{
			mentalese.NewTermRelationSet(mentalese.RelationSet{
				mentalese.NewRelation(true, mentalese.PredicateGoal, []mentalese.Term{
					mentalese.NewTermRelationSet(mentalese.RelationSet{
						mentalese.NewRelation(true, mentalese.PredicateRespond, []mentalese.Term{
							mentalese.NewTermString(input),
						}),
					}),
					mentalese.NewTermString(uuid),
				})}),
		}),
	}
	system.processRunner.RunRelationSet(set)

	return uuid
}

func (system *System) ReadActions(actionType string) mentalese.BindingSet {
	set := mentalese.RelationSet{
		mentalese.NewRelation(true, mentalese.PredicateAction, []mentalese.Term{
			mentalese.NewTermString(actionType),
			mentalese.NewTermVariable("Id"),
			mentalese.NewTermVariable("Content"),
		}),
	}

	return system.processRunner.RunRelationSet(set)
}

func (system *System) DeleteAction(actionId string) {
	set := mentalese.RelationSet{
		mentalese.NewRelation(true, mentalese.PredicateRetract, []mentalese.Term{
			mentalese.NewTermRelationSet(mentalese.RelationSet{
				mentalese.NewRelation(true, mentalese.PredicateAction, []mentalese.Term{
					mentalese.NewTermAnonymousVariable(),
					mentalese.NewTermString(actionId),
					mentalese.NewTermAnonymousVariable(),
				})}),
		}),
	}
	system.processRunner.RunRelationSet(set)
}

func (system *System) DeleteGoal(goalId string) {
	set := mentalese.RelationSet{
		mentalese.NewRelation(true, mentalese.PredicateRetract, []mentalese.Term{
			mentalese.NewTermRelationSet(mentalese.RelationSet{
				mentalese.NewRelation(true, mentalese.PredicateGoal, []mentalese.Term{
					mentalese.NewTermAnonymousVariable(),
					mentalese.NewTermString(goalId),
				})}),
		}),
	}
	system.processRunner.RunRelationSet(set)
}

func (system *System) Answer(input string) (string, *common.Options) {

	// find or create a goal
	goalId := system.getGoalId(input)

	// execute all goals
	system.Run()

	// get the goal's process
	process := system.processList.GetProcess(goalId)

	// build options for the user, if applicable
	options := system.buildOptions(process)

	// read answer action
	answer := system.readAnswer()

	// store the updated dialog context
	system.dialogContext.Store()

	return answer, options
}

func (system *System) getGoalId(input string) string {

	goalId := ""

	for _, process := range system.processList.GetProcesses() {
		if !process.IsDone() {

			userSelect := process.GetLastFrame().Relations
			binding := mentalese.NewBinding()
			binding.Set("Selection", mentalese.NewTermString(input))
			list := userSelect[0].Arguments[0]

			set := mentalese.RelationSet{
				mentalese.NewRelation(true, mentalese.PredicateAssert, []mentalese.Term{
					mentalese.NewTermRelationSet(
						mentalese.RelationSet{
							mentalese.NewRelation(true, mentalese.PredicateUserSelect, []mentalese.Term{
								list,
								mentalese.NewTermString(input),
							}),
						}),
				}),
			}
			system.processRunner.RunRelationSet(set)

			goalId = process.GoalId

			break
		}
	}

	if goalId == "" {
		goalId = system.CreateAnswerGoal(input)
	}

	return goalId
}

func (system *System) Run() {

	// find all goals
	set := mentalese.RelationSet{
		mentalese.NewRelation(true, mentalese.PredicateGoal, []mentalese.Term{
			mentalese.NewTermVariable("Goal"),
			mentalese.NewTermVariable("Id"),
		}),
	}
	bindings := system.processRunner.RunRelationSet(set)

	// go through all goals
	for _, binding := range bindings.GetAll() {
		goalId := binding.MustGet("Id").TermValue
		goalSet := binding.MustGet("Goal").TermValueRelationSet

		// run the process
		process := system.processList.GetOrCreateProcess(goalId, goalSet)
		system.processRunner.RunProcess(process)

		// delete goal when done
		if process.IsDone() {
			system.DeleteGoal(goalId)
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

func (system *System) readAnswer() string {
	actions := system.ReadActions(mentalese.ActionPrint)
	answer := ""
	if actions.GetLength() > 0 {
		action := actions.Get(0)
		answer = action.MustGet("Content").TermValue
		actionId := action.MustGet("Id").TermValue
		system.DeleteAction(actionId)
	}
	return answer
}

func (system *System) ResetSession() {
	system.dialogContext.Initialize()
	system.dialogContext.Store()

	system.solverAsync.ResetSession()
}
