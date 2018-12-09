package central

import (
	"nli-go/lib/mentalese"
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
func (solver ProblemSolver) SolveQuant(quant mentalese.Relation, nameStore *ResolvedNameStore, binding mentalese.Binding) []mentalese.Binding {

	solver.log.StartDebug("SolveQuant", quant, binding)

	rangeVariable := quant.Arguments[mentalese.Quantification_RangeVariableIndex]
	scopeSet := quant.Arguments[mentalese.Quantification_ScopeIndex].TermValueRelationSet

	// todo refactor this
	solver.quantLevel++

	// evaluate the scope
	scopeBindings := solver.SolveRelationSet(scopeSet, nameStore, []mentalese.Binding{binding})

	solver.quantLevel--

	if solver.quantLevel == 0 {
		// outermost quant: collect all bindings that quantify, if any
		scopeBindings = solver.validateNestedQuants(quant, nameStore, scopeBindings, []mentalese.Term{rangeVariable})
	}

	solver.log.EndDebug("SolveQuant", scopeBindings)

	return scopeBindings
}

// returns all bindings that quantify
func (solver ProblemSolver) validateNestedQuants(quant mentalese.Relation, nameStore *ResolvedNameStore, bindings []mentalese.Binding, rangeVariables[]mentalese.Term) []mentalese.Binding {

	resultBindings := bindings

	// validate lower levels
	scopeSet := quant.Arguments[mentalese.Quantification_ScopeIndex].TermValueRelationSet
	for _, relation := range scopeSet  {
		if relation.Predicate == mentalese.Predicate_Quant {

			rangeVariable := relation.Arguments[mentalese.Quantification_RangeVariableIndex]
			newRangeVariables := append(rangeVariables, rangeVariable)

			resultBindings = solver.validateNestedQuants(relation, nameStore, bindings, newRangeVariables)
			if len(resultBindings) == 0 {
				break
			}
		}
	}

	// validate current level
	if len(resultBindings) > 0 {
		resultBindings = solver.validateQuantification(quant, nameStore, resultBindings, rangeVariables)
	}

	return resultBindings
}

func (solver ProblemSolver) validateQuantification(quant mentalese.Relation, nameStore *ResolvedNameStore, bindings []mentalese.Binding, rangeVariables[]mentalese.Term) []mentalese.Binding {

	rangeVariable := quant.Arguments[mentalese.Quantification_RangeVariableIndex]
	rangeSet := quant.Arguments[mentalese.Quantification_RangeIndex].TermValueRelationSet
	quantifierSet := quant.Arguments[mentalese.Quantification_QuantifierIndex].TermValueRelationSet

	// lazy loading of range bindings
	var rangeBindings []mentalese.Binding
	rangeBindings = nil

	// group bindings by the combination of range values A2/B2/C1 = [bindings]
	groupedBindings := solver.groupBindings(quant, bindings, rangeVariables)

	var resultBindings []mentalese.Binding

	for _, bindingGroup := range groupedBindings {

		ok := false

		uniqueRangeValues := mentalese.CountUniqueValues(rangeVariable.TermValue, bindingGroup)

		if len(quantifierSet) == 1 {
			quantifier := quantifierSet[0]
			simpleQuantifier := quantifier.Arguments[1]

			if simpleQuantifier.TermType == mentalese.Term_number {

				// todo

				quantity, _ := strconv.Atoi(simpleQuantifier.TermValue)
				ok = quantity == uniqueRangeValues

			} else if simpleQuantifier.TermType == mentalese.Term_predicateAtom {

				if simpleQuantifier.TermValue == "all" {

					// load range bindings once
					if rangeBindings == nil {
						rangeBindings = solver.SolveRelationSet(rangeSet, nameStore, []mentalese.Binding{{}})
					}

					rangeCount := mentalese.CountUniqueValues(rangeVariable.TermValue, rangeBindings)
					ok = rangeCount == uniqueRangeValues
				}

				if simpleQuantifier.TermValue == "which" {
					ok = true
				}
			}
		} else {
			// how many; todo
			ok = true
		}

		if ok {
			resultBindings = append(resultBindings, bindingGroup...)
		}
	}

	return resultBindings
}

func (solver ProblemSolver) groupBindings(quant mentalese.Relation, bindings []mentalese.Binding, rangeVariables[]mentalese.Term) [][]mentalese.Binding {

	mappedGroupedBindings := map[string][]mentalese.Binding{}

	rangeVar := rangeVariables[len(rangeVariables) - 1]

	// group all bindings by range variables
	for _, binding := range bindings {

		// key = A1/B2/C1
		key := ""
		sep := ""
		for i := 0; i < len(rangeVariables) - 1; i++ {
			rangeVariable := rangeVariables[i]
			val, ok := binding[rangeVariable.TermValue]
			if ok {
				key += sep + val.TermValue
				sep = "/"
			}
		}

		_, ok := binding[rangeVar.TermValue]
		if ok {

			_, found := mappedGroupedBindings[key]
			if !found {
				mappedGroupedBindings[key] = []mentalese.Binding{}
			}
			mappedGroupedBindings[key] = append(mappedGroupedBindings[key], binding)
		}
	}

	// unpack grouped bindings into a simple array of arrays
	var groupedBindings [][]mentalese.Binding
	for _, bindingArray := range mappedGroupedBindings {
		groupedBindings = append(groupedBindings, bindingArray)
	}

	return groupedBindings
}
