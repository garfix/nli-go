package central

import (
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
)

type Answerer struct {
	solutions []mentalese.Solution
	matcher *mentalese.RelationMatcher
	solver *ProblemSolver
}

func NewAnswerer(matcher *mentalese.RelationMatcher) *Answerer {
	return &Answerer{solutions: []mentalese.Solution{}, matcher: matcher, solver: NewProblemSolver(matcher)}
}

func (answerer *Answerer) AddFactBase(source knowledge.FactBase) {
	answerer.solver.AddFactBase(source)
}

func (answerer *Answerer) AddRuleBase(source knowledge.RuleBase) {
	answerer.solver.AddRuleBase(source)
}

func (answerer *Answerer) AddMultipleBindingsBase(source knowledge.MultipleBindingsBase) {
	answerer.solver.AddMultipleBindingsBase(source)
}

func (answerer *Answerer) AddSolutions(solutions []mentalese.Solution) {
	answerer.solutions = append(answerer.solutions, solutions...)
}

// goal e.g. [ question(Q) child(S, O) name(S, 'Janice', fullName) numberOf(N, O) focus(Q, N) ]
// return e.g. [ child(S, O) gender(S, female) numberOf(N, O) ]
func (answerer Answerer) Answer(goal mentalese.RelationSet) mentalese.RelationSet {

	common.LogTree("Answer")

	answers := []mentalese.RelationSet{}

	// conditionBindings: map condition variables to goal variables
	solution, conditionBindings, found := answerer.findSolution(goal)
	if found {

		// resultBindings: map goal variables to answers
		resultBindings := answerer.solver.SolveRelationSet(goal, []mentalese.Binding{})

		// solutionBindings: map condition variables to results
		solutionBindings := []mentalese.Binding{}
		for _, conditionBinding := range conditionBindings {
			for _, resultBinding := range resultBindings {
				solutionBindings = append(solutionBindings, conditionBinding.Bind(resultBinding))
			}
		}

		// extend solution bindings by executing the preparation
		if !solution.Preparation.IsEmpty() {
			solutionBindings = answerer.solver.SolveRelationSet(solution.Preparation, solutionBindings)
		}

		// create answers relation sets by binding 'answer' to solutionBindings
		answers = answerer.matcher.BindRelationSetMultipleBindings(solution.Answer, solutionBindings)
	}

	singleAnswer := mentalese.RelationSet{}
	for _, answer := range answers {
		singleAnswer = singleAnswer.Merge(answer)
	}

	common.LogTree("Answer", singleAnswer)
	return singleAnswer
}

// Returns the solution whose condition matches the goal
func (answerer Answerer) findSolution(goal mentalese.RelationSet) (mentalese.Solution, []mentalese.Binding, bool) {

	common.LogTree("findSolution", goal)

	solution := mentalese.Solution{}
	bindings := []mentalese.Binding{}
	found := false

	for _, aSolution := range answerer.solutions  {

		bindings, _, found = answerer.matcher.MatchSequenceToSet(aSolution.Condition, goal, mentalese.Binding{})
		if found {
			solution = aSolution
			break
		}
	}

	common.LogTree("findSolution", solution, bindings, found)

	return solution, bindings, found
}
