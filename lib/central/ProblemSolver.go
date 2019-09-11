package central

import (
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"strings"
)

// The problem solver takes a relation set and a set of bindings
// and returns a set of new bindings
// It uses knowledge bases to find these bindings
//
// Structures used:
// RelationSet: describes the query
// SolutionRoute: describes a single path through a series of knowledge base calls
// RelationGroup: a single step in a solution route; it is the input for a specified knowledge base, having a calculation cost
type ProblemSolver struct {
	allKnowledgeBases   []knowledge.KnowledgeBase
	factBases           []knowledge.FactBase
	ruleBases           []knowledge.RuleBase
	aggregateBases      []knowledge.AggregateBase
	nestedStructureBase []knowledge.NestedStructureBase
	matcher             *mentalese.RelationMatcher
	optimizer           Optimizer
	modifier            *FactBaseModifier
	dialogContext 		*DialogContext
	log                 *common.SystemLog
	SolveDepth int
}

func NewProblemSolver(matcher *mentalese.RelationMatcher, dialogContext *DialogContext, log *common.SystemLog) *ProblemSolver {
	return &ProblemSolver{
		factBases:      []knowledge.FactBase{},
		ruleBases:      []knowledge.RuleBase{},
		aggregateBases: []knowledge.AggregateBase{},
		matcher:        matcher,
		optimizer:      NewOptimizer(matcher),
		modifier:       NewFactBaseModifier(log),
		dialogContext:  dialogContext,
		log:            log,
		SolveDepth:     0,
	}
}

func (solver *ProblemSolver) AddFactBase(factBase knowledge.FactBase) {
	solver.factBases = append(solver.factBases, factBase)
	solver.allKnowledgeBases = append(solver.allKnowledgeBases, factBase)
}

func (solver *ProblemSolver) AddFunctionBase(functionBase knowledge.FunctionBase) {
	solver.allKnowledgeBases = append(solver.allKnowledgeBases, functionBase)
}

func (solver *ProblemSolver) AddRuleBase(ruleBase knowledge.RuleBase) {
	solver.ruleBases = append(solver.ruleBases, ruleBase)
	solver.allKnowledgeBases = append(solver.allKnowledgeBases, ruleBase)
}

func (solver *ProblemSolver) AddMultipleBindingsBase(source knowledge.AggregateBase) {
	solver.aggregateBases = append(solver.aggregateBases, source)
	solver.allKnowledgeBases = append(solver.allKnowledgeBases, source)
}

func (solver *ProblemSolver) AddNestedStructureBase(base knowledge.NestedStructureBase) {
	solver.nestedStructureBase = append(solver.nestedStructureBase, base)
	solver.allKnowledgeBases = append(solver.allKnowledgeBases, base)
}

// set e.g. [ father(X, Y) father(Y, Z) ]
// bindings e.g. [{X: john, Z: jack} {}]
// return e.g. [
//  { X: john, Z: jack, Y: billy }
//  { X: john, Z: jack, Y: bob }
// ]
func (solver ProblemSolver) SolveRelationSet(set mentalese.RelationSet, keyCabinet *mentalese.KeyCabinet, bindings mentalese.Bindings) mentalese.Bindings {

	solver.log.StartDebug("SolveRelationSet", set, bindings)

	if keyCabinet == nil {
		keyCabinet = &mentalese.KeyCabinet{}
	}

	head := strings.Repeat("  ", solver.SolveDepth)

	solver.log.AddProduction(head + "Solve Set", set.String() + " " + keyCabinet.String() + " " + bindings.String())

	newBindings := mentalese.Bindings{}

	// remove duplicates because they cause unnecessary work and the optimizer can't deal with them
	set = set.RemoveDuplicates()

	// sort the relations to reduce the number of tuples retrieved from the fact bases
	solutionRoutes, remainingRelations, ok := solver.optimizer.CreateSolutionRoutes(set, solver.allKnowledgeBases, keyCabinet)

	solver.log.AddProduction(head + "Solution Routes", solutionRoutes.String())

	if !ok {

		solver.log.AddError("Cannot find these relations in any knowledge base: " + remainingRelations.String())

	} else {

		for _, solutionRoute := range solutionRoutes {
			newBindings = append(newBindings, solver.solveSingleSolutionRouteMultipleBindings(solutionRoute, keyCabinet, bindings)...)
		}
	}

	// remove duplicates because they cause unnecessary work and they cause problems for the generator
	newBindings = mentalese.UniqueBindings(newBindings)

	solver.log.AddProduction(head + "Result", newBindings.String())

	solver.log.EndDebug("SolveRelationSet", newBindings)

	return newBindings
}

func (solver ProblemSolver) solveSingleSolutionRouteMultipleBindings(solutionRoute knowledge.SolutionRoute, keyCabinet *mentalese.KeyCabinet, bindings mentalese.Bindings) mentalese.Bindings {

	newBindings := bindings

	solver.SolveDepth++

	head := strings.Repeat("  ", solver.SolveDepth)

	for _, relationGroup := range solutionRoute {

		solver.log.AddProduction(head + "Solve RelationGroup", relationGroup.String() + " " + keyCabinet.String() + " " + bindings.String())

		newBindings = solver.solveSingleRelationGroupMultipleBindings(relationGroup, keyCabinet, newBindings)

		solver.log.AddProduction(head + "Result", newBindings.String())

		if len(newBindings) == 0 {
			break
		}
	}

	solver.SolveDepth--

	return newBindings
}

func (solver ProblemSolver) findKnowledgeBaseByName(name string) knowledge.KnowledgeBase {
	for _, knowledgeBase := range solver.allKnowledgeBases {
		if knowledgeBase.GetName() == name {
			return knowledgeBase
		}
	}

	return nil
}

func (solver ProblemSolver) solveSingleRelationGroupMultipleBindings(relationGroup knowledge.RelationGroup, keyCabinet *mentalese.KeyCabinet, bindings mentalese.Bindings) mentalese.Bindings {

	solver.log.StartDebug("solveSingleRelationGroupMultipleBindings", relationGroup, bindings)

	newBindings := mentalese.Bindings{}

	knowledgeBase := solver.findKnowledgeBaseByName(relationGroup.KnowledgeBaseName)
	aggregateBase, isAggregateBase := knowledgeBase.(knowledge.AggregateBase)

	if isAggregateBase {

		mbBindings, ok := aggregateBase.Bind(relationGroup.Relations[0], bindings)

		if ok {
			newBindings = append(newBindings, mbBindings...)
		}

	} else {

		for _, binding := range bindings {
			groupBindings := solver.solveSingleRelationGroupSingleBinding(relationGroup, keyCabinet, binding)
			newBindings = append(newBindings, groupBindings...)
		}
	}

	solver.log.EndDebug("solveSingleRelationGroupMultipleBindings", newBindings)

	return newBindings
}

func (solver ProblemSolver) solveSingleRelationGroupSingleBinding(relationGroup knowledge.RelationGroup, keyCabinet *mentalese.KeyCabinet, binding mentalese.Binding) mentalese.Bindings {

	solver.log.StartDebug("solveSingleRelationGroupSingleBinding", relationGroup, binding)

	knowledgeBase := solver.findKnowledgeBaseByName(relationGroup.KnowledgeBaseName)
	factBase, isFactBase := knowledgeBase.(knowledge.FactBase)
	ruleBase, isRuleBase := knowledgeBase.(knowledge.RuleBase)
	functionBase, isFunctionBase := knowledgeBase.(knowledge.FunctionBase)
	_, isNestedStructureBase := knowledgeBase.(knowledge.NestedStructureBase)

	boundRelations := relationGroup.Relations.BindSingle(binding)

	var newBindings mentalese.Bindings

	if isFactBase {

		if len(boundRelations) == 1 && boundRelations[0].Predicate == mentalese.PredicateAssert {

			boundRelations = keyCabinet.BindToRelationSet(boundRelations, factBase.GetName())

			solver.modifier.Assert(boundRelations[0].Arguments[0].TermValueRelationSet, factBase, keyCabinet)
			newBindings = append(newBindings, binding)

		} else if len(boundRelations) == 1 && boundRelations[0].Predicate == mentalese.PredicateRetract {

			boundRelations = keyCabinet.BindToRelationSet(boundRelations, factBase.GetName())

			solver.modifier.Retract(boundRelations[0].Arguments[0].TermValueRelationSet, factBase, keyCabinet)
			newBindings = append(newBindings, binding)

		} else {

			boundRelations = keyCabinet.BindToRelationSet(boundRelations, factBase.GetName())

			sourceBindings := solver.FindFacts(factBase, boundRelations)

			for _, sourceBinding := range sourceBindings {

				combinedBinding := binding.Merge(sourceBinding)
				newBindings = append(newBindings, combinedBinding)
			}
		}

	} else if isFunctionBase {

		var relation = relationGroup.Relations[0]
		resultBinding, functionFound := functionBase.Execute(relation, binding)
		if functionFound {
			newBindings = append(newBindings, resultBinding)
		}

	} else if isRuleBase {

		newBindings = append(newBindings, solver.SolveSingleRelationSingleBindingSingleRuleBase(boundRelations[0], keyCabinet, binding, ruleBase)...)

	} else if isNestedStructureBase {

		newBindings = solver.SolveChildStructures(relationGroup.Relations[0], keyCabinet, binding)

	}

	solver.log.EndDebug("solveSingleRelationGroupSingleBinding", newBindings)

	return newBindings
}

func (solver ProblemSolver) SolveChildStructures(goal mentalese.Relation, keyCabinet *mentalese.KeyCabinet, binding mentalese.Binding) mentalese.Bindings {

	solver.log.StartDebug("NestedStructureBase BindChildStructures", goal, binding)

	var newBindings mentalese.Bindings

	if goal.Predicate == mentalese.PredicateQuant {

		newBindings = solver.SolveQuant(goal, keyCabinet, binding)

	} else if goal.Predicate == mentalese.PredicateSequence {

		newBindings = solver.SolveSeq(goal, keyCabinet, binding)

	}

	solver.log.EndDebug("NestedStructureBase BindChildStructures", newBindings)

	return newBindings
}


func (solver ProblemSolver) FindFacts(factBase knowledge.FactBase, goal mentalese.RelationSet) mentalese.Bindings {

	solver.log.StartDebug("FindFacts", goal)

	subgoalBindings := mentalese.Bindings{}

	for _, ds2db := range factBase.GetMappings() {

		// gender(14, G), gender(A, male) => internalBinding: A = 14
		internalBindingsX, match1 := solver.matcher.MatchSequenceToSet(ds2db.Pattern, goal, mentalese.Binding{})
		if match1 {

			internalBinding := internalBindingsX[0]

			// gender(14, G), gender(A, male) => externalBinding: G = male
			externalBindings, match2 := solver.matcher.MatchSequenceToSet(goal, ds2db.Pattern, mentalese.Binding{})
			if match2 {

				externalBinding := externalBindings[0]

				// create a version of the conditions with bound variables
				boundConditions := ds2db.Replacement.BindSingle(internalBinding)

				// match this bound version to the database
				internalBindings, match3 := solver.SolveMultipleRelationSingleFactBase(ds2db.Replacement, boundConditions, factBase)

				if match3 {
					for _, binding := range internalBindings {
						subgoalBindings = append(subgoalBindings, externalBinding.Intersection(binding))
					}
				}
			}
		}
	}

	solver.log.EndDebug("FindFacts", subgoalBindings)

	return subgoalBindings
}

func (solver ProblemSolver) SolveMultipleRelationSingleFactBase(unboundSequence []mentalese.Relation, boundSequence []mentalese.Relation, factBase knowledge.FactBase) (mentalese.Bindings, bool) {

	solver.log.StartDebug("SolveMultipleRelationSingleFactBase", boundSequence)

	// bindings using database level variables
	sequenceBindings := mentalese.Bindings{}
	match := true

	for i, relation := range boundSequence {

		relationBindings := mentalese.Bindings{}

		aggregateFunctionFound := false
		for _, aggregateBase := range solver.aggregateBases {
			newRelationBindings, ok := aggregateBase.Bind(unboundSequence[i], sequenceBindings)
			if ok {
				relationBindings = newRelationBindings
				aggregateFunctionFound = true
				break
			}
		}

		if !aggregateFunctionFound {

			if len(sequenceBindings) == 0 {

				resultBindings := factBase.MatchRelationToDatabase(relation)
				relationBindings = resultBindings

			} else {

				//// go through the bindings resulting from previous relation
				for _, binding := range sequenceBindings {

					boundRelation := relation.BindSingleRelationSingleBinding(binding)
					resultBindings := factBase.MatchRelationToDatabase(boundRelation)

					// found bindings must be extended with the bindings already present
					for _, resultBinding := range resultBindings {
						newRelationBinding := binding.Merge(resultBinding)
						relationBindings = append(relationBindings, newRelationBinding)
					}
				}
			}
		}

		sequenceBindings = relationBindings

		if len(sequenceBindings) == 0 {
			match = false
			break
		}
	}

	solver.log.EndDebug("SolveMultipleRelationSingleFactBase", sequenceBindings, match)

	return sequenceBindings, match
}

// goalRelation e.g. father('jack', Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver ProblemSolver) SolveSingleRelationSingleBindingSingleRuleBase(goalRelation mentalese.Relation, keyCabinet *mentalese.KeyCabinet, binding mentalese.Binding, ruleBase knowledge.RuleBase) mentalese.Bindings {

	solver.log.StartDebug("SolveSingleRelationSingleBindingSingleRuleBase", goalRelation, binding)

	for _, val := range binding {
		if val.TermType == mentalese.TermVariable {
			//panic("Variable bound to variable")
		}
	}

	inputVariables := goalRelation.GetVariableNames()

	goalBindings := mentalese.Bindings{}

	// match rules from the rule base to the goalRelation
	boundRelation := goalRelation.BindSingleRelationSingleBinding(binding)
	sourceSubgoalSets, sourceBindings := ruleBase.Bind(boundRelation)

	for i, sourceSubgoalSet := range sourceSubgoalSets {

		// sourceBinding: from subgoal variable to goal argument
		sourceBinding := sourceBindings[i]

		// rewrite the variables of subgoal set to those of goalRelation
		importedSubgoalSet := sourceSubgoalSet.ImportBinding(sourceBinding)

		subgoalResultBindings := mentalese.Bindings{binding}

		for _, subGoal := range importedSubgoalSet {

			subgoalResultBindings = solver.SolveRelationSet([]mentalese.Relation{subGoal}, keyCabinet, subgoalResultBindings)
		}

		for _, subgoalResultBinding := range subgoalResultBindings {

			// filter out the input variables
			filteredBinding := subgoalResultBinding.FilterVariablesByName(inputVariables)

			// make sure all variables of the original binding are present
			goalBinding := binding.Merge(filteredBinding)

			goalBindings = append(goalBindings, goalBinding)
		}
	}

	solver.log.EndDebug("SolveSingleRelationSingleBindingSingleRuleBase", goalBindings)

	return goalBindings
}
