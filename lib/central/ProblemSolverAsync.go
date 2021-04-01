package central

import (
	"fmt"
	"nli-go/lib/api"
	"nli-go/lib/central/goal"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

const handleLinkChar = "-"

type ProblemSolverAsync struct {
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
	modifier              *FactBaseModifier
	log                   *common.SystemLog
}

func NewProblemSolverAsync(matcher *RelationMatcher, log *common.SystemLog) *ProblemSolverAsync {
	variableGenerator := mentalese.NewVariableGenerator()
	async := ProblemSolverAsync{
		factBases:         []api.FactBase{},
		ruleBases:         []api.RuleBase{},
		functionBases:     []api.FunctionBase{},
		multiBindingBases: []api.MultiBindingBase{},
		solverFunctionBases: []api.SolverFunctionBase{},
		simpleFunctions:       map[string][]api.SimpleFunction{},
		multiBindingFunctions: map[string][]api.MultiBindingFunction{},
		matcher: matcher,
		variableGenerator: variableGenerator,
		relationHandlers: map[string][]api.RelationHandler{},
		modifier:          NewFactBaseModifier(log, variableGenerator),
		log: log,
	}

	return &async
}


func (solver *ProblemSolverAsync) AddFactBase(base api.FactBase) {
	solver.factBases = append(solver.factBases, base)
}

func (solver *ProblemSolverAsync) AddFunctionBase(base api.FunctionBase) {
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

func (solver *ProblemSolverAsync) AddRuleBase(base api.RuleBase) {
	solver.ruleBases = append(solver.ruleBases, base)
}

func (solver *ProblemSolverAsync) AddMultipleBindingBase(base api.MultiBindingBase) {
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

func (solver *ProblemSolverAsync) AddSolverFunctionBase(base api.SolverFunctionBase) {
	solver.solverFunctionBases = append(solver.solverFunctionBases, base)
}

func (solver *ProblemSolverAsync) PersistSessionBases() {
	for _, factBase := range solver.factBases {
		switch v := factBase.(type) {
		case api.SessionBasedFactBase:
			v.Persist()
		}
	}
	for _, ruleBase := range solver.ruleBases {
		switch v := ruleBase.(type) {
		case api.SessionBasedFactBase:
			v.Persist()
		}
	}
}

func (solver *ProblemSolverAsync) ResetSession() {
	for _, factBase := range solver.factBases {
		switch v := factBase.(type) {
		case api.SessionBasedFactBase:
			v.ResetSession()
		}
	}
	for _, ruleBase := range solver.ruleBases {
		switch v := ruleBase.(type) {
		case api.SessionBasedFactBase:
			v.ResetSession()
		}
	}
	// relations are indexed by instance; so we need to reindex at least these
	solver.Reindex()
}

func (s *ProblemSolverAsync) Reindex() {

	s.relationHandlers = map[string][]api.RelationHandler{}

	s.createFactBaseHandlers()
	s.createRuleHandlers()
	s.createSimpleFunctionBaseHandlers()
	s.createSolverFunctionBaseHandlers()
	s.createFactBaseModificationHandlers()
	s.createRuleBaseModificationHandlers()
}

func (s *ProblemSolverAsync) addRelationHandler(predicate string, handler api.RelationHandler) {
	_, found := s.relationHandlers[predicate]
	if !found {
		s.relationHandlers[predicate] = []api.RelationHandler{}
	}
	s.relationHandlers[predicate] = append(s.relationHandlers[predicate], handler)
}

func (s *ProblemSolverAsync) createFactBaseHandlers() {
	for _, base := range s.factBases {
		mappings := base.GetReadMappings()
		for _, mapping := range mappings {
			s.addRelationHandler(mapping.Goal.Predicate, s.createFactBaseClosure(base, mapping))
		}
	}
}

func (s *ProblemSolverAsync) createFactBaseClosure(base api.FactBase, mapping mentalese.Rule) api.RelationHandler{
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
		return s.findFactsSingleMapping(base, mapping, relation, binding)
	}
}

func (s *ProblemSolverAsync) createRuleHandlers() {
	for _, base := range s.ruleBases {
		for _, rule := range base.GetRules() {
			s.addRelationHandler(rule.Goal.Predicate, s.createRuleClosure(rule))
		}
	}
}

func (s *ProblemSolverAsync) createRuleClosure(rule mentalese.Rule) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

		_, match  := s.matcher.MatchTwoRelations(relation, rule.Goal, binding)
		if !match {
			return mentalese.NewBindingSet()
		}

		mapping, mappingOk := s.matcher.MatchTwoRelations(rule.Goal, relation, mentalese.NewBinding())
		// todo: necessary?
		if !mappingOk {
			return mentalese.NewBindingSet()
		}

		mappedPattern := rule.Pattern.ConvertVariables(mapping, s.variableGenerator)

		cursor := messenger.GetCursor()
		state := cursor.GetState("state", 0)

		// process child frame bindings
		if state == 0 {
			cursor.SetState("state", 1)
			// turn the cursor into a scope
			cursor.SetType(mentalese.FrameTypeScope)
			// push the child relations
			messenger.CreateChildStackFrame(mappedPattern, mentalese.InitBindingSet(binding))
		} else {
			return cursor.GetChildFrameResultBindings()
		}

		return mentalese.NewBindingSet()
	}
}

func (s *ProblemSolverAsync) createSimpleFunctionBaseHandlers() {
	for _, base := range s.functionBases {
		for predicate, function := range base.GetFunctions() {
			s.addRelationHandler(predicate, s.createSimpleFunctionClosure(function))
		}
	}
}

func (s *ProblemSolverAsync) createSimpleFunctionClosure(function api.SimpleFunction) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
		result, success := function(relation, binding)
		if success {
			return mentalese.InitBindingSet(result)
		} else {
			return mentalese.NewBindingSet()
		}
	}
}

func (s *ProblemSolverAsync) createSolverFunctionBaseHandlers() {
	for _, base := range s.solverFunctionBases {
		for predicate, function := range base.GetFunctions() {
			s.addRelationHandler(predicate, s.createSolverFunctionClosure(function))
		}
	}
}

func (s *ProblemSolverAsync) createSolverFunctionClosure(function api.SolverFunction) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
		return function(messenger, relation, binding)
	}
}

func (s *ProblemSolverAsync) createFactBaseModificationHandlers() {
	for _, base := range s.factBases {

		for _, mapping := range base.GetWriteMappings() {
			s.addRelationHandler(mentalese.PredicateAssert + handleLinkChar + mapping.Goal.Predicate, s.createAssertFactClosure(base))
			s.addRelationHandler(mentalese.PredicateRetract + handleLinkChar + mapping.Goal.Predicate, s.createRetractFactClosure(base))
		}
	}
}

func (s *ProblemSolverAsync) createAssertFactClosure(base api.FactBase) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
		if relation.Arguments[0].IsRelationSet() {
			localIdBinding := s.replaceSharedIdsByLocalIds(binding, base)
			boundRelation := relation.BindSingle(localIdBinding)
			singleRelation := boundRelation.Arguments[0].TermValueRelationSet[0]
			if singleRelation.IsBound() {
				found := s.modifier.Assert(singleRelation, base)
				if !found {
					return mentalese.NewBindingSet()
				}
			} else {
				s.log.AddError("Cannot assert unbound relation " + singleRelation.String())
				return mentalese.NewBindingSet()
			}
			newBinding := s.replaceLocalIdBySharedId(binding, base)
			return mentalese.InitBindingSet(newBinding)
		} else {
			return mentalese.NewBindingSet()
		}
	}
}

func (s *ProblemSolverAsync) createRetractFactClosure(base api.FactBase) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
		if relation.Arguments[0].IsRelationSet() {
			localIdBinding := s.replaceSharedIdsByLocalIds(binding, base)
			boundRelation := relation.BindSingle(localIdBinding)
			found := s.modifier.Retract(boundRelation.Arguments[0].TermValueRelationSet[0], base)
			if !found {
				return mentalese.NewBindingSet()
			}
			newBinding := s.replaceLocalIdBySharedId(binding, base)
			return mentalese.InitBindingSet(newBinding)
		} else {
			return mentalese.NewBindingSet()
		}
	}
}

func (s *ProblemSolverAsync) createRuleBaseModificationHandlers() {

	for _, base := range s.ruleBases {
		for _, predicate := range base.GetWritablePredicates() {
			s.addRelationHandler(mentalese.PredicateAssert + handleLinkChar + predicate, s.createAssertRuleClosure(base))
		}
	}
}

func (s *ProblemSolverAsync) createAssertRuleClosure(base api.RuleBase) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
		if relation.Arguments[0].IsRule() {
			rule := relation.Arguments[0].TermValueRule.BindSingle(binding)
			base.Assert(rule)
			s.Reindex()
			return mentalese.InitBindingSet(binding)
		} else {
			return mentalese.NewBindingSet()
		}
	}
}

func (s *ProblemSolverAsync) GetHandlers(relation mentalese.Relation) []api.RelationHandler {

	handle := relation.Predicate

	if handle == mentalese.PredicateAssert || handle == mentalese.PredicateRetract {
		object := relation.Arguments[0]
		if object.IsRule() {
			handle += handleLinkChar + relation.Arguments[0].TermValueRule.Goal.Predicate
		} else {
			handle += handleLinkChar + relation.Arguments[0].TermValueRelationSet[0].Predicate
		}
	}

	handlers, found := s.relationHandlers[handle]

	if found {
		return handlers
	} else {
		return []api.RelationHandler{}
	}
}

func (s *ProblemSolverAsync) SolveMultipleBindings(messenger *goal.Messenger, relation mentalese.Relation, bindings mentalese.BindingSet) (mentalese.BindingSet, bool) {

	newBindings := mentalese.NewBindingSet()
	multiFound := false

	functions, found := s.multiBindingFunctions[relation.Predicate]
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
func (solver *ProblemSolverAsync) FindFacts(factBase api.FactBase, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	dbBindings := mentalese.NewBindingSet()

	for _, ds2db := range factBase.GetReadMappings() {
		mappingBindings := solver.findFactsSingleMapping(factBase, ds2db, relation, binding)
		dbBindings.AddMultiple(mappingBindings)
	}

	return dbBindings
}

func (solver *ProblemSolverAsync) findFactsSingleMapping(factBase api.FactBase, ds2db mentalese.Rule, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	dbBindings := mentalese.NewBindingSet()

	activeBinding, match := solver.matcher.MatchTwoRelations(relation, ds2db.Goal, mentalese.NewBinding())
	if !match { return dbBindings }

	activeBinding2, match2 := solver.matcher.MatchTwoRelations(ds2db.Goal, relation, mentalese.NewBinding())
	if !match2 { return dbBindings }

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

func (solver *ProblemSolverAsync) solveMultipleRelationSingleFactBase(relations []mentalese.Relation, binding mentalese.Binding, factBase api.FactBase) mentalese.BindingSet {

	sequenceBindings := mentalese.InitBindingSet(binding)

	for _, relation := range relations {
		sequenceBindings = solver.solveSingleRelationSingleFactBase(relation, sequenceBindings, factBase)
	}

	return sequenceBindings
}

func (solver *ProblemSolverAsync) solveSingleRelationSingleFactBase(relation mentalese.Relation, bindings mentalese.BindingSet, factBase api.FactBase) mentalese.BindingSet {

	// todo: notice cannot handle second order predicates, since no handler

	relationBindings, multiFound := solver.SolveMultipleBindings(nil, relation, bindings)
	resultBindings := mentalese.NewBindingSet()

	if !multiFound {

		for _, binding := range bindings.GetAll() {

			// todo: notice cannot handle second order predicates, since no handler

			_, found := solver.simpleFunctions[relation.Predicate]
			if found {
				handlers := solver.GetHandlers(relation)
				resultBindings = handlers[0](nil, relation, binding)
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

func (solver *ProblemSolverAsync) replaceSharedIdsByLocalIds(binding mentalese.Binding, factBase api.FactBase) mentalese.Binding {

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

func (solver *ProblemSolverAsync) replaceLocalIdBySharedId(binding mentalese.Binding, factBase api.FactBase) mentalese.Binding {

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