package central

import (
	"nli-go/lib/api"
	"nli-go/lib/mentalese"
)

type KnowledgeBaseIndex struct {
	factBases             []api.FactBase
	ruleBases             []api.RuleBase
	functionBases         []api.FunctionBase
	multiBindingBases     []api.MultiBindingBase
	solverFunctionBases   []api.SolverFunctionBase
	knownPredicates       map[string]bool
	solverFunctions       map[string][]api.SolverFunction
	simpleFunctions       map[string][]api.SimpleFunction
	multiBindingFunctions map[string][]api.MultiBindingFunction
	ruleReadBases         map[string][]api.RuleBase
	ruleWriteBases        map[string][]api.RuleBase
	factReadBases         map[string][]api.FactBase
	factWriteBases		  map[string][]api.FactBase
}

func NewProblemSolverIndex() *KnowledgeBaseIndex {
	return &KnowledgeBaseIndex{
		factBases:         []api.FactBase{},
		ruleBases:         []api.RuleBase{},
		functionBases:     []api.FunctionBase{},
		multiBindingBases: []api.MultiBindingBase{},
		solverFunctionBases: []api.SolverFunctionBase{},
		knownPredicates:       map[string]bool{},
		solverFunctions:       map[string][]api.SolverFunction{},
		simpleFunctions:       map[string][]api.SimpleFunction{},
		multiBindingFunctions: map[string][]api.MultiBindingFunction{},
		ruleReadBases:         map[string][]api.RuleBase{},
		ruleWriteBases:         map[string][]api.RuleBase{},
		factReadBases:         map[string][]api.FactBase{},
		factWriteBases:         map[string][]api.FactBase{},
	}
}

func (solver *KnowledgeBaseIndex) AddFactBase(base api.FactBase) {
	solver.factBases = append(solver.factBases, base)
	rules := base.GetReadMappings()
	for _, rule := range rules {
		predicate := rule.Goal.Predicate
		solver.knownPredicates[predicate] = true
		_, found := solver.factReadBases[predicate]
		if !found {
			solver.factReadBases[predicate] = []api.FactBase{}
		}
		solver.factReadBases[predicate] = append(solver.factReadBases[predicate], base)
	}
	rules = base.GetWriteMappings()
	if len(rules) > 0 {
		solver.knownPredicates[mentalese.PredicateAssert] = true
		solver.knownPredicates[mentalese.PredicateRetract] = true
	}
	for _, rule := range rules {
		predicate := rule.Goal.Predicate
		solver.knownPredicates[predicate] = true
		_, found := solver.factReadBases[predicate]
		if !found {
			solver.factWriteBases[predicate] = []api.FactBase{}
		}
		solver.factWriteBases[predicate] = append(solver.factWriteBases[predicate], base)
	}
}

func (solver *KnowledgeBaseIndex) AddFunctionBase(base api.FunctionBase) {
	solver.functionBases = append(solver.functionBases, base)
	functions := base.GetFunctions()
	for predicate, function := range functions {
		solver.knownPredicates[predicate] = true
		_, found := solver.simpleFunctions[predicate]
		if !found {
			solver.simpleFunctions[predicate] = []api.SimpleFunction{}
		}
		solver.simpleFunctions[predicate] = append(solver.simpleFunctions[predicate], function)
	}
}

func (solver *KnowledgeBaseIndex) AddRuleBase(base api.RuleBase) {
	solver.ruleBases = append(solver.ruleBases, base)
	predicates := base.GetPredicates()
	if len(predicates) > 0 {
		solver.knownPredicates[mentalese.PredicateAssert] = true
		solver.knownPredicates[mentalese.PredicateRetract] = true
	}
	for _, predicate := range predicates {
		solver.knownPredicates[predicate] = true
		_, found := solver.ruleReadBases[predicate]
		if !found {
			solver.ruleReadBases[predicate] = []api.RuleBase{}
		}
		solver.ruleReadBases[predicate] = append(solver.ruleReadBases[predicate], base)
	}
}

func (solver *KnowledgeBaseIndex) AddMultipleBindingBase(base api.MultiBindingBase) {
	solver.multiBindingBases = append(solver.multiBindingBases, base)
	functions := base.GetFunctions()
	for predicate, function := range functions {
		solver.knownPredicates[predicate] = true
		_, found := solver.multiBindingFunctions[predicate]
		if !found {
			solver.multiBindingFunctions[predicate] = []api.MultiBindingFunction{}
		}
		solver.multiBindingFunctions[predicate] = append(solver.multiBindingFunctions[predicate], function)
	}
}

func (solver *KnowledgeBaseIndex) AddSolverFunctionBase(base api.SolverFunctionBase) {
	solver.solverFunctionBases = append(solver.solverFunctionBases, base)
	functions := base.GetFunctions()
	for predicate, function := range functions {
		solver.knownPredicates[predicate] = true
		_, found := solver.solverFunctions[predicate]
		if !found {
			solver.solverFunctions[predicate] = []api.SolverFunction{}
		}
		solver.solverFunctions[predicate] = append(solver.solverFunctions[predicate], function)
	}
}

func (solver *KnowledgeBaseIndex) reindexRules() {
	// todo: should really reindex knownPredicates too
	solver.ruleReadBases = map[string][]api.RuleBase{}
	for _, base := range solver.ruleBases {
		predicates := base.GetPredicates()
		for _, predicate := range predicates {
			solver.knownPredicates[predicate] = true
			_, found := solver.ruleReadBases[predicate]
			if !found {
				solver.ruleReadBases[predicate] = []api.RuleBase{}
			}
			solver.ruleReadBases[predicate] = append(solver.ruleReadBases[predicate], base)
		}
	}
}