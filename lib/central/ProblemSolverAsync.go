package central

import (
	"nli-go/lib/api"
	"nli-go/lib/central/goal"
	"nli-go/lib/mentalese"
)

const handleLinkChar = "-"

type ProblemSolverAsync struct {
	solver *ProblemSolver
	relationHandlers      map[string][]api.RelationHandler
}

func NewProblemSolverAsync(solver *ProblemSolver) *ProblemSolverAsync {
	async := ProblemSolverAsync{
		solver: solver,
		relationHandlers:      map[string][]api.RelationHandler{},
	}

	return &async
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
	for _, base := range s.solver.index.factBases {
		rules := base.GetReadMappings()
		for _, rule := range rules {
			s.addRelationHandler(rule.Goal.Predicate, s.createFactBaseClosure(base))
		}
	}
}

func (s *ProblemSolverAsync) createFactBaseClosure(base api.FactBase) api.RelationHandler{
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
		return s.solver.FindFacts(base, relation, binding)
	}
}

func (s *ProblemSolverAsync) createRuleHandlers() {
	for _, base := range s.solver.index.ruleBases {
		for _, rule := range base.GetRules() {
			s.addRelationHandler(rule.Goal.Predicate, s.createRuleClosure(rule))
		}
	}
}

func (s *ProblemSolverAsync) createRuleClosure(rule mentalese.Rule) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

		_, match  := s.solver.matcher.MatchTwoRelations(relation, rule.Goal, binding)
		if !match {
			return mentalese.NewBindingSet()
		}

		mapping, mappingOk := s.solver.matcher.MatchTwoRelations(rule.Goal, relation, mentalese.NewBinding())
		// todo: necessary?
		if !mappingOk {
			return mentalese.NewBindingSet()
		}

		mappedPattern := rule.Pattern.ConvertVariables(mapping, s.solver.variableGenerator)

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
	for _, base := range s.solver.index.functionBases {
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
	for _, base := range s.solver.index.solverFunctionBases {
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
	for _, base := range s.solver.index.factBases {

		for _, mapping := range base.GetWriteMappings() {
			s.addRelationHandler(mentalese.PredicateAssert + handleLinkChar + mapping.Goal.Predicate, s.createAssertFactClosure(base))
			s.addRelationHandler(mentalese.PredicateRetract + handleLinkChar + mapping.Goal.Predicate, s.createRetractFactClosure(base))
		}
	}
}

func (s *ProblemSolverAsync) createAssertFactClosure(base api.FactBase) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
		if relation.Arguments[0].IsRelationSet() {
			localIdBinding := s.solver.replaceSharedIdsByLocalIds(binding, base)
			boundRelation := relation.BindSingle(localIdBinding)
			singleRelation := boundRelation.Arguments[0].TermValueRelationSet[0]
			if singleRelation.IsBound() {
				found := s.solver.modifier.Assert(singleRelation, base)
				if !found {
					return mentalese.NewBindingSet()
				}
			} else {
				s.solver.log.AddError("Cannot assert unbound relation " + singleRelation.String())
				return mentalese.NewBindingSet()
			}
			newBinding := s.solver.replaceLocalIdBySharedId(binding, base)
			return mentalese.InitBindingSet(newBinding)
		} else {
			return mentalese.NewBindingSet()
		}
	}
}

func (s *ProblemSolverAsync) createRetractFactClosure(base api.FactBase) api.RelationHandler {
	return func(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
		if relation.Arguments[0].IsRelationSet() {
			localIdBinding := s.solver.replaceSharedIdsByLocalIds(binding, base)
			boundRelation := relation.BindSingle(localIdBinding)
			found := s.solver.modifier.Retract(boundRelation.Arguments[0].TermValueRelationSet[0], base)
			if !found {
				return mentalese.NewBindingSet()
			}
			newBinding := s.solver.replaceLocalIdBySharedId(binding, base)
			return mentalese.InitBindingSet(newBinding)
		} else {
			return mentalese.NewBindingSet()
		}
	}
}

func (s *ProblemSolverAsync) createRuleBaseModificationHandlers() {

	for _, base := range s.solver.index.ruleBases {
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
			s.solver.index.reindexRules() // todo remove
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

	functions, found := s.solver.index.multiBindingFunctions[relation.Predicate]
	if found {
		for _, function := range functions {
			newBindings = function(messenger, relation, bindings)
			messenger.AddOutBindings(newBindings)
			multiFound = true
		}
	}

	return newBindings, multiFound
}

func (s *ProblemSolverAsync) SolveSingleRelationSingleBinding(messenger *goal.Messenger, relation mentalese.Relation, binding mentalese.Binding) {

	_, found := s.solver.index.knownPredicates[relation.Predicate]
		if !found {
			s.solver.log.AddError("Predicate not supported by any knowledge base: " + relation.Predicate)
			return
		}

	// go through all simple fact bases
	factBases, f4 := s.solver.index.factReadBases[relation.Predicate]
	if f4 {
		for _, factBase := range factBases {
			messenger.AddOutBindings(s.solver.FindFacts(factBase, relation, binding))
		}
	}

	// go through all rule bases
	bases, f3 := s.solver.index.ruleReadBases[relation.Predicate]
	if f3 {
		for _, base := range bases {
			s.solveSingleRelationSingleBindingSingleRuleBase(messenger, relation, binding, base)
		}
	}

	// go through all simple function bases
	functions1, f1 := s.solver.index.simpleFunctions[relation.Predicate]
	if f1 {
		for _, function := range functions1 {
			resultBinding, success := function(relation, binding)
			if success {
				messenger.AddOutBinding(resultBinding)
			}
		}
	}

	// go through all solver functions
	functions2, f2 := s.solver.index.solverFunctions[relation.Predicate]
	if f2 {
		for _, function := range functions2 {
			result := function(messenger, relation, binding)
			messenger.AddOutBindings(result)
		}
	}

	// do assert / retract
	s.solver.modifyKnowledgeBase(relation, binding)
}

func (s *ProblemSolverAsync) solveSingleRelationSingleBindingSingleRuleBase(messenger *goal.Messenger, goalRelation mentalese.Relation, binding mentalese.Binding, ruleBase api.RuleBase) {

	// match rules from the rule base to the goalRelation
	rules := ruleBase.GetRulesForRelation(goalRelation, binding)
	sourceSubgoalSets := []mentalese.RelationSet{}
	for _, rule := range rules {
		aBinding, _ := s.solver.matcher.MatchTwoRelations(goalRelation, rule.Goal, binding)
		bBinding, _ := s.solver.matcher.MatchTwoRelations(rule.Goal, goalRelation, mentalese.NewBinding())
		boundRule := rule.BindSingle(bBinding)
		boundRule = boundRule.InstantiateUnboundVariables(aBinding, s.solver.variableGenerator)
		sourceSubgoalSets = append(sourceSubgoalSets, boundRule.Pattern)
	}

	scopedBinding := mentalese.NewBinding().Merge(binding)

	cursor := messenger.GetCursor()

	// Build the rule index
	currentRuleIndex := cursor.GetState("rule", 0)

	// process child frame bindings
	if currentRuleIndex > 0 {
		messenger.AddOutBindings(cursor.GetChildFrameResultBindings())
	}

	if currentRuleIndex < len(rules) {
		cursor.SetState("rule", currentRuleIndex + 1)
		messenger.CreateChildStackFrame(sourceSubgoalSets[currentRuleIndex], mentalese.InitBindingSet(scopedBinding))
	}
}
