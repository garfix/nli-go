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

// Returns the solutions whose condition matches the goal, and a set of bindings per solution
func (answerer Answerer) FindSolutions(goal mentalese.RelationSet) []mentalese.Solution {

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

func (answerer Answerer) Build(template mentalese.RelationSet, bindings mentalese.BindingSet) mentalese.RelationSet {

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
