package global

import (
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
	solver                *central.ProblemSolver
	solverAsync           *central.ProblemSolverAsync
	answerer              *central.Answerer
	generator             *generate.Generator
	surfacer              *generate.SurfaceRepresentation
	processRunner         *central.ProcessRunner
}

// Low-level function to inspect the internal state of the system
func (system *System) Query(relations string) mentalese.BindingSet {
	set := system.internalGrammarParser.CreateRelationSet(relations)
	result := system.solver.SolveRelationSet(set, mentalese.InitBindingSet( mentalese.NewBinding()))

	system.dialogContext.Store()

	return result
}

func (system *System) CreateAnswerGoal(input string) string {

	uuid := common.CreateUuid()

	// go:assert(go:goal(go:answer(input, Id)))
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
	system.solver.SolveRelationSet(set, mentalese.InitBindingSet(mentalese.NewBinding()))

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

	return system.solver.SolveRelationSet(set, mentalese.InitBindingSet(mentalese.NewBinding()))
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
	system.solver.SolveRelationSet(set, mentalese.InitBindingSet(mentalese.NewBinding()))
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
	system.solver.SolveRelationSet(set, mentalese.InitBindingSet(mentalese.NewBinding()))
}

func (system *System) Run() {
	// find all goals
	set := mentalese.RelationSet{
		mentalese.NewRelation(true, mentalese.PredicateGoal, []mentalese.Term{
			mentalese.NewTermVariable("Goal"),
			mentalese.NewTermVariable("Id"),
		}),
	}
	// find processes
	bindings := system.solver.SolveRelationSet(set, mentalese.InitBindingSet(mentalese.NewBinding()))
	for _, binding := range bindings.GetAll() {
		goalId := binding.MustGet("Id").TermValue
		goalSet := binding.MustGet("Goal").TermValueRelationSet
		system.processRunner.RunProcess(goalId, goalSet)
	}
}

func (system *System) Answer(input string) (string, *common.Options) {

	options := common.NewOptions()
	goalId := ""

	for _, process := range system.processRunner.GetProcesses() {
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
			system.solver.SolveRelationSet(set, mentalese.InitBindingSet(mentalese.NewBinding()))

			goalId = process.GoalId

			break
		}
	}

	if goalId == "" {
		goalId = system.CreateAnswerGoal(input)
	}

	system.Run()

	process := system.processRunner.GetProcessByGoalId(goalId)
	lastFrame := process.GetLastFrame()
	if lastFrame != nil {
		relation := lastFrame.Relations[0]
		if relation.Predicate == mentalese.PredicateUserSelect {
			for i, value := range relation.Arguments[0].TermValueList.GetValues() {
				options.AddOption(strconv.Itoa(i), value)
			}
		}
	} else {
		system.DeleteGoal(goalId)
	}

	actions := system.ReadActions(mentalese.ActionPrint)
	answer := ""
	if actions.GetLength() > 0 {
		action := actions.Get(0)
		answer = action.MustGet("Content").TermValue
		actionId := action.MustGet("Id").TermValue
		system.DeleteAction(actionId)
	}

	return answer, options
}

func (system *System) ResetSession() {
	system.dialogContext.Initialize()
	system.dialogContext.Store()

	system.solver.ResetSession()
}
