package central

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

// The answerer takes a relation set in domain format
// and returns a relation set in domain format
// It uses Solution structures to determine how to act
type Answerer struct {
	solutions []mentalese.Solution
	matcher   *mentalese.RelationMatcher
	solver    *ProblemSolver
	builder   *RelationSetBuilder
	log       *common.SystemLog
}

func NewAnswerer(matcher *mentalese.RelationMatcher, solver *ProblemSolver, log *common.SystemLog) *Answerer {

	builder := NewRelationSetBuilder()
	builder.addGenerator(NewSystemGenerator())

	return &Answerer{
		solutions: []mentalese.Solution{},
		matcher:   matcher,
		solver:    solver,
		builder:   builder,
		log:       log,
	}
}

func (answerer *Answerer) AddSolutions(solutions []mentalese.Solution) {
	answerer.solutions = append(answerer.solutions, solutions...)
}

// goal e.g. [ question(Q) child(S, O) name(S, 'Janice', fullName) number_of(N, O) focus(Q, N) ]
// return e.g. [ child(S, O) gender(S, female) number_of(N, O) ]
func (answerer Answerer) Answer(goal mentalese.RelationSet) mentalese.RelationSet {

	answerer.log.StartDebug("Answer")

	answer := mentalese.RelationSet{}

	// conditionBindings: map condition variables to goal variables
	solution, conditionBindings, found := answerer.findSolution(goal)
	if !found {

		answerer.log.AddError("Answerer could not find a solution.")

	} else {

		// resultBindings: map goal variables to answers
		resultBindings := answerer.solver.SolveRelationSet(goal, []mentalese.Binding{{}})

		// choose a handler based on whether there were results
		resultHandler := solution.NoResults
		if len(resultBindings) > 0 {
			resultHandler = solution.SomeResults
		}

		// solutionBindings: map condition variables to results
		var solutionBindings []mentalese.Binding
		for _, conditionBinding := range conditionBindings {
			for _, resultBinding := range resultBindings {
				solutionBindings = append(solutionBindings, conditionBinding.Bind(resultBinding))
			}
		}

		// extend solution bindings by executing the preparation
		if !resultHandler.Preparation.IsEmpty() {
			solutionBindings = answerer.solver.SolveRelationSet(resultHandler.Preparation, solutionBindings)
		}

		// create answer relation sets by binding 'answer' to solutionBindings
		answer = answerer.builder.Build(resultHandler.Answer, solutionBindings)
	}

	answerer.log.EndDebug("Answer", answer)
	return answer
}

// Returns the solution whose condition matches the goal
func (answerer Answerer) findSolution(goal mentalese.RelationSet) (mentalese.Solution, []mentalese.Binding, bool) {

	answerer.log.StartDebug("findSolution", goal)

	solution := mentalese.Solution{}
	bindings := []mentalese.Binding{}
	found := false

	for _, aSolution := range answerer.solutions {

		unScopedGoal := goal.UnScope()

		bindings, found = answerer.matcher.MatchSequenceToSet(aSolution.Condition, unScopedGoal, mentalese.Binding{})
		if found {
			solution = aSolution
			break
		}
	}

	answerer.log.EndDebug("findSolution", solution, bindings, found)

	return solution, bindings, found
}
