package nested

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
)

// find(quant() quant(), relationset)
func (base *SystemNestedStructureBase) SolveFind(find mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {
	if len(find.Arguments) != 2 {
		panic("find(quants, scope) needs two arguments")
	}
	return base.solveQuantifiedRelations(find, binding, true)
}

// do(quant() quant(), relationset)
func (base *SystemNestedStructureBase) SolveDo(find mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {
	if len(find.Arguments) != 2 {
		panic("do(quants, scope) needs two arguments")
	}
	return base.solveQuantifiedRelations(find, binding, false)
}

func (base *SystemNestedStructureBase) solveQuantifiedRelations(find mentalese.Relation, binding mentalese.Binding, continueAfterEnough bool) mentalese.Bindings {

	quants := find.Arguments[0].TermValueRelationSet
	scope := find.Arguments[1].TermValueRelationSet

	return base.solveQuants(quants, scope, binding, continueAfterEnough)
}

func (base *SystemNestedStructureBase) solveQuants(quants mentalese.RelationSet, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) mentalese.Bindings {

	quant := quants[0]

	if quant.Predicate != mentalese.PredicateQuant {
		base.log.AddError("First argument of a `do` or `find` must contain only `quant`s")
		return mentalese.Bindings{}
	}

	rangeSet := quant.Arguments[mentalese.QuantRangeSetIndex].TermValueRelationSet

	rangeBindings := base.solver.SolveRelationSet(rangeSet, mentalese.Bindings{binding})
	isReference := false

	for _, r := range rangeSet {
		if r.Predicate == mentalese.PredicateBackReference {
			isReference = true
		}
	}

	scopeBindings := mentalese.Bindings{}
	resultCount := 0

	if len(quants) == 1 {
		scopeBindings = base.solveScope(quant, scopeSet, rangeBindings, isReference, continueAfterEnough)
		resultCount = len(scopeBindings)
	} else {
		for _, rangeBinding := range rangeBindings {
			singleScopeBindings := base.solveQuants(quants[1:], scopeSet, rangeBinding, continueAfterEnough)

			if len(singleScopeBindings) > 0 {
				scopeBindings = append(scopeBindings, singleScopeBindings...)
				resultCount++
			}
		}
	}

	success := base.tryQuantifier(quant, rangeBindings, scopeBindings, isReference)

	if success {
		return scopeBindings
	} else {
		return mentalese.Bindings{}
	}
}

func (base *SystemNestedStructureBase) solveScope(quant mentalese.Relation, scopeSet []mentalese.Relation, rangeBindings mentalese.Bindings, isReference bool, continueAfterEnough bool)  mentalese.Bindings {

	rangeVariable := quant.Arguments[mentalese.QuantRangeVariableIndex].TermValue
	scopeBindings := mentalese.Bindings{}
	groupedScopeBindings := []mentalese.Bindings{}

	for _, rangeBinding := range rangeBindings {
		singleScopeBindings := base.solver.SolveRelationSet(scopeSet, mentalese.Bindings{ rangeBinding })

		if len(singleScopeBindings) > 0 {
			groupedScopeBindings = append(groupedScopeBindings, singleScopeBindings)
			scopeBindings = append(scopeBindings, singleScopeBindings...)
		}

		value, found := rangeBinding[rangeVariable]
		if found && value.IsId() {
			group := central.EntityReferenceGroup{central.CreateEntityReference(value.TermValue, value.TermEntityType) }
			base.dialogContext.AnaphoraQueue.AddReferenceGroup(group)
		}

		if base.tryQuantifier(quant, rangeBindings, scopeBindings, isReference) {
			if !continueAfterEnough {
				break
			}
		}
	}

	return scopeBindings
}

func (base *SystemNestedStructureBase) tryQuantifier(quant mentalese.Relation, rangeBindings mentalese.Bindings, scopeBindings mentalese.Bindings, isReference bool) bool {

	firstArgument := quant.Arguments[mentalese.QuantQuantifierIndex]

	rangeVariable := quant.Arguments[mentalese.QuantRangeVariableIndex].TermValue

	rangeCount := rangeBindings.GetDistinctValueCount(rangeVariable)
	scopeCount := scopeBindings.GetDistinctValueCount(rangeVariable)

	// special case: the existential quantifier `some`
	if firstArgument.IsAtom() && firstArgument.TermValue == mentalese.PredicateQuantifierSome {
		if scopeCount == 0 {
			base.log.AddProduction("Do/Find", "Quantifier Some mismatch: no results")
			return false
		} else {
			return true
		}
	}

	if !firstArgument.IsRelationSet() ||
		len(firstArgument.TermValueRelationSet) != 1 ||
		firstArgument.TermValueRelationSet[0].Predicate != mentalese.PredicateQuantifier {
		base.log.AddError("First argument of a `quant` must be a `quantifier`")
		return false
	}

	quantifier := quant.Arguments[mentalese.QuantQuantifierIndex].TermValueRelationSet[0]

	scopeCountVariable := quantifier.Arguments[mentalese.QuantifierResultCountVariableIndex].TermValue
	rangeCountVariable := quantifier.Arguments[mentalese.QuantifierRangeCountVariableIndex].TermValue
	quantifierSet := quantifier.Arguments[mentalese.QuantifierSetIndex].TermValueRelationSet

	rangeVal := mentalese.NewString(strconv.Itoa(rangeCount))
	resultVal := mentalese.NewString(strconv.Itoa(scopeCount))

	quantifierBindings := base.solver.SolveRelationSet(quantifierSet, mentalese.Bindings{
		{ rangeCountVariable: rangeVal, scopeCountVariable: resultVal },
	})

	success := len(quantifierBindings) > 0

	if !success {
		base.log.AddProduction("Do/Find", "Quantifier mismatch")
		base.log.AddProduction("Do/Find", "  Range count: "+rangeCountVariable+" = "+strconv.Itoa(rangeCount))
		base.log.AddProduction("Do/Find", "  Scope count: "+scopeCountVariable+" = "+strconv.Itoa(scopeCount))
		base.log.AddProduction("Do/Find", "  Quantifier check: "+quantifierSet.String())
	}

	return success
}

func (base *SystemNestedStructureBase) quickAcceptabilityCheck(variable string, entityType string, relations mentalese.RelationSet) bool {

	accepted := false

	for _, relation := range relations {
		for i, argument := range relation.Arguments {
			if argument.IsVariable() && argument.TermValue == variable {
				argumentEntityType := base.predicates.GetEntityType(relation.Predicate, i)
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
func (base *SystemNestedStructureBase) rangeIndexClarification(rangeBindings mentalese.Bindings, rangeVariable string) (int, bool) {

	options := common.NewOptions()

	answer, found := base.dialogContext.GetAnswerToOpenQuestion()

	if found {

		index, _ := strconv.Atoi(answer)
		base.dialogContext.RemoveAnswerToOpenQuestion()

		return index, true

	} else {

		base.log.SetClarificationRequest("I don't understand which one you mean", options)
		return 0, false
	}
}