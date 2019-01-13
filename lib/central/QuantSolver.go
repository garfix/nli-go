package central

import (
	"nli-go/lib/mentalese"
)

// This part of the problem solver produces a set of bindings given a quant and a binding
// A quant is a scoped quantification like this (this is nested quant)
//
//     quant(S1, [ isa(S1, parent) ], D1, [ isa(D1, every) ], [
//         quant(O1, [ isa(O1, child) ], D2, [ isa(D2, none) ], [
//             have_child(S1, O1)
//         ])
//     ])

func (solver ProblemSolver) SolveQuant(quant mentalese.Relation, nameStore *ResolvedNameStore, binding mentalese.Binding) []mentalese.Binding {
	// solve the range
	rangeBindings := solver.SolveRelationSet(quant.Arguments[mentalese.Quantification_RangeIndex].TermValueRelationSet, nameStore, []mentalese.Binding{binding})

	combinedScopeBindings := [][]mentalese.Binding{}

	// for each entity in the range, solve the scoped relations
	for _, rangeBinding := range rangeBindings {

		scopeBindings := solver.SolveRelationSet(quant.Arguments[mentalese.Quantification_ScopeIndex].TermValueRelationSet, nameStore, []mentalese.Binding{rangeBinding})

		if len(scopeBindings) > 0 {
			combinedScopeBindings = append(combinedScopeBindings, scopeBindings)
		}
	}

	// todo: this only works for EVERY at the moment; check the quantifier from the quant!

	if len(combinedScopeBindings) == len(rangeBindings) {
		return rangeBindings
	} else {
		return []mentalese.Binding{}
	}
}
