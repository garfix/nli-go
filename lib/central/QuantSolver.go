package central

import (
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

func (solver ProblemSolver) SolveQuant(quant mentalese.Relation, nameStore *mentalese.ResolvedNameStore, binding mentalese.Binding) []mentalese.Binding {
	// solve the range
	rangeBindings := solver.SolveRelationSet(quant.Arguments[mentalese.Quantification_RangeIndex].TermValueRelationSet, nameStore, []mentalese.Binding{binding})

	combinedScopeBindings := [][]mentalese.Binding{}

	quantifier := quant.Arguments[mentalese.Quantification_QuantifierIndex]
	count := 0

	// pick the number from the quantifier, if applicable
	if quantifier.TermValueRelationSet[0].Predicate == mentalese.PredicateNumber {
		numberRelation := quantifier.TermValueRelationSet[0]
		count, _ = strconv.Atoi(numberRelation.Arguments[1].TermValue)
	}

	// this now works for EVERY and for NUMBER at the moment; check the quantifier from the quant!

	// for each entity in the range, solve the scoped relations
	index := 0
	for _, rangeBinding := range rangeBindings {

		scopeBindings := solver.SolveRelationSet(quant.Arguments[mentalese.Quantification_ScopeIndex].TermValueRelationSet, nameStore, []mentalese.Binding{rangeBinding})

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
			return []mentalese.Binding{}
		}
	} else {
		// EVERY
		if len(combinedScopeBindings) == len(rangeBindings) {
			return rangeBindings
		} else {
			return []mentalese.Binding{}
		}
	}
}
