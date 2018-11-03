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
func (answerer Answerer) Answer(goal mentalese.RelationSet) mentalese.RelationSet {

	answerer.log.StartDebug("Answer")

	answer := mentalese.RelationSet{}
	transformer := mentalese.NewRelationTransformer(answerer.matcher, answerer.log)

	// conditionBindings: map condition variables to goal variables
	allSolutions, allConditionBindings := answerer.findSolutions(goal)

	if len(allSolutions) == 0 {

		answerer.log.AddError("Answerer could not find a solution.")

	} else {

		for i, solution := range allSolutions {

			conditionBindings := allConditionBindings[i]

			// add transformation variables
			conditionBindings = answerer.addVariablesFromTransformationReplacements(conditionBindings, solution.Transformations)

			// apply transformation, if available
			transformedGoal := transformer.Replace(solution.Transformations, goal)

			// scope here, just before finding the solution
			quantifierScoper := mentalese.NewQuantifierScoper(answerer.log)
			transformedGoal = quantifierScoper.Scope(transformedGoal)

			// resultBindings: map goal variables to answers
			resultBindings := answerer.solver.SolveRelationSet(transformedGoal, []mentalese.Binding{{}})

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

			// stop after the first solution
			break
		}
	}

	answerer.log.EndDebug("Answer", answer)
	return answer
}

func (answerer Answerer) addVariablesFromTransformationReplacements(conditionBindings []mentalese.Binding, transformations []mentalese.RelationTransformation) []mentalese.Binding {

	var newBindings []mentalese.Binding
	for _, binding := range conditionBindings {

		newBinding := binding.Copy()

		for _, transformation := range transformations {
			for _, relation := range transformation.Replacement {
				for _, argument := range relation.Arguments {
					if argument.TermType == mentalese.Term_variable {
						if !binding.ContainsVariable(argument.TermValue) {
							newBinding = newBinding.Merge(mentalese.Binding{argument.TermValue: mentalese.Term{TermType: mentalese.Term_variable, TermValue: argument.TermValue}})
						}
					}
				}
			}
		}

		newBindings = append(newBindings, newBinding)
	}

	return newBindings
}

// Returns the solutions whose condition matches the goal, and a set of bindings per solution
func (answerer Answerer) findSolutions(goal mentalese.RelationSet) ([]mentalese.Solution, [][]mentalese.Binding) {

	answerer.log.StartDebug("findSolutions", goal)

	var solutions []mentalese.Solution
	var allBindings [][]mentalese.Binding

	for _, aSolution := range answerer.solutions {

		unScopedGoal := goal.UnScope()

		bindings, found := answerer.matcher.MatchSequenceToSet(aSolution.Condition, unScopedGoal, mentalese.Binding{})
		if found {
			solutions = append(solutions, aSolution)
			allBindings = append(allBindings, bindings)
		}
	}

	answerer.log.EndDebug("findSolutions", solutions, allBindings)

	return solutions, allBindings
}
