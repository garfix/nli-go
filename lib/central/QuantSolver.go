package central

import (
    "nli-go/lib/mentalese"
    "nli-go/lib/common"
    "strconv"
)

// This part of the problem solver produces a set of bindings given a quant and a binding
// A quant is a scoped quantification like this (this is nested quant)
//
//     quant(S1, [ isa(S1, parent) ], D1, [ isa(D1, every) ], [
//         quant(O1, [ isa(O1, child) ], D2, [ isa(D2, 2) ], [
//             have_child(S1, O1)
//         ])
//     ])
func (solver ProblemSolver) SolveQuant(quant mentalese.Relation, binding mentalese.Binding) []mentalese.Binding {

    common.LogTree("SolveQuant", quant, binding)

    rangeVariable := quant.Arguments[mentalese.Quantification_RangeVariableIndex]
    rangeSet := quant.Arguments[mentalese.Quantification_RangeIndex].TermValueRelationSet
//    quantifierVariable := quant.Arguments[mentalese.Quantification_QuantifierVariableIndex]
    quantifierSet := quant.Arguments[mentalese.Quantification_QuantifierIndex].TermValueRelationSet
    scopeSet := quant.Arguments[mentalese.Quantification_ScopeIndex].TermValueRelationSet

    // bind the range to variable bindings
    rangeBindings := solver.SolveMultipleRelationsSingleBinding(rangeSet, mentalese.Binding{})

    // evaluate the scope for each of the variable bindings
    scopeBindings := []mentalese.Binding{}
    for _, rangeBinding := range rangeBindings {
        scopeBinding := binding.Merge(rangeBinding)
        scopeBindings = append(scopeBindings, solver.SolveMultipleRelationsSingleBinding(scopeSet, scopeBinding)...)
    }

    // validate with the quantifier
    quantBindings := scopeBindings
    if !solver.validate(quantifierSet, rangeVariable, rangeBindings, scopeBindings) {
        quantBindings = []mentalese.Binding{}
    }

    common.LogTree("SolveQuant", quantBindings)

    return quantBindings
}

// Checks whether the quantity of scope with respect to range is according to the quantifier
func (solver ProblemSolver) validate(quantifierSet mentalese.RelationSet, rangeVariable mentalese.Term, rangeBindings []mentalese.Binding, scopeBindings []mentalese.Binding) bool {

    common.LogTree("validate", quantifierSet, rangeVariable, rangeBindings, scopeBindings)

    ok := false
    rangeCount := mentalese.CountUniqueValues(rangeVariable.TermValue, rangeBindings)
    scopeCount := mentalese.CountUniqueValues(rangeVariable.TermValue, scopeBindings)

    ok = scopeCount >= 1

    if len(quantifierSet) == 1 {
        quantifier := quantifierSet[0]
        simpleQuantifier := quantifier.Arguments[1]

        if simpleQuantifier.TermType == mentalese.Term_number {

            quantity, _ := strconv.Atoi(simpleQuantifier.TermValue)
            ok = quantity == scopeCount

        } else if simpleQuantifier.TermType == mentalese.Term_predicateAtom {

            if simpleQuantifier.TermValue == "all" {
                ok = scopeCount == rangeCount
            }
        }
    }

    common.LogTree("validate", ok)

    return ok
}