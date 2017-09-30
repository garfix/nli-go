package central

import (
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
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

func NewAnswerer(matcher *mentalese.RelationMatcher, log *common.SystemLog) *Answerer {

	builder := NewRelationSetBuilder()
	builder.addGenerator(NewSystemGenerator())

	return &Answerer{
		solutions: []mentalese.Solution{},
		matcher:   matcher,
		solver:    NewProblemSolver(matcher, log),
		builder:   builder,
		log:       log,
	}
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

func (answerer *Answerer) AddNestedStructureBase(base knowledge.NestedStructureBase) {
	answerer.solver.AddNestedStructureBase(base)
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
		solutionBindings := []mentalese.Binding{}
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

		unscopedGoal := answerer.Unscope(goal)

		bindings, found = answerer.matcher.MatchSequenceToSet(aSolution.Condition, unscopedGoal, mentalese.Binding{})
		if found {
			solution = aSolution
			break
		}
	}

	answerer.log.EndDebug("findSolution", solution, bindings, found)

	return solution, bindings, found
}

func (Answerer Answerer) Unscope(relations mentalese.RelationSet) mentalese.RelationSet {

	unscoped := mentalese.RelationSet{}

	for _, relation := range relations {

		copy := relation.Copy()

		if relation.Predicate == mentalese.Predicate_Quant {
			// unscope the relation sets
			for i, argument := range relation.Arguments {
				if argument.IsRelationSet() {

					scopedSet := copy.Arguments[i].TermValueRelationSet
					copy.Arguments[i].TermValueRelationSet = mentalese.RelationSet{}

					// recurse into the scope
					unscoped = append(unscoped, Answerer.Unscope(scopedSet)...)
				}
			}
		}

		unscoped = append(unscoped, copy)
	}

	return unscoped
}
