package central

import (
	"fmt"
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

const handleLinkChar = "-"

type ProblemSolver struct {
	factBases             []api.FactBase
	ruleBases             []api.RuleBase
	functionBases         []api.FunctionBase
	multiBindingBases     []api.MultiBindingBase
	solverFunctionBases   []api.SolverFunctionBase
	simpleFunctions       map[string][]api.SimpleFunction
	multiBindingFunctions map[string][]api.MultiBindingFunction
	matcher               *RelationMatcher
	variableGenerator     *mentalese.VariableGenerator
	relationHandlers      map[string][]api.RelationHandler
	functions             map[string]mentalese.Rule
	modifier              *FactBaseModifier
	log                   *common.SystemLog
}

func NewProblemSolver(matcher *RelationMatcher, variableGenerator *mentalese.VariableGenerator, log *common.SystemLog) *ProblemSolver {
	solver := ProblemSolver{
		factBases:             []api.FactBase{},
		ruleBases:             []api.RuleBase{},
		functionBases:         []api.FunctionBase{},
		multiBindingBases:     []api.MultiBindingBase{},
		solverFunctionBases:   []api.SolverFunctionBase{},
		simpleFunctions:       map[string][]api.SimpleFunction{},
		multiBindingFunctions: map[string][]api.MultiBindingFunction{},
		matcher:               matcher,
		variableGenerator:     variableGenerator,
		relationHandlers:      map[string][]api.RelationHandler{},
		functions:             map[string]mentalese.Rule{},
		log:                   log,
	}

	return &solver
}

func (solver *ProblemSolver) SetModifier(modifier *FactBaseModifier) {
	solver.modifier = modifier
}

func (solver *ProblemSolver) AddFactBase(base api.FactBase) {
	solver.factBases = append(solver.factBases, base)
}

func (solver *ProblemSolver) AddFunctionBase(base api.FunctionBase) {
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

func (solver *ProblemSolver) AddRuleBase(base api.RuleBase) {
	solver.ruleBases = append(solver.ruleBases, base)
}

func (solver *ProblemSolver) AddMultipleBindingBase(base api.MultiBindingBase) {
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

func (solver *ProblemSolver) AddSolverFunctionBase(base api.SolverFunctionBase) {
	solver.solverFunctionBases = append(solver.solverFunctionBases, base)
}

func (solver *ProblemSolver) Reindex() {

	solver.relationHandlers = map[string][]api.RelationHandler{}

	solver.createFactBaseHandlers()
	solver.createRuleHandlers()
	solver.createSimpleFunctionBaseHandlers()
	solver.createSolverFunctionBaseHandlers()
	solver.createFactBaseModificationHandlers()
	solver.createRuleBaseModificationHandlers()
}

func (solver *ProblemSolver) addRelationHandler(predicate string, handler api.RelationHandler) {
	_, found := solver.relationHandlers[predicate]
	if !found {
		solver.relationHandlers[predicate] = []api.RelationHandler{}
	}
	solver.relationHandlers[predicate] = append(solver.relationHandlers[predicate], handler)
}

func (solver *ProblemSolver) createFactBaseHandlers() {
	for _, base := range solver.factBases {
		mappings := base.GetReadMappings()
		for _, mapping := range mappings {
			solver.addRelationHandler(mapping.Goal.Predicate, solver.createFactBaseClosure(base, mapping))
		}
	}
}

func (solver *ProblemSolver) createFactBaseClosure(base api.FactBase, mapping mentalese.Rule) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
		return solver.findFactsSingleMapping(base, mapping, relation, binding)
	}
}

func (solver *ProblemSolver) createRuleHandlers() {
	for _, base := range solver.ruleBases {
		for _, rule := range base.GetRules() {
			if rule.IsFunction {
				solver.functions[rule.Goal.Predicate] = rule
			}
			solver.addRelationHandler(rule.Goal.Predicate, solver.createRuleClosure(rule))
		}
	}
}

func (solver *ProblemSolver) createRuleClosure(rule mentalese.Rule) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

		b1, match := solver.matcher.MatchTwoRelations(relation, rule.Goal, binding)
		if !match {
			return mentalese.NewBindingSet()
		}

		mapping, mappingOk := solver.matcher.MatchTwoRelations(rule.Goal, relation, mentalese.NewBinding())
		// todo: necessary?
		if !mappingOk {
			return mentalese.NewBindingSet()
		}

		binding = binding.Merge(b1.RemoveVariables())

		mappedPattern := rule.Pattern.ConvertVariables(mapping, solver.variableGenerator)

		cursor := messenger.GetCursor()
		cursor.SetType(mentalese.FrameTypeScope)

		newBindings := messenger.ExecuteChildStackFrame(mappedPattern, mentalese.InitBindingSet(binding))
		return newBindings

	}
}

func (solver *ProblemSolver) createSimpleFunctionBaseHandlers() {
	for _, base := range solver.functionBases {
		for predicate, function := range base.GetFunctions() {
			solver.addRelationHandler(predicate, solver.createSimpleFunctionClosure(function))
		}
	}
}

func (solver *ProblemSolver) createSimpleFunctionClosure(function api.SimpleFunction) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
		result, success := function(messenger, relation, binding)
		if success {
			return mentalese.InitBindingSet(result)
		} else {
			return mentalese.NewBindingSet()
		}
	}
}

func (solver *ProblemSolver) createSolverFunctionBaseHandlers() {
	for _, base := range solver.solverFunctionBases {
		for predicate, function := range base.GetFunctions() {
			solver.addRelationHandler(predicate, solver.createSolverFunctionClosure(function))
		}
	}
}

func (solver *ProblemSolver) createSolverFunctionClosure(function api.SolverFunction) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
		return function(messenger, relation, binding)
	}
}

func (solver *ProblemSolver) createFactBaseModificationHandlers() {
	for _, base := range solver.factBases {

		for _, mapping := range base.GetWriteMappings() {
			solver.addRelationHandler(mentalese.PredicateAssert+handleLinkChar+mapping.Goal.Predicate, solver.createAssertFactClosure(base))
			solver.addRelationHandler(mentalese.PredicateRetract+handleLinkChar+mapping.Goal.Predicate, solver.createRetractFactClosure(base))
		}
	}
}

func (solver *ProblemSolver) createAssertFactClosure(base api.FactBase) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
		if solver.modifier == nil {
			return mentalese.NewBindingSet()
		}
		if relation.Arguments[0].IsRelationSet() {
			localIdBinding := solver.replaceSharedIdsByLocalIds(binding, base)
			boundRelation := relation.BindSingle(localIdBinding)
			singleRelation := boundRelation.Arguments[0].TermValueRelationSet[0]
			if singleRelation.IsBound() {
				found := solver.modifier.Assert(singleRelation, base)
				if !found {
					return mentalese.NewBindingSet()
				}
			} else {
				solver.log.AddError("Cannot assert unbound relation " + singleRelation.String())
				return mentalese.NewBindingSet()
			}
			newBinding := solver.replaceLocalIdBySharedId(binding, base)
			return mentalese.InitBindingSet(newBinding)
		} else {
			return mentalese.NewBindingSet()
		}
	}
}

func (solver *ProblemSolver) createRetractFactClosure(base api.FactBase) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
		if solver.modifier == nil {
			return mentalese.NewBindingSet()
		}
		if relation.Arguments[0].IsRelationSet() {
			localIdBinding := solver.replaceSharedIdsByLocalIds(binding, base)
			boundRelation := relation.BindSingle(localIdBinding)
			found := solver.modifier.Retract(boundRelation.Arguments[0].TermValueRelationSet[0], base)
			if !found {
				return mentalese.NewBindingSet()
			}
			newBinding := solver.replaceLocalIdBySharedId(binding, base)
			return mentalese.InitBindingSet(newBinding)
		} else {
			return mentalese.NewBindingSet()
		}
	}
}

func (solver *ProblemSolver) createRuleBaseModificationHandlers() {

	for _, base := range solver.ruleBases {
		for _, predicate := range base.GetWritablePredicates() {
			solver.addRelationHandler(mentalese.PredicateAssert+handleLinkChar+predicate, solver.createAssertRuleClosure(base))
		}
	}
}

func (solver *ProblemSolver) createAssertRuleClosure(base api.RuleBase) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
		if relation.Arguments[0].IsRule() {
			rule := relation.Arguments[0].TermValueRule.BindSingle(binding)
			base.Assert(rule)
			solver.Reindex()
			return mentalese.InitBindingSet(binding)
		} else {
			return mentalese.NewBindingSet()
		}
	}
}

func (solver *ProblemSolver) GetHandlers(relation mentalese.Relation) []api.RelationHandler {

	handle := relation.Predicate

	if handle == mentalese.PredicateAssert || handle == mentalese.PredicateRetract {
		object := relation.Arguments[0]
		if object.IsRule() {
			handle += handleLinkChar + relation.Arguments[0].TermValueRule.Goal.Predicate
		} else {
			handle += handleLinkChar + relation.Arguments[0].TermValueRelationSet[0].Predicate
		}
	}

	handlers, found := solver.relationHandlers[handle]

	if found {
		return handlers
	} else {
		return []api.RelationHandler{}
	}
}

func (solver *ProblemSolver) SolveMultipleBindings(messenger *Messenger, relation mentalese.Relation, bindings mentalese.BindingSet) (mentalese.BindingSet, bool) {

	newBindings := mentalese.NewBindingSet()
	multiFound := false

	functions, found := solver.multiBindingFunctions[relation.Predicate]
	if found {
		for _, function := range functions {
			newBindings = function(messenger, relation, bindings)
			multiFound = true
			break
		}
	}

	return newBindings, multiFound
}

// Creates bindings for the free variables in 'relations', by resolving them in factBase
func (solver *ProblemSolver) FindFacts(factBase api.FactBase, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	dbBindings := mentalese.NewBindingSet()

	for _, ds2db := range factBase.GetReadMappings() {
		mappingBindings := solver.findFactsSingleMapping(factBase, ds2db, relation, binding)
		dbBindings.AddMultiple(mappingBindings)
	}

	return dbBindings
}

func (solver *ProblemSolver) findFactsSingleMapping(factBase api.FactBase, ds2db mentalese.Rule, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	dbBindings := mentalese.NewBindingSet()

	activeBinding, match := solver.matcher.MatchTwoRelations(relation, ds2db.Goal, mentalese.NewBinding())
	if !match {
		return dbBindings
	}

	activeBinding2, match2 := solver.matcher.MatchTwoRelations(ds2db.Goal, relation, mentalese.NewBinding())
	if !match2 {
		return dbBindings
	}

	dbRelations := ds2db.Pattern.ConvertVariables(activeBinding2, solver.variableGenerator)

	localIdBinding := solver.replaceSharedIdsByLocalIds(binding, factBase)

	relevantBinding := localIdBinding.Select(dbRelations.GetVariableNames())
	newDbBindings := solver.solveMultipleRelationSingleFactBase(dbRelations, relevantBinding, factBase)

	for _, newDbBinding := range newDbBindings.GetAll() {

		dbBinding := activeBinding.Merge(newDbBinding)

		combinedBinding := localIdBinding.Merge(dbBinding.Select(relation.GetVariableNames()))
		sharedBinding := solver.replaceLocalIdBySharedId(combinedBinding, factBase)
		dbBindings.Add(sharedBinding)
	}

	return dbBindings
}

func (solver *ProblemSolver) solveMultipleRelationSingleFactBase(relations []mentalese.Relation, binding mentalese.Binding, factBase api.FactBase) mentalese.BindingSet {

	sequenceBindings := mentalese.InitBindingSet(binding)

	for _, relation := range relations {
		sequenceBindings = solver.solveSingleRelationSingleFactBase(relation, sequenceBindings, factBase)
	}

	return sequenceBindings
}

func (solver *ProblemSolver) solveSingleRelationSingleFactBase(relation mentalese.Relation, bindings mentalese.BindingSet, factBase api.FactBase) mentalese.BindingSet {

	// todo: notice cannot handle second order predicates, since no handler

	relationBindings, multiFound := solver.SolveMultipleBindings(nil, relation, bindings)
	resultBindings := mentalese.NewBindingSet()

	if !multiFound {

		for _, binding := range bindings.GetAll() {

			// todo: notice cannot handle second order predicates, since no handler

			_, found := solver.simpleFunctions[relation.Predicate]
			if found {
				handlers := solver.GetHandlers(relation)
				simpleMessenger := NewSimpleMessenger()
				resultBindings = handlers[0](simpleMessenger, relation, binding)
			} else {
				resultBindings = factBase.MatchRelationToDatabase(relation, binding)
			}

			// found bindings must be extended with the bindings already present
			for _, resultBinding := range resultBindings.GetAll() {
				newRelationBinding := binding.Merge(resultBinding)
				relationBindings.Add(newRelationBinding)
			}
		}
	}

	return relationBindings
}

func (solver *ProblemSolver) replaceSharedIdsByLocalIds(binding mentalese.Binding, factBase api.FactBase) mentalese.Binding {

	newBinding := mentalese.NewBinding()

	for key, value := range binding.GetAll() {
		newValue := value

		if value.IsId() {
			sharedId := value.TermValue
			sort := value.TermSort
			if sort != "" {
				localId := factBase.GetLocalId(sharedId, sort)
				if localId == "" {
					solver.log.AddError(fmt.Sprintf("Local id %s not found for %s in fact base %s", sharedId, sort, factBase.GetName()))
					return newBinding
				}
				newValue = mentalese.NewTermId(localId, sort)
			}
		}

		newBinding.Set(key, newValue)
	}

	return newBinding
}

func (solver *ProblemSolver) replaceLocalIdBySharedId(binding mentalese.Binding, factBase api.FactBase) mentalese.Binding {

	newBinding := mentalese.NewBinding()

	for key, value := range binding.GetAll() {
		newValue := value

		if value.IsId() {
			localId := value.TermValue
			sort := value.TermSort
			if sort != "" {
				sharedId := factBase.GetSharedId(localId, sort)
				if sharedId == "" {
					solver.log.AddError(fmt.Sprintf("Shared id %s not found for %s in fact base %s", localId, sort, factBase.GetName()))
					return newBinding
				}
				newValue = mentalese.NewTermId(sharedId, sort)
			}
		}

		newBinding.Set(key, newValue)
	}

	return newBinding
}
