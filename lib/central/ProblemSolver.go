package central

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/knowledge"
	"nli-go/lib/common"
)

type ProblemSolver struct {
	factBases []knowledge.FactBase
	ruleBases []knowledge.RuleBase
	multipleBindingsBases []knowledge.MultipleBindingsBase
	matcher *mentalese.RelationMatcher
}

func NewProblemSolver(matcher *mentalese.RelationMatcher) *ProblemSolver {
	return &ProblemSolver{
		factBases: []knowledge.FactBase{},
		ruleBases: []knowledge.RuleBase{},
		multipleBindingsBases: []knowledge.MultipleBindingsBase{},
		matcher: matcher,
	}
}

func (solver *ProblemSolver) AddFactBase(factBase knowledge.FactBase) {
	solver.factBases = append(solver.factBases, factBase)
}

func (solver *ProblemSolver) AddRuleBase(ruleBase knowledge.RuleBase) {
	solver.ruleBases = append(solver.ruleBases, ruleBase)
}

func (solver *ProblemSolver) AddMultipleBindingsBase(source knowledge.MultipleBindingsBase) {
	solver.multipleBindingsBases = append(solver.multipleBindingsBases, source)
}

// goals e.g. [ father(X, Y) father(Y, Z) ]
// return e.g. [
//  [ father('john', 'jack') father('jack', 'joe') ]
//  [ father('bob', 'jonathan') father('jonathan', 'bill') ]
// ]
func (solver ProblemSolver) Solve(goals []mentalese.Relation) []mentalese.RelationSet {

	common.LogTree("Solve")
	bindings := solver.SolveRelationSet(goals, []mentalese.Binding{})
	solutions := solver.matcher.BindRelationSetMultipleBindings(goals, bindings)

	common.LogTree("Solve", solutions)
	return solutions
}

// set e.g. [ father(X, Y) father(Y, Z) ]
// bindings e.g. [{X: john, Z: jack} {}]
// return e.g. [
//  { X: john, Z: jack, Y: billy }
//  { X: john, Z: jack, Y: bob }
// ]
func (solver ProblemSolver) SolveRelationSet(set mentalese.RelationSet, bindings []mentalese.Binding) []mentalese.Binding {

	common.LogTree("SolveRelationSet", set, bindings)

	for _, relation := range set {
		bindings = solver.SolveSingleRelationMultipleBindings(relation, bindings)

		if len(bindings) == 0 {
			break
		}
	}

	common.LogTree("SolveRelationSet", bindings)

	return bindings
}

// goal e.g. father(Y, Z)
// bindings e.g. {
//  { {X='john', Y='jack'} }
//  { {X='bob', Y='jonathan'} }
// }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) SolveSingleRelationMultipleBindings(goalRelation mentalese.Relation, bindings []mentalese.Binding) []mentalese.Binding {

	common.LogTree("SolveSingleRelationMultipleBindings", goalRelation, bindings)

	newBindings := []mentalese.Binding{}
	multiFound := false

	for _ , multipleBindingsBase := range solver.multipleBindingsBases {
		newBindings, multiFound = multipleBindingsBase.Bind(goalRelation, bindings)
		if multiFound {
			break
		}
	}

	if !multiFound {

		if len(bindings) == 0 {
			newBindings = solver.SolveSingleRelationSingleBinding(goalRelation, mentalese.Binding{})
		} else {
			for _, binding := range bindings {
				newBindings = append(newBindings, solver.SolveSingleRelationSingleBinding(goalRelation, binding)...)
			}
		}
	}

	common.LogTree("SolveSingleRelationMultipleBindings", newBindings)

	return newBindings
}

// goalRelation e.g. father(Y, Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) SolveSingleRelationSingleBinding(goalRelation mentalese.Relation, binding mentalese.Binding) []mentalese.Binding {

	common.LogTree("SolveSingleRelationSingleBinding", goalRelation, binding)

	newBindings := []mentalese.Binding{}

	// go through all fact bases
	for _, factBase := range solver.factBases {
		newBindings = append(newBindings, solver.SolveSingleRelationSingleBindingSingleFactBase(goalRelation, binding, factBase)...)
	}

	// go through all rule bases
	for _, ruleBase := range solver.ruleBases {
		newBindings = append(newBindings, solver.SolveSingleRelationSingleBindingSingleRuleBase(goalRelation, binding, ruleBase)...)
	}

	common.LogTree("SolveSingleRelationSingleBinding", newBindings)

	return newBindings
}

// boundRelation e.g. father('jack', Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) SolveSingleRelationSingleBindingSingleFactBase(goalRelation mentalese.Relation, binding mentalese.Binding, factBase knowledge.FactBase) []mentalese.Binding {

	common.LogTree("SolveSingleRelationSingleBindingSingleFactBase", goalRelation, binding)

	for key, val := range binding {
		if val.TermType == mentalese.Term_variable {
			panic("Variable bound to variable " + key);
		}
	}

	boundRelation := solver.matcher.BindSingleRelationSingleBinding(goalRelation, binding)

	newBindings := []mentalese.Binding{}

	// boundRelation e.g. father(X, 'john')
	// sourceBindings e.g. {
	//    { X='Jack' },
	// }
	sourceBindings := factBase.Bind(boundRelation)

	for _, sourceBinding := range sourceBindings {

		combinedBinding := binding.Merge(sourceBinding)
		newBindings = append(newBindings, combinedBinding)
	}

	common.LogTree("SolveSingleRelationSingleBindingSingleFactBase", newBindings)

	return newBindings
}

// goalRelation e.g. father('jack', Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) SolveSingleRelationSingleBindingSingleRuleBase(goalRelation mentalese.Relation, binding mentalese.Binding, ruleBase knowledge.RuleBase) []mentalese.Binding {

	common.LogTree("SolveSingleRelationSingleBindingSingleRuleBase", goalRelation, binding)

	for _, val := range binding {
		if val.TermType == mentalese.Term_variable {
			panic("Variable bound to variable");
		}
	}

	goalBindings := []mentalese.Binding{}

	// match rules from the rule base to the goalRelation
	boundRelation := solver.matcher.BindSingleRelationSingleBinding(goalRelation, binding)
	sourceSubgoalSets, sourceBindings := ruleBase.Bind(boundRelation)

	for i, sourceSubgoalSet := range sourceSubgoalSets {

		// sourceBinding: from subgoal variable to goal argument
		sourceBinding := sourceBindings[i]

		// subgoalBinding: from subgoal variable to goal constant
		subgoalBinding := sourceBinding.RemoveVariables()

		subgoalResultBindings := solver.SolveMultipleRelationsSingleBinding(sourceSubgoalSet, subgoalBinding)

		// subgoalResultBinding: from subgoal variables to constants (contains temporary variables)
		for _, subgoalResultBinding := range subgoalResultBindings {

			// sourceBinding e.g. { A:X, B:'yellow' }
			// after swap { X:A }
			// bind { X:A } with { A:'red', B:'yellow', C:'blue' } results in { X:'red' }
			convertedBinding := sourceBinding.Swap().Bind(subgoalResultBinding)

			// start extending the new binding with goalRelation variables as keys
			goalBinding := binding.Merge(convertedBinding)
			goalBindings = append(goalBindings, goalBinding)
		}
	}

	common.LogTree("SolveSingleRelationSingleBindingSingleRuleBase", goalBindings)

	return goalBindings
}

// goal e.g. { father(X, Y), father(Y, Z)}
// bindings {X='john', Y='jack'}
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) SolveMultipleRelationsSingleBinding(goals []mentalese.Relation, binding mentalese.Binding) []mentalese.Binding {

	common.LogTree("SolveMultipleRelationsSingleBinding", goals, binding)

	bindings := []mentalese.Binding{binding}

	for _, goal := range goals {
		bindings = solver.SolveSingleRelationMultipleBindings(goal, bindings)
	}

	common.LogTree("SolveMultipleRelationsSingleBinding", bindings)

	return bindings
}