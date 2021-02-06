package central

import (
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

// The answerer takes a relation set in domain format
// and returns a relation set in domain format
// It uses Solution structures to determine how to act
type Answerer struct {
	solutions []mentalese.Solution
	matcher   *RelationMatcher
	solver    *ProblemSolver
	solverAsync *ProblemSolverAsync
	log       *common.SystemLog
}

func NewAnswerer(matcher *RelationMatcher, solver *ProblemSolver, solverAsync *ProblemSolverAsync, log *common.SystemLog) *Answerer {

	return &Answerer{
		solutions: []mentalese.Solution{},
		matcher:   matcher,
		solver:    solver,
		solverAsync: solverAsync,
		log:       log,
	}
}

func (answerer *Answerer) AddSolutions(solutions []mentalese.Solution) {
	answerer.solutions = append(answerer.solutions, solutions...)
}

// goal e.g. [ question(Q) child(S, O) SharedId(S, 'Janice', fullName) number_of(O, N) focus(Q, N) ]
// return e.g. [ child(S, O) gender(S, female) number_of(O, N) ]
func (answerer Answerer) Answer(messenger api.ProcessMessenger, goal mentalese.RelationSet, bindings mentalese.BindingSet) mentalese.RelationSet {

	answer := mentalese.RelationSet{}
	transformer := NewRelationTransformer(answerer.matcher, answerer.log)

	allSolutions := answerer.findSolutions(goal)

	if len(allSolutions) == 0 {

		answerer.log.AddError("There are no solutions for this problem")

	} else {

		for i, solution := range allSolutions {

			answerer.log.AddProduction("Solution", solution.Condition.String())

			// apply transformation, if available
			transformedGoal := transformer.Replace(solution.Transformations, goal)

			// resultBindings: map goal variables to answers
			resultBindings := answerer.solver.SolveRelationSet(transformedGoal, bindings)

			// no results? try the next solution (if there is one)
			if resultBindings.IsEmpty() {

				// stack trace
				answerer.log.AddProduction("Stack trace", answerer.solver.callStack.String())

				if i < len(allSolutions) - 1 {
					continue
				}
			}

			group := EntityReferenceGroup{}
			for _, id := range resultBindings.GetIds(solution.Result.TermValue) {
				group = append(group, CreateEntityReference(id.TermValue, id.TermSort))
			}
			answerer.solver.dialogContext.AnaphoraQueue.AddReferenceGroup(group)

			// find a handler
			conditionedBindings := resultBindings
			var resultHandler *mentalese.ResultHandler
			for _, response := range solution.Responses {
				if !response.Condition.IsEmpty() {
					conditionBindings := answerer.solver.SolveRelationSet(response.Condition, resultBindings)
					if conditionBindings.IsEmpty() {
						continue
					} else {
						conditionedBindings = conditionBindings
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
			var solutionBindings = conditionedBindings

			// extend solution bindings by executing the preparation
			if !resultHandler.Preparation.IsEmpty() {
				solutionBindings = answerer.solver.SolveRelationSet(resultHandler.Preparation, conditionedBindings)
			}

			// create answer relation sets by binding 'answer' to solutionBindings
			answer = answerer.build(resultHandler.Answer, solutionBindings)

			// stop after the first solution
			break
		}
	}

	return answer
}

// Returns the solutions whose condition matches the goal, and a set of bindings per solution
func (answerer Answerer) findSolutions(goal mentalese.RelationSet) []mentalese.Solution {

	var solutions []mentalese.Solution

	for _, aSolution := range answerer.solutions {

		unScopedGoal := goal.UnScope()

		bindings, found := answerer.matcher.MatchSequenceToSet(aSolution.Condition, unScopedGoal, mentalese.NewBinding())
		if found {

			for _, binding := range bindings.GetAll() {
				boundSolution := aSolution.BindSingle(binding)
				solutions = append(solutions, boundSolution)
			}
		}
	}

	return solutions
}

func (answerer Answerer) build(template mentalese.RelationSet, bindings mentalese.BindingSet) mentalese.RelationSet {

	newSet := mentalese.RelationSet{}

	if bindings.IsEmpty() {
		newSet = template
	} else {

		sets := template.BindRelationSetMultipleBindings(bindings)

		newSet = mentalese.RelationSet{}
		for _, set := range sets {
			newSet = newSet.Merge(set)
		}
	}

	return newSet
}
