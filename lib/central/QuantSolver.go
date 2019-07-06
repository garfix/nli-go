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

func (solver ProblemSolver) SolveQuant(quant mentalese.Relation, keyCabinet *mentalese.KeyCabinet, binding mentalese.Binding) mentalese.Bindings {
	// solve the range
	rangeBindings := solver.SolveRelationSet(quant.Arguments[mentalese.QuantRangeIndex].TermValueRelationSet, keyCabinet, mentalese.Bindings{binding})

	combinedScopeBindings := []mentalese.Bindings{}

	rangeVariable := quant.Arguments[mentalese.QuantRangeVariableIndex].TermValue
	quantifier := quant.Arguments[mentalese.QuantQuantifierIndex]
	count := 0

	// pick the number from the quantifier, if applicable
	if quantifier.TermValueRelationSet[0].Predicate == mentalese.PredicateNumber {
		numberRelation := quantifier.TermValueRelationSet[0]
		count, _ = strconv.Atoi(numberRelation.Arguments[1].TermValue)
	}

	if quantifier.TermValueRelationSet[0].Predicate == mentalese.PredicateIsa &&
		quantifier.TermValueRelationSet[0].Arguments[1].TermType == mentalese.TermPredicateAtom &&
		quantifier.TermValueRelationSet[0].Arguments[1].TermValue == mentalese.AtomThe {
		count = 1

		if len(rangeBindings) != 1 {

			rangeIndex, found := solver.rangeIndexClarification(rangeBindings, rangeVariable)
			if found {
				rangeBindings = rangeBindings[rangeIndex:rangeIndex + 1]
			} else {
				return mentalese.Bindings{}
			}
		}
	}

	// this now works for EVERY and for NUMBER at the moment; check the quantifier from the quant!

	// for each entity in the range, solve the scoped relations
	index := 0
	for _, rangeBinding := range rangeBindings {

		scopeBindings := solver.SolveRelationSet(quant.Arguments[mentalese.QuantScopeIndex].TermValueRelationSet, keyCabinet, mentalese.Bindings{rangeBinding})

		if len(scopeBindings) > 0 {
			combinedScopeBindings = append(combinedScopeBindings, scopeBindings)
		}

		index++
		if count > 0 && index == count {
			break
		}
	}

	if count > 0 {
		// NUMBER
		if len(combinedScopeBindings) == count {
			return rangeBindings
		} else {
			return mentalese.Bindings{}
		}
	} else {
		// EVERY
		if len(combinedScopeBindings) == len(rangeBindings) {
			return rangeBindings
		} else {
			return mentalese.Bindings{}
		}
	}
}

// ask the user which of the specified entities he/she means
func (solver ProblemSolver) rangeIndexClarification(rangeBindings mentalese.Bindings, rangeVariable string) (int, bool) {

	options := common.NewOptions()

	//for i, rangeBinding := range rangeBindings {
	//	val := rangeBinding[rangeVariable].TermValue
	//	options.AddOption(strconv.Itoa(i), val)
	//}

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