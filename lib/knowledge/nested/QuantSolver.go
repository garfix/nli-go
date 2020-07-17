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

	return base.solveQuants(quants[0], quants[1:], scope, binding, continueAfterEnough)
}

func (base *SystemNestedStructureBase) solveQuants(quant mentalese.Relation, restQuants mentalese.RelationSet, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) mentalese.Bindings {

	base.log.StartProduction("Quant", quant.String())

	result := mentalese.Bindings{}

	if quant.Predicate == mentalese.PredicateXor {

		result = base.SolveXorQuant(quant, restQuants, scopeSet, binding, continueAfterEnough)

	} else if quant.Predicate == mentalese.PredicateOr {

		result = base.SolveOrQuant(quant, restQuants, scopeSet, binding, continueAfterEnough)

	} else if quant.Predicate == mentalese.PredicateAnd {

		result = base.SolveAndQuant(quant, restQuants, scopeSet, binding, continueAfterEnough)

	} else if quant.Predicate != mentalese.PredicateQuant {
		base.log.AddError("First argument of a `do` or `find` must contain only `quant`s")
		return mentalese.Bindings{}
	} else {

		result = base.solveSimpleQuant(quant, restQuants, scopeSet, binding, continueAfterEnough)

	}

	base.log.EndProduction("Quant", result.String())

	return result
}

func (base *SystemNestedStructureBase) solveSimpleQuant(quant mentalese.Relation, restQuants mentalese.RelationSet, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) mentalese.Bindings {

	rangeSet := quant.Arguments[mentalese.QuantRangeSetIndex].TermValueRelationSet
	rangeBindings := base.solver.SolveRelationSet(rangeSet, mentalese.Bindings{binding})

	scopeBindings := mentalese.Bindings{}

	if len(restQuants) == 0 {
		scopeBindings = base.solveScope(quant, scopeSet, rangeBindings, continueAfterEnough)
	} else {
		for _, rangeBinding := range rangeBindings {
			singleScopeBindings := base.solveQuants(restQuants[0], restQuants[1:], scopeSet, rangeBinding, continueAfterEnough)
			scopeBindings = append(scopeBindings, singleScopeBindings...)
		}
	}

	success := base.tryQuantifier(quant, rangeBindings, scopeBindings, true)

	if success {
		return scopeBindings
	} else {
		return mentalese.Bindings{}
	}
}

func (base *SystemNestedStructureBase) SolveAndQuant(xorQuant mentalese.Relation, restQuants mentalese.RelationSet, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) mentalese.Bindings {
	leftQuant := xorQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
	rightQuant := xorQuant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]

	leftResultBindings := base.solveQuants(leftQuant, restQuants, scopeSet, binding, continueAfterEnough)

	resultBindings := mentalese.Bindings{}
	for _, leftResultBinding := range leftResultBindings {
		rightResultBindings := base.solveQuants(rightQuant, restQuants, scopeSet, leftResultBinding, continueAfterEnough)
		resultBindings = append(resultBindings, rightResultBindings...)
	}

	return resultBindings
}

func (base *SystemNestedStructureBase) SolveOrQuant(orQuant mentalese.Relation, restQuants mentalese.RelationSet, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) mentalese.Bindings {
	leftQuant := orQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
	rightQuant := orQuant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]
	leftResultBindings := base.solveQuants(leftQuant, restQuants, scopeSet, binding, continueAfterEnough)
	rightResultBindings := base.solveQuants(rightQuant, restQuants, scopeSet, binding, continueAfterEnough)

	return append(leftResultBindings, rightResultBindings...)
}

func (base *SystemNestedStructureBase) SolveXorQuant(xorQuant mentalese.Relation, restQuants mentalese.RelationSet, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) mentalese.Bindings {
	leftQuant := xorQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
	resultBindings := base.solveQuants(leftQuant, restQuants, scopeSet, binding, continueAfterEnough)
	if len(resultBindings) == 0 {
		rightQuant := xorQuant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]
		resultBindings = base.solveQuants(rightQuant, restQuants, scopeSet, binding, continueAfterEnough)
	}

	return resultBindings
}

func (base *SystemNestedStructureBase) solveScope(quant mentalese.Relation, scopeSet []mentalese.Relation, rangeBindings mentalese.Bindings, continueAfterEnough bool)  mentalese.Bindings {

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

		if base.tryQuantifier(quant, rangeBindings, scopeBindings, false) {
			if !continueAfterEnough {
				break
			}
		}
	}

	return scopeBindings
}

func (base *SystemNestedStructureBase) tryQuantifier(quant mentalese.Relation, rangeBindings mentalese.Bindings, scopeBindings mentalese.Bindings, final bool) bool {

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

	// special case: the existential quantifier `none`
	if firstArgument.IsRelationSet() && len(firstArgument.TermValueRelationSet) == 0 {
		if final {
			return true
		} else {
			return false
		}
	}

	if !firstArgument.IsRelationSet() ||
		firstArgument.TermValueRelationSet[0].Predicate != mentalese.PredicateQuantifier {
		base.log.AddError("First argument of a `quant` must be a `quantifier`, but is " + firstArgument.String())
		return false
	}

	quantifier := firstArgument.TermValueRelationSet[0]

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