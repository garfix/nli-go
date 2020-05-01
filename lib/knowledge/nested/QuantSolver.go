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
	quantifierSet := quant.Arguments[mentalese.QuantQuantifierSetIndex].TermValueRelationSet
	rangeVariable := quant.Arguments[mentalese.QuantRangeVariableIndex].TermValue
	rangeSet := quant.Arguments[mentalese.QuantRangeSetIndex].TermValueRelationSet

	rangeBindings := base.collectRangeBindings(quantifierSet, rangeVariable, rangeSet, binding)
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

	rangeVariable := quant.Arguments[mentalese.QuantRangeVariableIndex].TermValue
	rangeCountVariable := quant.Arguments[mentalese.QuantRangeCountVariableIndex].TermValue
	scopeCountVariable := quant.Arguments[mentalese.QuantResultCountVariableIndex].TermValue
	quantifierSet := quant.Arguments[mentalese.QuantQuantifierSetIndex].TermValueRelationSet

	rangeCount := rangeBindings.GetDistinctValueCount(rangeVariable)
	scopeCount := scopeBindings.GetDistinctValueCount(rangeVariable)

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

	// R = count(range)
	// S = count(scope)

	// the => equals(S, 1)
	// some => greater_than(S, 0)
	// all => equals(S, R)
	// number(2) => equals(S, 2)
	// two or three => or(number(S, 2), number(S, 3))

	//isTheQuantifier := quantifierSet[0].Predicate == mentalese.PredicateThe
	//isAllQuantifier := quantifierSet[0].Predicate == mentalese.PredicateAll
	//isSomeQuantifier := quantifierSet[0].Predicate == mentalese.PredicateSome
	//
	//
	//count := 0
	//
	//// pick the number from the quantifierSet, if applicable
	//if quantifierSet[0].Predicate == mentalese.PredicateNumber {
	//	numberRelation := quantifierSet[0]
	//	count, _ = strconv.Atoi(numberRelation.Arguments[0].TermValue)
	//}
	//
	//if isTheQuantifier || isReference {
	//	count = 1
	//}
	//
	//if isTheQuantifier {
	//	// THE
	//	if scopeCount == 1 {
	//		return true
	//	} else {
	//		return false
	//	}
	//} else if isAllQuantifier {
	//	// ALL
	//	if scopeCount == len(rangeBindings) {
	//		return true
	//	} else {
	//		return false
	//	}
	//} else if isSomeQuantifier {
	//	// SOME
	//	if scopeCount > 0 {
	//		return true
	//	} else {
	//		return false
	//	}
	//} else if count > 0 {
	//	// A NUMBER
	//	if scopeCount == count {
	//		return true
	//	} else {
	//		return false
	//	}
	//} else {
	//	// NO SIMPLE QUANTIFIER
	//	// todo
	//	return true
	//}
}

func (base *SystemNestedStructureBase) collectRangeBindings(quantifier mentalese.RelationSet, rangeVariable string, rangeSet mentalese.RelationSet, binding mentalese.Binding) mentalese.Bindings {
	rangeBindings := []mentalese.Binding{}

	isTheQuantifier := quantifier[0].Predicate == "the"
	isReference := len(rangeSet) > 0 && rangeSet[0].Predicate == mentalese.PredicateBackReference

	if isTheQuantifier || isReference {

		// try the anaphora queue first
		refFound := false
		for _, group := range *base.dialogContext.AnaphoraQueue {

			ref := group[0]

			refBinding := binding.Merge(mentalese.Binding{ rangeVariable: mentalese.NewId(ref.Id, ref.EntityType)})
			//  empty range set for "it"
			if len(rangeSet) == 0 {
				refFound = true
				rangeBindings = mentalese.Bindings{ refBinding }
				break
			}
			if !base.quickAcceptabilityCheck(rangeVariable, ref.EntityType, rangeSet) {
				continue
			}
			testRangeBindings := base.solver.SolveRelationSet(rangeSet, mentalese.Bindings{refBinding})
			if len(testRangeBindings) == 1 {
				refFound = true
				rangeBindings = testRangeBindings
				break
			}
		}

		if !refFound {
			rangeBindings = base.solver.SolveRelationSet(rangeSet, mentalese.Bindings{binding})
		}

		if len(rangeBindings) != 1 {

			rangeIndex, found := base.rangeIndexClarification(rangeBindings, rangeVariable)
			if found {
				rangeBindings = rangeBindings[rangeIndex:rangeIndex + 1]
			} else {
				return mentalese.Bindings{}
			}
		}

	} else {
		rangeBindings = base.solver.SolveRelationSet(rangeSet, mentalese.Bindings{binding})
	}

	return rangeBindings
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