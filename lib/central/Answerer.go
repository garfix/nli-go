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

func NewAnswerer() *Answerer {
	return &Answerer{solutions: []mentalese.Solution{}, matcher:mentalese.NewRelationMatcher(), solver: NewProblemSolver()}
}

func (answerer *Answerer) AddKnowledgeBase(source knowledge.KnowledgeBase) {
	answerer.solver.AddKnowledgeBase(source)
}

func (answerer *Answerer) AddMultipleBindingsBase(source knowledge.MultipleBindingsBase) {
	answerer.solver.AddMultipleBindingsBase(source)
}

func (solver *Answerer) AddSolutions(solutions []mentalese.Solution) {
	solver.solutions = append(solver.solutions, solutions...)
}

// goal e.g. [ question(Q) child(S, O) name(S, 'Janice', fullName) numberOf(O, N) focus(Q, N) ]
// return e.g. [ child(S, O) gender(S, female) numberOf(O, N) ]
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

	return solution, bindings, found
}
