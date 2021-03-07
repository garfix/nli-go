package central

import (
	"nli-go/lib/api"
)

type KnowledgeBaseIndex struct {
	factBases             []api.FactBase
	ruleBases             []api.RuleBase
	functionBases         []api.FunctionBase
	multiBindingBases     []api.MultiBindingBase
	solverFunctionBases   []api.SolverFunctionBase
	simpleFunctions       map[string][]api.SimpleFunction
	multiBindingFunctions map[string][]api.MultiBindingFunction
}

func NewProblemSolverIndex() *KnowledgeBaseIndex {
	return &KnowledgeBaseIndex{
		factBases:         []api.FactBase{},
		ruleBases:         []api.RuleBase{},
		functionBases:     []api.FunctionBase{},
		multiBindingBases: []api.MultiBindingBase{},
		solverFunctionBases: []api.SolverFunctionBase{},
		simpleFunctions:       map[string][]api.SimpleFunction{},
		multiBindingFunctions: map[string][]api.MultiBindingFunction{},
	}
}

func (solver *KnowledgeBaseIndex) AddFactBase(base api.FactBase) {
	solver.factBases = append(solver.factBases, base)
}

func (solver *KnowledgeBaseIndex) AddFunctionBase(base api.FunctionBase) {
	solver.functionBases = append(solver.functionBases, base)
	functions := base.GetFunctions()
	for predicate, function := range functions {
		_, found := solver.simpleFunctions[predicate]
		if !found {
			solver.simpleFunctions[predicate] = []api.SimpleFunction{}
		}
		solver.simpleFunctions[predicate] = append(solver.simpleFunctions[predicate], function)
	}
}

func (solver *KnowledgeBaseIndex) AddRuleBase(base api.RuleBase) {
	solver.ruleBases = append(solver.ruleBases, base)
}

func (solver *KnowledgeBaseIndex) AddMultipleBindingBase(base api.MultiBindingBase) {
	solver.multiBindingBases = append(solver.multiBindingBases, base)
	functions := base.GetFunctions()
	for predicate, function := range functions {
		_, found := solver.multiBindingFunctions[predicate]
		if !found {
			solver.multiBindingFunctions[predicate] = []api.MultiBindingFunction{}
		}
		solver.multiBindingFunctions[predicate] = append(solver.multiBindingFunctions[predicate], function)
	}
}

func (solver *KnowledgeBaseIndex) AddSolverFunctionBase(base api.SolverFunctionBase) {
	solver.solverFunctionBases = append(solver.solverFunctionBases, base)
}
