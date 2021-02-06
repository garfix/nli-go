package central

import (
	"nli-go/lib/api"
	"nli-go/lib/central/goal"
	"nli-go/lib/mentalese"
	"strconv"
)

type ProblemSolverAsync struct {
	solver *ProblemSolver
}

func NewProblemSolverAsync(solver *ProblemSolver) *ProblemSolverAsync {
	return &ProblemSolverAsync{
		solver: solver,
	}
}

func (s *ProblemSolverAsync) solveMultipleBindings(relation mentalese.Relation, bindings mentalese.BindingSet) (mentalese.BindingSet, bool) {

	newBindings := mentalese.NewBindingSet()
	multiFound := false

	functions, found := s.solver.index.multiBindingFunctions[relation.Predicate]
	if found {
		for _, function := range functions {
			newBindings = function(relation, bindings)
			multiFound = true
		}
	}

	return newBindings, multiFound
}

func (s *ProblemSolverAsync) SolveSingleRelationSingleBinding(messenger *goal.Messenger) {

	relation := messenger.GetRelation()
	binding := messenger.GetInBinding()

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
			result := function(relation, binding)
			messenger.AddOutBindings(result)
		}
	}

	// do assert / retract
	// todo
	s.solver.modifyKnowledgeBase(relation, binding)
}

func (s *ProblemSolverAsync) solveSingleRelationSingleBindingSingleRuleBase(messenger *goal.Messenger, goalRelation mentalese.Relation, binding mentalese.Binding, ruleBase api.RuleBase) {

	subgoalResultBindings := mentalese.BindingSet{}
	inputVariables := goalRelation.GetVariableNames()

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
	currentRuleIndex := 0
	ruleBinding, ruleBindingFound := cursor.State.Get("rule")
	if ruleBindingFound {
		currentRuleIndex, _ = ruleBinding.GetIntValue()
	}

	// process child frame bindings
	if currentRuleIndex > 0 {
		subgoalResultBindings = cursor.StepBindings
		for _, childResult := range cursor.ChildFrameResultBindings.GetAll() {
// todo: don't need to do this here
			// filter out the input variables
			filteredBinding := childResult.FilterVariablesByName(inputVariables)
			// make sure all variables of the original binding are present
			goalBinding := scopedBinding.Merge(filteredBinding)
			subgoalResultBindings.Add(goalBinding)
		}
		cursor.StepBindings = subgoalResultBindings
	}

	if currentRuleIndex < len(rules) {
		cursor.State.Set("rule", mentalese.NewTermString(strconv.Itoa(currentRuleIndex + 1)))
		messenger.CreateChildStackFrame(sourceSubgoalSets[currentRuleIndex], mentalese.InitBindingSet(scopedBinding))
	}
}
