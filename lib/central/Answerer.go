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

// goal e.g. [ question(Q) child(S, O) SharedId(S, 'Janice', fullName) number_of(O, N) focus(Q, N) ]
// return e.g. [ child(S, O) gender(S, female) number_of(O, N) ]
func (answerer Answerer) Answer(goal mentalese.RelationSet, bindings mentalese.Bindings) mentalese.RelationSet {

	answerer.log.StartDebug("Answer")

// todo remove
	keyCabinet := mentalese.NewKeyCabinet()

	answer := mentalese.RelationSet{}
	transformer := mentalese.NewRelationTransformer(answerer.matcher, answerer.log)

	// conditionBindings: map condition variables to goal variables
	allSolutions := answerer.findSolutions(goal)

	if len(allSolutions) == 0 {

		answerer.log.AddError("There are no solutions for this problem")

	} else {

		for i, solution := range allSolutions {

			// apply transformation, if available
			transformedGoal := transformer.Replace(solution.Transformations, goal)

			// resultBindings: map goal variables to answers
			resultBindings := answerer.solver.SolveRelationSet(transformedGoal, keyCabinet, bindings)

			// no results? try the next solution (if there is one)
			if len(resultBindings) == 0 {
				if i < len(allSolutions) - 1 {
					continue
				}
			}

			// find a handler
			condionedBindings := resultBindings
			var resultHandler *mentalese.ResultHandler
			for _, response := range solution.Responses {
				if !response.Condition.IsEmpty() {
					conditionBindings := answerer.solver.SolveRelationSet(response.Condition, keyCabinet, resultBindings)
					if len(conditionBindings) == 0 {
						continue
					} else {
						condionedBindings = conditionBindings
					}
				}
				resultHandler = &response
				break
			}

			if resultHandler == nil {
				answerer.log.AddError("No solution had its conditions fulfilled")
				break
			}

			// solutionBindings: map condition variables to results
			var solutionBindings = condionedBindings

			// extend solution bindings by executing the preparation
			if !resultHandler.Preparation.IsEmpty() {
				solutionBindings = answerer.solver.SolveRelationSet(resultHandler.Preparation, keyCabinet, condionedBindings)
			}

			// create answer relation sets by binding 'answer' to solutionBindings
			answer = answerer.builder.Build(resultHandler.Answer, solutionBindings)

			// stop after the first solution
			break
		}
	}

	answerer.log.EndDebug("Answer", answer)
	return answer
}

// Returns the solutions whose condition matches the goal, and a set of bindings per solution
func (answerer Answerer) findSolutions(goal mentalese.RelationSet) []mentalese.Solution {

	answerer.log.StartDebug("findSolutions", goal)

	var solutions []mentalese.Solution

	for _, aSolution := range answerer.solutions {

		unScopedGoal := goal.UnScope()

		bindings, found := answerer.matcher.MatchSequenceToSet(aSolution.Condition, unScopedGoal, mentalese.Binding{})
		if found {

			for _, binding := range bindings {
				boundSolution := aSolution.BindSingle(binding)
				solutions = append(solutions, boundSolution)
			}
		}
	}

	answerer.log.EndDebug("findSolutions", solutions)

	return solutions
}
