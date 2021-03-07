package central

import (
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

// The problem solver takes a relation set and a set of bindings
// and returns a set of new bindings
// It uses knowledge bases to find these bindings
type ProblemSolver struct {
	index    			  *KnowledgeBaseIndex
	matcher               *RelationMatcher
	variableGenerator     *mentalese.VariableGenerator
	modifier              *FactBaseModifier
	dialogContext         *DialogContext
	log                   *common.SystemLog
}

func NewProblemSolver(matcher *RelationMatcher, dialogContext *DialogContext, log *common.SystemLog) *ProblemSolver {
	variableGenerator := mentalese.NewVariableGenerator()
	return &ProblemSolver{
		index: 			   NewProblemSolverIndex(),
		variableGenerator: variableGenerator,
		modifier:          NewFactBaseModifier(log, variableGenerator),
		matcher:           matcher,
		dialogContext:     dialogContext,
		log:               log,
	}
}

func (solver *ProblemSolver) AddFactBase(base api.FactBase) {
	solver.index.AddFactBase(base)
}

func (solver *ProblemSolver) AddFunctionBase(base api.FunctionBase) {
	solver.index.AddFunctionBase(base)
}

func (solver *ProblemSolver) AddRuleBase(base api.RuleBase) {
	solver.index.AddRuleBase(base)
}

func (solver *ProblemSolver) AddMultipleBindingBase(base api.MultiBindingBase) {
	solver.index.AddMultipleBindingBase(base)
}

func (solver *ProblemSolver) AddSolverFunctionBase(base api.SolverFunctionBase) {
	solver.index.AddSolverFunctionBase(base)
}

func (solver *ProblemSolver) ResetSession() {
	for _, factBase := range solver.index.factBases {
		switch v := factBase.(type) {
		case api.SessionBasedFactBase:
			v.ResetSession()
		}
	}
}

