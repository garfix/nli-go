package nested

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"strconv"
)

// quant_check(quant() quant(), relationset)
func (base *SystemNestedStructureBase) SolveQuantCheck(find mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {
	if len(find.Arguments) != 2 {
		panic("quant_check(quants, scope) needs two arguments")
	}
	return base.solveQuantifiedRelations(find, binding, true)
}

// quant_foreach(quant() quant(), relationset)
func (base *SystemNestedStructureBase) SolveQuantForeach(find mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {
	if len(find.Arguments) != 2 {
		panic("quant_foreach(quants, scope) needs two arguments")
	}
	return base.solveQuantifiedRelations(find, binding, false)
}

func (base *SystemNestedStructureBase) SolveQuantOrderedList(quantList mentalese.Relation, binding mentalese.Binding) mentalese.Bindings {

	bound := quantList.BindSingle(binding)

	if !knowledge.Validate(bound, "rav", base.log) { return nil }

	quant := bound.Arguments[0].TermValueRelationSet[0]
	orderFunction := bound.Arguments[1].TermValue
	listVariable := bound.Arguments[2].TermValue

	list := base.getQuantifiedEntities(quant, orderFunction, binding)

	newBinding := binding.Copy()
	newBinding[listVariable] = mentalese.NewTermList(list)

	return mentalese.Bindings{ newBinding }
}

func (base *SystemNestedStructureBase) getQuantifiedEntities(quant mentalese.Relation, orderFunction string, binding mentalese.Binding) mentalese.TermList {

	quantifiedEntities := mentalese.TermList{}

	if quant.Predicate == mentalese.PredicateOr {

		leftQuant := quant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
		rightQuant := quant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]

		leftEntities := base.getEntities(leftQuant, orderFunction, binding)
		rightEntities := base.getEntities(rightQuant, orderFunction, binding)
		combinedEntities := append(leftEntities, rightEntities...)
		uniqueEntities := unique(combinedEntities)
		orderedEntities := base.entityQuickSort(uniqueEntities, orderFunction)
		quantifiedEntities = base.applyQuantifierForOr(leftQuant, rightQuant, leftEntities, rightEntities, orderedEntities)

	} else if quant.Predicate == mentalese.PredicateAnd {

		leftQuant := quant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
		rightQuant := quant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]

		leftEntities := base.getEntities(leftQuant, orderFunction, binding)
		rightEntities := base.getEntities(rightQuant, orderFunction, binding)
		combinedEntities := append(leftEntities, rightEntities...)
		uniqueEntities := unique(combinedEntities)
		orderedEntities := base.entityQuickSort(uniqueEntities, orderFunction)
		quantifiedEntities = base.applyQuantifierForAnd(leftQuant, rightQuant, leftEntities, rightEntities, orderedEntities)

	} else if quant.Predicate != mentalese.PredicateQuant {

		base.log.AddError("First argument of a `quant_list` must be a `quant`")
		return mentalese.TermList{}

	} else {

		entities := base.getEntities(quant, orderFunction, binding)
		orderedEntities := base.entityQuickSort(entities, orderFunction)
		quantifiedEntities = base.applyQuantifier(quant, orderedEntities)

	}

	return quantifiedEntities
}

func (base *SystemNestedStructureBase) getEntities(quant mentalese.Relation, orderFunction string, binding mentalese.Binding) []mentalese.Term {

	if quant.Predicate != mentalese.PredicateQuant {
		return base.getQuantifiedEntities(quant, orderFunction, binding)
	}

	rangeSet := quant.Arguments[mentalese.QuantRangeSetIndex].TermValueRelationSet
	rangeVariable := quant.Arguments[mentalese.QuantRangeVariableIndex].TermValue
	rangeBindings := base.solver.SolveRelationSet(rangeSet, mentalese.Bindings{binding})
	return rangeBindings.GetIds(rangeVariable)
}

func unique(values []mentalese.Term) []mentalese.Term {
	uniq := []mentalese.Term{}

	for i := 0; i < len(values); i++ {
		value := values[i].TermValue
		if !containsId(uniq, value) {
			uniq = append(uniq, values[i])
		}
	}

	return uniq
}

func containsId(values []mentalese.Term, value string) bool {
	for i := 0; i < len(values); i++ {
		if values[i].TermValue == value {
			return true
		}
	}
	return false
}

// select either the left branch or the right branch, based on the entities and the quantifiers
func (base *SystemNestedStructureBase) applyQuantifierForOr(leftQuant mentalese.Relation, rightQuant mentalese.Relation, leftValues []mentalese.Term, rightValues []mentalese.Term, orderedValues []mentalese.Term) []mentalese.Term {

	leftScopeCount := 0
	rightScopeCount := 0
	selectedLeftIds := []mentalese.Term{}
	selectedRightIds := []mentalese.Term{}
	selectedIds := []mentalese.Term{}
	ok := false

	for i := 0; i < len(orderedValues); i++ {
		value := orderedValues[i]
		if containsId(leftValues, value.TermValue) {
			leftScopeCount++
			selectedLeftIds = append(selectedLeftIds, value)
			if leftQuant.Predicate != mentalese.PredicateQuant {
				ok = leftScopeCount == len(leftValues)
			} else {
				ok = base.tryQuantifier(leftQuant, len(leftValues), leftScopeCount, true)
			}
			if ok {
				selectedIds = selectedLeftIds
				break
			}
		}
		if containsId(rightValues, value.TermValue) {
			rightScopeCount++
			selectedRightIds = append(selectedRightIds, value)
			if rightQuant.Predicate != mentalese.PredicateQuant {
				ok = rightScopeCount == len(rightValues)
			} else {
				ok = base.tryQuantifier(rightQuant, len(rightValues), rightScopeCount, true)
			}
			if ok {
				selectedIds = selectedRightIds
				break
			}
		}
	}

	return selectedIds
}

// select either the left branch or the right branch, based on the entities and the quantifiers
func (base *SystemNestedStructureBase) applyQuantifierForAnd(leftQuant mentalese.Relation, rightQuant mentalese.Relation, leftValues []mentalese.Term, rightValues []mentalese.Term, orderedValues []mentalese.Term) []mentalese.Term {

	leftScopeCount := 0
	rightScopeCount := 0
	leftDone := false
	rightDone := false
	selectedIds := []mentalese.Term{}
	ok := false

	for i := 0; i < len(orderedValues); i++ {
		term := orderedValues[i]
		value := term.TermValue
		if !leftDone {
			if containsId(leftValues, value) {
				selectedIds = append(selectedIds, term)
				leftScopeCount++
				if leftQuant.Predicate != mentalese.PredicateQuant {
					ok = leftScopeCount == len(leftValues)
				} else {
					ok = base.tryQuantifier(leftQuant, len(leftValues), leftScopeCount, true)
				}
				if ok {
					leftDone = true
				}
			}
		}
		if !rightDone {
			if containsId(rightValues, value) {
				selectedIds = append(selectedIds, term)
				rightScopeCount++
				if rightQuant.Predicate != mentalese.PredicateQuant {
					ok = rightScopeCount == len(rightValues)
				} else {
					ok = base.tryQuantifier(rightQuant, len(rightValues), rightScopeCount, true)
				}
				if ok {
					rightDone = true
				}
			}
		}
		if leftDone && rightDone { break }
	}

	return selectedIds
}

func (base *SystemNestedStructureBase) applyQuantifier(quant mentalese.Relation, rangeValues []mentalese.Term) []mentalese.Term {
	rangeCount := len(rangeValues)
	scopeCount := 0
	for i := 0; i < rangeCount; i++ {
		ok := base.tryQuantifier(quant, rangeCount, i, true)
		if ok {
			scopeCount = i
			break
		}
	}

	return rangeValues[0:scopeCount]
}

func (base *SystemNestedStructureBase) solveQuantifiedRelations(find mentalese.Relation, binding mentalese.Binding, continueAfterEnough bool) mentalese.Bindings {

	quants := find.Arguments[0].TermValueRelationSet
	scope := find.Arguments[1].TermValueRelationSet

	return base.solveQuants(quants[0], scope, binding, continueAfterEnough)
}

func (base *SystemNestedStructureBase) solveQuants(quant mentalese.Relation, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) mentalese.Bindings {

	base.log.StartProduction("Quant", quant.String())

	result := mentalese.Bindings{}

	if quant.Predicate == mentalese.PredicateXor {

		result = base.SolveXorQuant(quant, scopeSet, binding, continueAfterEnough)

	} else if quant.Predicate == mentalese.PredicateOr {

		result = base.SolveOrQuant(quant, scopeSet, binding, continueAfterEnough)

	} else if quant.Predicate == mentalese.PredicateAnd {

		result = base.SolveAndQuant(quant, scopeSet, binding, continueAfterEnough)

	} else if quant.Predicate != mentalese.PredicateQuant {
		base.log.AddError("First argument of a `do` or `find` must contain only `quant`s")
		return mentalese.Bindings{}
	} else {

		result = base.solveSimpleQuant(quant, scopeSet, binding, continueAfterEnough)

	}

	base.log.EndProduction("Quant", result.String())

	return result
}

func (base *SystemNestedStructureBase) solveSimpleQuant(quant mentalese.Relation, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) mentalese.Bindings {

	rangeSet := quant.Arguments[mentalese.QuantRangeSetIndex].TermValueRelationSet
	rangeBindings := base.solver.SolveRelationSet(rangeSet, mentalese.Bindings{binding})

	scopeBindings := mentalese.Bindings{}

	scopeBindings = base.solveScope(quant, scopeSet, rangeBindings, continueAfterEnough)

	rangeVariable := quant.Arguments[mentalese.QuantRangeVariableIndex].TermValue

	rangeCount := rangeBindings.GetDistinctValueCount(rangeVariable)
	scopeCount := scopeBindings.GetDistinctValueCount(rangeVariable)

	success := base.tryQuantifier(quant, rangeCount, scopeCount, true)

	if success {
		return scopeBindings
	} else {
		return mentalese.Bindings{}
	}
}

func (base *SystemNestedStructureBase) SolveAndQuant(xorQuant mentalese.Relation, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) mentalese.Bindings {

	leftQuant := xorQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
	rightQuant := xorQuant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]

	leftResultBindings := base.solveQuants(leftQuant, scopeSet, binding, continueAfterEnough)

	resultBindings := mentalese.Bindings{}
	for _, leftResultBinding := range leftResultBindings {
		rightResultBindings := base.solveQuants(rightQuant, scopeSet, leftResultBinding, continueAfterEnough)
		resultBindings = append(resultBindings, rightResultBindings...)
	}

	return resultBindings
}

func (base *SystemNestedStructureBase) SolveOrQuant(orQuant mentalese.Relation, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) mentalese.Bindings {
	leftQuant := orQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
	rightQuant := orQuant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]
	leftResultBindings := base.solveQuants(leftQuant, scopeSet, binding, continueAfterEnough)
	rightResultBindings := base.solveQuants(rightQuant, scopeSet, binding, continueAfterEnough)

	return append(leftResultBindings, rightResultBindings...)
}

func (base *SystemNestedStructureBase) SolveXorQuant(xorQuant mentalese.Relation, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) mentalese.Bindings {
	leftQuant := xorQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
	resultBindings := base.solveQuants(leftQuant, scopeSet, binding, continueAfterEnough)
	if len(resultBindings) == 0 {
		rightQuant := xorQuant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]
		resultBindings = base.solveQuants(rightQuant, scopeSet, binding, continueAfterEnough)
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

		rangeCount := rangeBindings.GetDistinctValueCount(rangeVariable)
		scopeCount := scopeBindings.GetDistinctValueCount(rangeVariable)

		if base.tryQuantifier(quant, rangeCount, scopeCount, false) {
			if !continueAfterEnough {
				break
			}
		}
	}

	return scopeBindings
}

func (base *SystemNestedStructureBase) tryQuantifier(quant mentalese.Relation, rangeCount int, scopeCount int, final bool) bool {

	firstArgument := quant.Arguments[mentalese.QuantQuantifierIndex]

	// special case: the existential quantifier `some`
	if firstArgument.IsAtom() && firstArgument.TermValue == mentalese.AtomSome {
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

	rangeVal := mentalese.NewTermString(strconv.Itoa(rangeCount))
	resultVal := mentalese.NewTermString(strconv.Itoa(scopeCount))

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
				argumentEntityType := base.meta.GetEntityType(relation.Predicate, i)

				if argumentEntityType == "" || base.meta.MatchesSort(argumentEntityType, entityType) {
					accepted = true
					break
				}

				//if  argumentEntityType == "" || argumentEntityType == entityType {
				//	accepted = true
				//	break
				//}
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