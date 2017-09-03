package central

import (
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
)

type ProblemSolver struct {
	allKnowledgeBases     []knowledge.KnowledgeBase
	factBases             []knowledge.FactBase
	ruleBases             []knowledge.RuleBase
	multipleBindingsBases []knowledge.MultipleBindingsBase
	matcher               *mentalese.RelationMatcher
	optimizer			  Optimizer
	log                   *common.SystemLog
}

func NewProblemSolver(matcher *mentalese.RelationMatcher, log *common.SystemLog) *ProblemSolver {
	return &ProblemSolver{
		factBases:             []knowledge.FactBase{},
		ruleBases:             []knowledge.RuleBase{},
		multipleBindingsBases: []knowledge.MultipleBindingsBase{},
		matcher:               matcher,
		optimizer: 			   Optimizer{},
		log:                   log,
	}
}

func (solver *ProblemSolver) AddFactBase(factBase knowledge.FactBase) {
	solver.factBases = append(solver.factBases, factBase)
	solver.allKnowledgeBases = append(solver.allKnowledgeBases, factBase)
}

func (solver *ProblemSolver) AddRuleBase(ruleBase knowledge.RuleBase) {
	solver.ruleBases = append(solver.ruleBases, ruleBase)
	solver.allKnowledgeBases = append(solver.allKnowledgeBases, ruleBase)
}

func (solver *ProblemSolver) AddMultipleBindingsBase(source knowledge.MultipleBindingsBase) {
	solver.multipleBindingsBases = append(solver.multipleBindingsBases, source)
	solver.allKnowledgeBases = append(solver.allKnowledgeBases, source)
}

// Checks if all relations in set are handled by some knowledge store
func (solver *ProblemSolver) CheckMappings(set mentalese.RelationSet) (bool, mentalese.Relation) {

	relation := mentalese.Relation{}

	for _, relation = range set {

		found := false

		if relation.Predicate == mentalese.Predicate_Quant {

			quant := relation
			rangeSet := quant.Arguments[mentalese.Quantification_RangeIndex].TermValueRelationSet
			scopeSet := quant.Arguments[mentalese.Quantification_ScopeIndex].TermValueRelationSet

			found, unfoundSubRelation := solver.CheckMappings(rangeSet)
			if !found {
				return false, unfoundSubRelation
			}

			found, unfoundSubRelation = solver.CheckMappings(scopeSet)
			if !found {
				return false, unfoundSubRelation
			}

			// the 'quant' relation itself needs no mapping
			continue
		}

		for _, kb := range solver.allKnowledgeBases {
			if kb.Knows(relation) {
				found = true
			}
		}

		if !found {
			return false, relation
		}

	}

	return true, relation
}

// goals e.g. [ father(X, Y) father(Y, Z) ]
// return e.g. [
//  [ father('john', 'jack') father('jack', 'joe') ]
//  [ father('bob', 'jonathan') father('jonathan', 'bill') ]
// ]
func (solver ProblemSolver) Solve(goals []mentalese.Relation) []mentalese.RelationSet {

// NOTE: this function is only used by a test

	solver.log.StartDebug("Solve")
	bindings := solver.SolveRelationSet(goals, []mentalese.Binding{})
	solutions := solver.matcher.BindRelationSetMultipleBindings(goals, bindings)

	solver.log.EndDebug("Solve", solutions)
	return solutions
}

// set e.g. [ father(X, Y) father(Y, Z) ]
// bindings e.g. [{X: john, Z: jack} {}]
// return e.g. [
//  { X: john, Z: jack, Y: billy }
//  { X: john, Z: jack, Y: bob }
// ]
func (solver ProblemSolver) SolveRelationSet(set mentalese.RelationSet, bindings []mentalese.Binding) []mentalese.Binding {

	solver.log.StartDebug("SolveRelationSet", set, bindings)

	// sort the relations to reduce the number of tuples retrieved from the fact bases
	sortedRelations := solver.optimizer.Optimize(set, solver.factBases)

	for _, relation := range sortedRelations {
		bindings = solver.SolveSingleRelationMultipleBindings(relation, bindings)

		if len(bindings) == 0 {
			break
		}
	}

	solver.log.EndDebug("SolveRelationSet", bindings)

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

	solver.log.StartDebug("SolveSingleRelationMultipleBindings", goalRelation, bindings)

	newBindings := []mentalese.Binding{}
	multiFound := false

	for _, multipleBindingsBase := range solver.multipleBindingsBases {
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

	solver.log.EndDebug("SolveSingleRelationMultipleBindings", newBindings)

	return newBindings
}

// goalRelation e.g. father(Y, Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) SolveSingleRelationSingleBinding(goalRelation mentalese.Relation, binding mentalese.Binding) []mentalese.Binding {

	solver.log.StartDebug("SolveSingleRelationSingleBinding", goalRelation, binding)

	newBindings := []mentalese.Binding{}

	// scoped quantification
	if goalRelation.Predicate == mentalese.Predicate_Quant {
		newBindings = append(newBindings, solver.SolveQuant(goalRelation, binding)...)
	} else {
		// go through all fact bases
		for _, factBase := range solver.factBases {
			newBindings = append(newBindings, solver.SolveSingleRelationSingleBindingSingleFactBase(goalRelation, binding, factBase)...)
		}

		// go through all rule bases
		for _, ruleBase := range solver.ruleBases {
			newBindings = append(newBindings, solver.SolveSingleRelationSingleBindingSingleRuleBase(goalRelation, binding, ruleBase)...)
		}
	}

	solver.log.EndDebug("SolveSingleRelationSingleBinding", newBindings)

	return newBindings
}

// boundRelation e.g. father('jack', Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) SolveSingleRelationSingleBindingSingleFactBase(goalRelation mentalese.Relation, binding mentalese.Binding, factBase knowledge.FactBase) []mentalese.Binding {

	solver.log.StartDebug("SolveSingleRelationSingleBindingSingleFactBase", goalRelation, binding)

	//for key, val := range binding {
	//	if val.TermType == mentalese.Term_variable {
	//		panic("Variable bound to variable " + key);
	//	}
	//}

	boundRelation := solver.matcher.BindSingleRelationSingleBinding(goalRelation, binding)

	newBindings := []mentalese.Binding{}

	// boundRelation e.g. father(X, 'john')
	// sourceBindings e.g. {
	//    { X='Jack' },
	// }
	sourceBindings := solver.FindFacts(factBase, boundRelation)

	for _, sourceBinding := range sourceBindings {

		combinedBinding := binding.Merge(sourceBinding)
		newBindings = append(newBindings, combinedBinding)
	}

	solver.log.EndDebug("SolveSingleRelationSingleBindingSingleFactBase", newBindings)

	return newBindings
}

func (solver ProblemSolver) FindFacts(factBase knowledge.FactBase, goal mentalese.Relation) []mentalese.Binding {

	solver.log.StartDebug("FindFacts", goal)

	subgoalBindings := []mentalese.Binding{}

	for _, ds2db := range factBase.GetMappings() {

		// gender(14, G), gender(A, male) => externalBinding: G = male
		externalBinding, match := solver.matcher.MatchTwoRelations(goal, ds2db.DsSource, mentalese.Binding{})
		if match {

			// gender(14, G), gender(A, male) => internalBinding: A = 14
			internalBinding, _ := solver.matcher.MatchTwoRelations(ds2db.DsSource, goal, mentalese.Binding{})

			// create a version of the conditions with bound variables
			boundConditions := solver.matcher.BindRelationSetSingleBinding(ds2db.DbTarget, internalBinding)

			// match this bound version to the database
			internalBindings, match := factBase.Bind(boundConditions)

			if match {
				for _, binding := range internalBindings {
					subgoalBindings = append(subgoalBindings, externalBinding.Intersection(binding))
				}
			}
		}
	}

	solver.log.EndDebug("FindFacts", subgoalBindings)

	return subgoalBindings
}

// goalRelation e.g. father('jack', Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) SolveSingleRelationSingleBindingSingleRuleBase(goalRelation mentalese.Relation, binding mentalese.Binding, ruleBase knowledge.RuleBase) []mentalese.Binding {

	solver.log.StartDebug("SolveSingleRelationSingleBindingSingleRuleBase", goalRelation, binding)

	for _, val := range binding {
		if val.TermType == mentalese.Term_variable {
			panic("Variable bound to variable")
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

	solver.log.EndDebug("SolveSingleRelationSingleBindingSingleRuleBase", goalBindings)

	return goalBindings
}

// goal e.g. { father(X, Y), father(Y, Z)}
// bindings {X='john', Y='jack'}
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) SolveMultipleRelationsSingleBinding(goals []mentalese.Relation, binding mentalese.Binding) []mentalese.Binding {

	solver.log.StartDebug("SolveMultipleRelationsSingleBinding", goals, binding)

	bindings := []mentalese.Binding{binding}

	for _, goal := range goals {
		bindings = solver.SolveSingleRelationMultipleBindings(goal, bindings)
	}

	solver.log.EndDebug("SolveMultipleRelationsSingleBinding", bindings)

	return bindings
}
