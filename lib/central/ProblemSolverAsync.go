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

//func (s *ProblemSolverAsync) SolveRelation(process *goal.Process) {
//
//	relation := process.GetLastFrame().
//
//	if s.solver.log.Active() { s.solver.log.StartDebug("Solve Relation", relation.String() + " " + fmt.Sprint(bindings)) }
//
//	_, found := s.solver.index.knownPredicates[relation.Predicate]
//	if !found {
//		s.solver.log.AddError("Predicate not supported by any knowledge base: " + relation.Predicate)
//		return
//	}
//
//	newBindings, multiFound := s.solveMultipleBindings(relation, bindings)
//
//	if !multiFound {
//		for _, binding := range bindings.GetAll() {
//			newBindings.AddMultiple(s.SolveSingleRelationSingleBinding(process, relation, binding))
//		}
//	}
//
//	if s.solver.log.Active() { s.solver.log.EndDebug("Solve Relation", relation.String() + ": " + fmt.Sprint(newBindings)) }
//
//	return newBindings
//}

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

func (s *ProblemSolverAsync) SolveSingleRelationSingleBinding(process *goal.Process) {

	frame := process.GetLastFrame()
	relation := frame.GetCurrentRelation()
	binding := frame.GetCurrentBinding()

	_, found := s.solver.index.knownPredicates[relation.Predicate]
		if !found {
			s.solver.log.AddError("Predicate not supported by any knowledge base: " + relation.Predicate)
			return
		}

	relationVariables := relation.GetVariableNames()
	simpleBinding := binding.FilterVariablesByName(relationVariables)

	s.solver.callStack.PushSingle(relation, binding)

	newBindings := mentalese.NewBindingSet()

	// go through all simple fact bases
	factBases, f4 := s.solver.index.factReadBases[relation.Predicate]
	if f4 {
		for _, factBase := range factBases {
			newBindings.AddMultiple(s.solver.FindFacts(factBase, relation, simpleBinding))
		}
	}

	// go through all rule bases
	bases, f3 := s.solver.index.ruleReadBases[relation.Predicate]
	if f3 {
		for _, base := range bases {
			s.solveSingleRelationSingleBindingSingleRuleBase(process, relation, simpleBinding, base)
		}
	}

	// go through all simple function bases
	functions1, f1 := s.solver.index.simpleFunctions[relation.Predicate]
	if f1 {
		for _, function := range functions1 {
			resultBinding, success := function(relation, simpleBinding)
			if success {
				newBindings.Add(resultBinding)
			}
		}
	}

	// go through all solver functions
	functions2, f2 := s.solver.index.solverFunctions[relation.Predicate]
	if f2 {
		for _, function := range functions2 {
			result := function(relation, simpleBinding)
			frame.OutBindings.AddMultiple(result)
		}
	}

	// do assert / retract
	newBindings.AddMultiple(s.solver.modifyKnowledgeBase(relation, simpleBinding))

	// compose the result set
	completedBindings := mentalese.NewBindingSet()
	for _, newBinding := range newBindings.GetAll() {
		// remove temporary variables
		essentialResultBinding := newBinding.FilterVariablesByName(relationVariables)
		// combine the source binding with the clean results
		completedBinding := binding.Merge(essentialResultBinding)
		completedBindings.Add(completedBinding)
	}

	s.solver.callStack.Pop(newBindings)
}

func (s *ProblemSolverAsync) solveSingleRelationSingleBindingSingleRuleBase(process *goal.Process, goalRelation mentalese.Relation, binding mentalese.Binding, ruleBase api.RuleBase) {

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

	// build the rule index
	currentRuleIndex := 0
	ruleBinding, ruleBindingFound := process.GetCursor().State.Get("rule")
	if ruleBindingFound {
		currentRuleIndex, _ = ruleBinding.GetIntValue()
	}

	// process child frame bindings
	if currentRuleIndex > 0 {
		subgoalResultBindings = process.GetCursor().OutBindings
		for _, childResult := range process.GetCursor().ChildFrameResultBindings.GetAll() {
			// filter out the input variables
			filteredBinding := childResult.FilterVariablesByName(inputVariables)
			// make sure all variables of the original binding are present
			goalBinding := scopedBinding.Merge(filteredBinding)
			subgoalResultBindings.Add(goalBinding)
		}
		process.GetCursor().OutBindings = subgoalResultBindings
	}

	if currentRuleIndex < len(rules) {
		process.GetCursor().State.Set("rule", mentalese.NewTermString(strconv.Itoa(currentRuleIndex + 1)))
		process.PushFrame(sourceSubgoalSets[currentRuleIndex], mentalese.InitBindingSet(scopedBinding))
	}

	//for i, sourceSubgoalSet := range sourceSubgoalSets {
	//
	//	scope := mentalese.NewScope()
	//	s.solver.scopeStack.PushFrame(scope)
	//
	//	scopedBinding := mentalese.NewScopedBinding(scope).Merge(binding)
	//	subgoalResultBindings := mentalese.InitBindingSet(scopedBinding)
	//
	//	binding := mentalese.NewBinding()
	//	binding.Set("rule", mentalese.NewTermString(strconv.Itoa(i)))
	//	process.PushFrame(sourceSubgoalSet, binding)
	//
	//	//for _, subGoal := range sourceSubgoalSet {
	//	//
	//	//	subgoalResultBindings = solver.SolveRelationSet([]mentalese.Relation{subGoal}, subgoalResultBindings)
	//	//
	//	//	if subgoalResultBindings.IsEmpty() {
	//	//		break
	//	//	}
	//	//}
	//
	//	for _, subgoalResultBinding := range subgoalResultBindings.GetAll() {
	//
	//		// filter out the input variables
	//		filteredBinding := subgoalResultBinding.FilterVariablesByName(inputVariables)
	//
	//		// make sure all variables of the original binding are present
	//		goalBinding := scopedBinding.Merge(filteredBinding)
	//
	//		goalBindings.Add(goalBinding)
	//	}
	//
	//	s.solver.scopeStack.Pop()
	//}

}
