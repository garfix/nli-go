package central

import (
	"nli-go/lib/api"
	"nli-go/lib/mentalese"
)

type ProblemSolverAsync struct {
	solver *ProblemSolver
}

func NewProblemSolverAsync(solver *ProblemSolver) *ProblemSolverAsync {
	return &ProblemSolverAsync{
		solver: solver,
	}
}

func (s *ProblemSolverAsync) SolveMultipleBindings(messenger api.ProcessMessenger, relation mentalese.Relation, bindings mentalese.BindingSet) (mentalese.BindingSet, bool) {

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

//func (s *ProblemSolverAsync) SolveSingleRelationSingleBinding(messenger api.ProcessMessenger) {
func (s *ProblemSolverAsync) SolveSingleRelationSingleBinding(messenger api.ProcessMessenger, relation mentalese.Relation, binding mentalese.Binding) {

	//relation := messenger.GetRelation()
	//binding := messenger.GetInBinding()

	_, found := s.solver.index.knownPredicates[relation.Predicate]
		if !found {
			s.solver.log.AddError("Predicate not supported by any knowledge base: " + relation.Predicate)
			return
		}

	// go through all simple fact bases
	factBases, f4 := s.solver.index.factReadBases[relation.Predicate]
	if f4 {
		for _, factBase := range factBases {
			// todo
			s.solver.FindFacts(factBase, relation, binding)
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
	// todo
	s.solver.modifyKnowledgeBase(relation, binding)
}

func (s *ProblemSolverAsync) solveSingleRelationSingleBindingSingleRuleBase(messenger api.ProcessMessenger, goalRelation mentalese.Relation, binding mentalese.Binding, ruleBase api.RuleBase) {

	// match rules from the rule base to the goalRelation
	rules := ruleBase.GetRules(goalRelation, binding)
	sourceSubgoalSets := []mentalese.RelationSet{}
	for _, rule := range rules {
		aBinding, _ := s.solver.matcher.MatchTwoRelations(goalRelation, rule.Goal, binding)
		bBinding, _ := s.solver.matcher.MatchTwoRelations(rule.Goal, goalRelation, mentalese.NewBinding())
		boundRule := rule.BindSingle(bBinding)
		boundRule = boundRule.InstantiateUnboundVariables(aBinding, s.solver.variableGenerator)
		sourceSubgoalSets = append(sourceSubgoalSets, boundRule.Pattern)
	}

	scope := mentalese.NewScope()
	scopedBinding := mentalese.NewScopedBinding(scope).Merge(binding)

	cursor := messenger.GetCursor()

	// build the rule index
	currentRuleIndex := cursor.GetState("rule", 0)

	// process child frame bindings
	if currentRuleIndex > 0 {
	//	cursor.AddStepBindings(cursor.GetChildFrameResultBindings())
		messenger.AddOutBindings(cursor.GetChildFrameResultBindings())
	}

	if currentRuleIndex < len(rules) {
		cursor.SetState("rule", currentRuleIndex + 1)
		messenger.CreateChildStackFrame(sourceSubgoalSets[currentRuleIndex], mentalese.InitBindingSet(scopedBinding))
	}
}
