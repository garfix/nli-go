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

// goal e.g. [ question(Q) child(S, O) EntityId(S, 'Janice', fullName) number_of(N, O) focus(Q, N) ]
// return e.g. [ child(S, O) gender(S, female) number_of(N, O) ]
func (answerer Answerer) Answer(goal mentalese.RelationSet, keyCabinet *mentalese.KeyCabinet) mentalese.RelationSet {

	answerer.log.StartDebug("Answer")

	answer := mentalese.RelationSet{}
	transformer := mentalese.NewRelationTransformer(answerer.matcher, answerer.log)

	// scope here, just before finding the solution
	quantifierScoper := mentalese.NewQuantifierScoper(answerer.log)
	scopedGoal := quantifierScoper.Scope(goal)

	answerer.log.AddProduction("Scoped", scopedGoal.String())

	// apply sequences
	sequenceApplier := mentalese.NewSequenceApplier(answerer.log)
	sequencedGoal := sequenceApplier.ApplySequences(scopedGoal)

	answerer.log.AddProduction("With Sequences", sequencedGoal.String())

	// conditionBindings: map condition variables to goal variables
	allSolutions := answerer.findSolutions(sequencedGoal)

	if len(allSolutions) == 0 {

		answerer.log.AddError("Answerer could not find a solution.")

	} else {

		for i, solution := range allSolutions {

			// apply transformation, if available
			transformedGoal := transformer.Replace(solution.Transformations, sequencedGoal)

			// resultBindings: map goal variables to answers
			resultBindings := answerer.solver.SolveRelationSet(transformedGoal, keyCabinet, mentalese.Bindings{{}})

			// choose a handler based on whether there were results
			resultHandler := solution.NoResults
			if len(resultBindings) > 0 {
				resultHandler = solution.SomeResults
			} else {
				// no results? try the next solution (if there is one)
				if i < len(allSolutions) - 1 {
					continue
				}
			}

			// solutionBindings: map condition variables to results
			var solutionBindings = resultBindings

			// extend solution bindings by executing the preparation
			if !resultHandler.Preparation.IsEmpty() {
				solutionBindings = answerer.solver.SolveRelationSet(resultHandler.Preparation, keyCabinet, resultBindings)
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
