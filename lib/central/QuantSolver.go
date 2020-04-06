package central

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
)

// This part of the problem solver produces a set of bindings given a quant and a binding
// A quant is a scoped quantification like this (this is nested quant)
//
//     quant(S1, [ isa(S1, parent) ], D1, [ isa(D1, every) ], [
//         quant(O1, [ isa(O1, child) ], D2, [ isa(D2, none) ], [
//             have_child(S1, O1)
//         ])
//     ])

func (solver ProblemSolver) SolveQuant(quant mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	groupedScopeBindings := []mentalese.Bindings{}
	scopeBindings := mentalese.Bindings{}

	rangeVariable := quant.Arguments[mentalese.QuantRangeVariableIndex].TermValue
	quantifier := quant.Arguments[mentalese.QuantQuantifierIndex]
	rangeSet := quant.Arguments[mentalese.QuantRangeIndex]
	count := 0

	// pick the number from the quantifier, if applicable
	if quantifier.TermValueRelationSet[0].Predicate == mentalese.PredicateNumber {
		numberRelation := quantifier.TermValueRelationSet[0]
		count, _ = strconv.Atoi(numberRelation.Arguments[1].TermValue)
	}

	// solve the range
	rangeBindings := []mentalese.Binding{}

	isTheQuantifier := quantifier.TermValueRelationSet[0].Predicate == mentalese.PredicateThe
	isAllQuantifier := quantifier.TermValueRelationSet[0].Predicate == mentalese.PredicateAll
	isSomeQuantifier := quantifier.TermValueRelationSet[0].Predicate == mentalese.PredicateSome

	isReference := len(rangeSet.TermValueRelationSet) > 0 && rangeSet.TermValueRelationSet[0].Predicate == mentalese.PredicateReference

	if isTheQuantifier || isReference {
		count = 1

		// try the anaphora queue first
		refFound := false
		for _, group := range *solver.dialogContext.AnaphoraQueue {

			ref := group[0]

			refBinding := binding.Merge(mentalese.Binding{ rangeVariable: mentalese.NewId(ref.Id, ref.EntityType)})
			rangeSet := quant.Arguments[mentalese.QuantRangeIndex].TermValueRelationSet
			//  empty range set for "it"
			if len(rangeSet) == 0 {
				refFound = true
				rangeBindings = mentalese.Bindings{ refBinding }
				break
			}
			if !solver.quickAcceptabilityCheck(rangeVariable, ref.EntityType, rangeSet) {
				continue
			}
			testRangeBindings := solver.SolveRelationSet(rangeSet, mentalese.Bindings{refBinding})
			if len(testRangeBindings) == 1 {
				refFound = true
				rangeBindings = testRangeBindings
				break
			}
		}

		if !refFound {
			rangeBindings = solver.SolveRelationSet(quant.Arguments[mentalese.QuantRangeIndex].TermValueRelationSet, mentalese.Bindings{binding})
		}

		if len(rangeBindings) != 1 {

			rangeIndex, found := solver.rangeIndexClarification(rangeBindings, rangeVariable)
			if found {
				rangeBindings = rangeBindings[rangeIndex:rangeIndex + 1]
			} else {
				return mentalese.Bindings{}
			}
		}

	} else {
		rangeBindings = solver.SolveRelationSet(quant.Arguments[mentalese.QuantRangeIndex].TermValueRelationSet, mentalese.Bindings{binding})
	}

	// this now works for EVERY and for NUMBER at the moment; check the quantifier from the quant!

	// for each entity in the range, solve the scoped relations
	index := 0
	for _, rangeBinding := range rangeBindings {

		singleScopeBindings := solver.SolveRelationSet(quant.Arguments[mentalese.QuantScopeIndex].TermValueRelationSet, mentalese.Bindings{rangeBinding})

		if len(singleScopeBindings) > 0 {
			groupedScopeBindings = append(groupedScopeBindings, singleScopeBindings)
			scopeBindings = append(scopeBindings, singleScopeBindings...)
		}

		value, found := rangeBinding[rangeVariable]
		if found && value.IsId() {
			group := EntityReferenceGroup{ CreateEntityReference(value.TermValue, value.TermEntityType) }
			solver.dialogContext.AnaphoraQueue.AddReferenceGroup(group)
		}

		index++
		if count > 0 && index == count {
			break
		}
	}

	if isTheQuantifier {
		// THE
		if len(groupedScopeBindings) == 1 {
			return scopeBindings
		} else {
			return mentalese.Bindings{}
		}
	} else if isAllQuantifier {
		// ALL
		if len(groupedScopeBindings) == len(rangeBindings) {
			return scopeBindings
		} else {
			return mentalese.Bindings{}
		}
	} else if isSomeQuantifier {
		// SOME
		if len(groupedScopeBindings) > 0 {
			return scopeBindings
		} else {
			return mentalese.Bindings{}
		}
	} else if count > 0 {
		// A NUMBER
		if len(groupedScopeBindings) == count {
			return scopeBindings
		} else {
			return mentalese.Bindings{}
		}
	} else {
		// NO SIMPLE QUANTIFIER
		// todo
		return scopeBindings
	}
}

func (solver ProblemSolver) quickAcceptabilityCheck(variable string, entityType string, relations mentalese.RelationSet) bool {

	accepted := false

	for _, relation := range relations {
		for i, argument := range relation.Arguments {
			if argument.IsVariable() && argument.TermValue == variable {
				argumentEntityType := solver.predicates.GetEntityType(relation.Predicate, i)
				if  argumentEntityType == "" || argumentEntityType == entityType {
					accepted = true
					break
				}
			}
		}
	}

	return accepted
}

// ask the user which of the specified entities he/she means
func (solver ProblemSolver) rangeIndexClarification(rangeBindings mentalese.Bindings, rangeVariable string) (int, bool) {

	options := common.NewOptions()

	answer, found := solver.dialogContext.GetAnswerToOpenQuestion()

	if found {

		index, _ := strconv.Atoi(answer)
		solver.dialogContext.RemoveAnswerToOpenQuestion()

		return index, true

	} else {

		solver.log.SetClarificationRequest("I don't understand which one you mean", options)
		return 0, false
	}
}