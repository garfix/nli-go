package function

import (
	"nli-go/lib/api"
	"nli-go/lib/central"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"strconv"
)

func (base *SystemSolverFunctionBase) solveAsync(messenger api.ProcessMessenger, set mentalese.RelationSet, bindings mentalese.BindingSet) (mentalese.BindingSet, bool) {

	return messenger.ExecuteChildStackFrameAsync(set, bindings)
}

// quant_check(quant() quant(), relationset)
func (base *SystemSolverFunctionBase) quantCheck(messenger api.ProcessMessenger, find mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
	if len(find.Arguments) != 2 {
		panic("quant_check(quants, scope) needs two arguments")
	}

	messenger.GetCursor().SetState("childIndex", 0)

	result, loading := base.solveQuantifiedRelations(messenger, find, binding, true)
	if loading {
		return mentalese.NewBindingSet()
	} else {
		return result
	}
}

// quant_foreach(quant() quant(), relationset)
func (base *SystemSolverFunctionBase) quantForeach(messenger api.ProcessMessenger, find mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
	if len(find.Arguments) != 2 {
		panic("quant_foreach(quants, scope) needs two arguments")
	}

	cursor := messenger.GetCursor()
	cursor.SetType(mentalese.FrameTypeLoop)
	cursor.SetState("childIndex", 0)

	result, loading := base.solveQuantifiedRelations(messenger, find, binding, false)
	if loading {
		return mentalese.NewBindingSet()
	} else {
		return result
	}
}

func (base *SystemSolverFunctionBase) quantOrderedList(messenger api.ProcessMessenger, quantList mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := quantList.BindSingle(binding)

	if !knowledge.Validate(bound, "rav", base.log) { return mentalese.NewBindingSet() }

	cursor := messenger.GetCursor()
	cursor.SetState("childIndex", 0)

	quant := bound.Arguments[0].TermValueRelationSet[0]
	orderFunction := bound.Arguments[1].TermValue
	listVariable := bound.Arguments[2].TermValue

	list, loading := base.getQuantifiedEntities(messenger, quant, orderFunction, binding)
	if loading {
		return mentalese.NewBindingSet()
	}

	newBinding := binding.Copy()
	newBinding.Set(listVariable, mentalese.NewTermList(list))

	return mentalese.InitBindingSet(newBinding)
}

func (base *SystemSolverFunctionBase) quantOrderSingle(quant mentalese.Relation, orderFunction string) mentalese.RelationSet {

	orderedQuant := quant.Copy()

	if quant.Predicate == mentalese.PredicateQuant {
		for len(orderedQuant.Arguments) < 3 {
			orderedQuant.Arguments = append(orderedQuant.Arguments, mentalese.NewTermAnonymousVariable())
		}
		orderedQuant.Arguments[2] = mentalese.NewTermAtom(orderFunction)
	} else {
		leftQuant := orderedQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
		rightQuant := orderedQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]

		orderedQuant.Arguments[mentalese.SeqFirstOperandIndex] = mentalese.NewTermRelationSet( base.quantOrderSingle(leftQuant, orderFunction) )
		orderedQuant.Arguments[mentalese.SeqSecondOperandIndex] = mentalese.NewTermRelationSet( base.quantOrderSingle(rightQuant, orderFunction) )
	}

	return mentalese.RelationSet{ orderedQuant }
}

func (base *SystemSolverFunctionBase) getQuantifiedEntities(messenger api.ProcessMessenger, quant mentalese.Relation, orderFunction string, binding mentalese.Binding) (mentalese.TermList, bool) {

	quantifiedEntities := mentalese.TermList{}

	if quant.Predicate == mentalese.PredicateOr {

		leftQuant := quant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
		rightQuant := quant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]

		leftEntities, loading := base.getEntities(messenger, leftQuant, orderFunction, binding)
		if loading {
			return mentalese.TermList{}, true
		}
		rightEntities, loading := base.getEntities(messenger, rightQuant, orderFunction, binding)
		if loading {
			return mentalese.TermList{}, true
		}
		combinedEntities := append(leftEntities, rightEntities...)
		uniqueEntities := unique(combinedEntities)
		orderedEntities, loading := base.entityQuickSort(messenger, uniqueEntities, orderFunction)
		if loading {
			return mentalese.TermList{}, true
		}
		quantifiedEntities, loading = base.applyQuantifierForOr(messenger, leftQuant, rightQuant, leftEntities, rightEntities, orderedEntities)
		if loading {
			return mentalese.TermList{}, true
		}

	} else if quant.Predicate == mentalese.PredicateAnd {

		leftQuant := quant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
		rightQuant := quant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]

		leftEntities, loading := base.getEntities(messenger, leftQuant, orderFunction, binding)
		if loading {
			return mentalese.TermList{}, true
		}
		rightEntities, loading := base.getEntities(messenger, rightQuant, orderFunction, binding)
		if loading {
			return mentalese.TermList{}, true
		}
		combinedEntities := append(leftEntities, rightEntities...)
		uniqueEntities := unique(combinedEntities)
		orderedEntities, loading := base.entityQuickSort(messenger, uniqueEntities, orderFunction)
		if loading {
			return mentalese.TermList{}, true
		}
		quantifiedEntities, loading = base.applyQuantifierForAnd(messenger, leftQuant, rightQuant, leftEntities, rightEntities, orderedEntities)
		if loading {
			return mentalese.TermList{}, true
		}

	} else if quant.Predicate != mentalese.PredicateQuant {

		base.log.AddError("First argument of a `quant_list` must be a `quant`")
		return mentalese.TermList{}, false

	} else {

		entities, loading := base.getEntities(messenger, quant, orderFunction, binding)
		if loading {
			return mentalese.TermList{}, true
		}
		orderedEntities, loading := base.entityQuickSort(messenger, entities, orderFunction)
		if loading {
			return mentalese.TermList{}, true
		}
		quantifiedEntities, loading = base.applyQuantifier(messenger, quant, orderedEntities)
		if loading {
			return mentalese.TermList{}, true
		}

	}

	return quantifiedEntities, false
}

func (base *SystemSolverFunctionBase) getEntities(messenger api.ProcessMessenger, quant mentalese.Relation, orderFunction string, binding mentalese.Binding) ([]mentalese.Term, bool) {

	if quant.Predicate != mentalese.PredicateQuant {
		return base.getQuantifiedEntities(messenger, quant, orderFunction, binding)
	}

	rangeSet := quant.Arguments[mentalese.QuantRangeSetIndex].TermValueRelationSet
	rangeVariable := quant.Arguments[mentalese.QuantRangeVariableIndex].TermValue
	//rangeBindings := base.solver.SolveRelationSet(rangeSet, mentalese.InitBindingSet(binding))
	rangeBindings, loading := base.solveAsync(messenger, rangeSet, mentalese.InitBindingSet(binding))
	return rangeBindings.GetIds(rangeVariable), loading
}

func unique(values []mentalese.Term) []mentalese.Term {
	uniq := []mentalese.Term{}

	for i := 0; i < len(values); i++ {
		value := values[i]
		if !containsId(uniq, value) {
			uniq = append(uniq, values[i])
		}
	}

	return uniq
}

func containsId(values []mentalese.Term, value mentalese.Term) bool {
	for i := 0; i < len(values); i++ {
		if values[i].Equals(value) {
			return true
		}
	}
	return false
}

// select either the left branch or the right branch, based on the entities and the quantifiers
func (base *SystemSolverFunctionBase) applyQuantifierForOr(messenger api.ProcessMessenger, leftQuant mentalese.Relation, rightQuant mentalese.Relation, leftValues []mentalese.Term, rightValues []mentalese.Term, orderedValues []mentalese.Term) ([]mentalese.Term, bool) {

	leftScopeCount := 0
	rightScopeCount := 0
	selectedLeftIds := []mentalese.Term{}
	selectedRightIds := []mentalese.Term{}
	selectedIds := []mentalese.Term{}
	ok := false
	loading := false

	for i := 0; i < len(orderedValues); i++ {
		value := orderedValues[i]
		if containsId(leftValues, value) {
			leftScopeCount++
			selectedLeftIds = append(selectedLeftIds, value)
			if leftQuant.Predicate != mentalese.PredicateQuant {
				ok = leftScopeCount == len(leftValues)
			} else {
				ok, loading = base.tryQuantifier(messenger, leftQuant, len(leftValues), leftScopeCount, true)
				if loading {
					return []mentalese.Term{}, true
				}
			}
			if ok {
				selectedIds = selectedLeftIds
				break
			}
		}
		if containsId(rightValues, value) {
			rightScopeCount++
			selectedRightIds = append(selectedRightIds, value)
			if rightQuant.Predicate != mentalese.PredicateQuant {
				ok = rightScopeCount == len(rightValues)
			} else {
				ok, loading = base.tryQuantifier(messenger, rightQuant, len(rightValues), rightScopeCount, true)
				if loading {
					return []mentalese.Term{}, true
				}
			}
			if ok {
				selectedIds = selectedRightIds
				break
			}
		}
	}

	return selectedIds, false
}

// select either the left branch or the right branch, based on the entities and the quantifiers
func (base *SystemSolverFunctionBase) applyQuantifierForAnd(messenger api.ProcessMessenger, leftQuant mentalese.Relation, rightQuant mentalese.Relation, leftValues []mentalese.Term, rightValues []mentalese.Term, orderedValues []mentalese.Term) ([]mentalese.Term, bool) {

	leftScopeCount := 0
	rightScopeCount := 0
	leftDone := false
	rightDone := false
	selectedIds := []mentalese.Term{}
	ok := false
	loading:= false

	for i := 0; i < len(orderedValues); i++ {
		term := orderedValues[i]
		if !leftDone {
			if containsId(leftValues, term) {
				selectedIds = append(selectedIds, term)
				leftScopeCount++
				if leftQuant.Predicate != mentalese.PredicateQuant {
					ok = leftScopeCount == len(leftValues)
				} else {
					ok, loading = base.tryQuantifier(messenger, leftQuant, len(leftValues), leftScopeCount, true)
					if loading {
						return []mentalese.Term{}, true
					}
				}
				if ok {
					leftDone = true
				}
			}
		}
		if !rightDone {
			if containsId(rightValues, term) {
				selectedIds = append(selectedIds, term)
				rightScopeCount++
				if rightQuant.Predicate != mentalese.PredicateQuant {
					ok = rightScopeCount == len(rightValues)
				} else {
					ok, loading = base.tryQuantifier(messenger, rightQuant, len(rightValues), rightScopeCount, true)
					if loading {
						return []mentalese.Term{}, true
					}
				}
				if ok {
					rightDone = true
				}
			}
		}
		if leftDone && rightDone { break }
	}

	return selectedIds, false
}

func (base *SystemSolverFunctionBase) applyQuantifier(messenger api.ProcessMessenger, quant mentalese.Relation, rangeValues []mentalese.Term) ([]mentalese.Term, bool) {
	rangeCount := len(rangeValues)
	scopeCount := 0
	for i := 0; i <= rangeCount; i++ {
		ok, loading := base.tryQuantifier(messenger, quant, rangeCount, i, true)
		if loading {
			return []mentalese.Term{}, true
		}
		if ok {
			scopeCount = i
			break
		}
	}

	return rangeValues[0:scopeCount], false
}

func (base *SystemSolverFunctionBase) solveQuantifiedRelations(messenger api.ProcessMessenger, find mentalese.Relation, binding mentalese.Binding, continueAfterEnough bool) (mentalese.BindingSet, bool) {

	quants := find.Arguments[0].TermValueRelationSet
	scope := find.Arguments[1].TermValueRelationSet

	return base.solveQuants(messenger, quants[0], scope, binding, continueAfterEnough)
}

func (base *SystemSolverFunctionBase) solveQuants(messenger api.ProcessMessenger, quant mentalese.Relation, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) (mentalese.BindingSet, bool) {

	result := mentalese.NewBindingSet()
	loading := false

	if quant.Predicate == mentalese.PredicateXor {

		result, loading = base.SolveXorQuant(messenger, quant, scopeSet, binding, continueAfterEnough)

	} else if quant.Predicate == mentalese.PredicateOr {

		result, loading = base.SolveOrQuant(messenger, quant, scopeSet, binding, continueAfterEnough)

	} else if quant.Predicate == mentalese.PredicateAnd {

		result, loading  = base.SolveAndQuant(messenger, quant, scopeSet, binding, continueAfterEnough)

	} else if quant.Predicate != mentalese.PredicateQuant {
		base.log.AddError("First argument of a `do` or `find` must contain only `quant`s")
		return mentalese.NewBindingSet(), false
	} else {

		result, loading = base.solveSimpleQuant(messenger, quant, scopeSet, binding, continueAfterEnough)

	}

	return result, loading
}

func (base *SystemSolverFunctionBase) solveSimpleQuant(messenger api.ProcessMessenger, quant mentalese.Relation, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) (mentalese.BindingSet, bool) {

	rangeSet := quant.Arguments[mentalese.QuantRangeSetIndex].TermValueRelationSet
	//rangeBindings := base.solver.SolveRelationSet(rangeSet, mentalese.InitBindingSet(binding))
	rangeBindings, loading := base.solveAsync(messenger, rangeSet, mentalese.InitBindingSet(binding))
	if loading {
		return mentalese.NewBindingSet(), true
	}

	scopeBindings, loading := base.solveScope(messenger, quant, scopeSet, rangeBindings, continueAfterEnough)
	if loading {
		return mentalese.NewBindingSet(), true
	}

	rangeVariable := quant.Arguments[mentalese.QuantRangeVariableIndex].TermValue

	rangeCount := rangeBindings.GetDistinctValueCount(rangeVariable)
	scopeCount := scopeBindings.GetDistinctValueCount(rangeVariable)

	success, loading := base.tryQuantifier(messenger, quant, rangeCount, scopeCount, true)
	if loading {
		return mentalese.NewBindingSet(), true
	}

	if success {
		return scopeBindings, false
	} else {
		return mentalese.NewBindingSet(), false
	}
}

func (base *SystemSolverFunctionBase) SolveAndQuant(messenger api.ProcessMessenger, xorQuant mentalese.Relation, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) (mentalese.BindingSet, bool) {

	leftQuant := xorQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
	rightQuant := xorQuant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]

	leftResultBindings, loading := base.solveQuants(messenger, leftQuant, scopeSet, binding, continueAfterEnough)
	if loading {
		return mentalese.NewBindingSet(), true
	}

	resultBindings := mentalese.NewBindingSet()
	for _, leftResultBinding := range leftResultBindings.GetAll() {
		rightResultBindings, loading := base.solveQuants(messenger, rightQuant, scopeSet, leftResultBinding, continueAfterEnough)
		if loading {
			return mentalese.NewBindingSet(), true
		}
		resultBindings.AddMultiple(rightResultBindings)
	}

	return resultBindings, false
}

func (base *SystemSolverFunctionBase) SolveOrQuant(messenger api.ProcessMessenger, orQuant mentalese.Relation, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) (mentalese.BindingSet, bool) {
	leftQuant := orQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
	rightQuant := orQuant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]
	leftResultBindings, loading := base.solveQuants(messenger, leftQuant, scopeSet, binding, continueAfterEnough)
	if loading {
		return mentalese.NewBindingSet(), true
	}
	rightResultBindings, loading := base.solveQuants(messenger, rightQuant, scopeSet, binding, continueAfterEnough)
	if loading {
		return mentalese.NewBindingSet(), true
	}

	newBindings := leftResultBindings.Copy()
	newBindings.AddMultiple(rightResultBindings)
	return newBindings, false
}

func (base *SystemSolverFunctionBase) SolveXorQuant(messenger api.ProcessMessenger, xorQuant mentalese.Relation, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool) (mentalese.BindingSet, bool) {
	leftQuant := xorQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
	resultBindings, loading := base.solveQuants(messenger, leftQuant, scopeSet, binding, continueAfterEnough)
	if loading {
		return mentalese.NewBindingSet(), true
	}
	if resultBindings.IsEmpty() {
		rightQuant := xorQuant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]
		resultBindings, loading = base.solveQuants(messenger, rightQuant, scopeSet, binding, continueAfterEnough)
		if loading {
			return mentalese.NewBindingSet(), true
		}
	}

	return resultBindings, false
}

func (base *SystemSolverFunctionBase) solveScope(messenger api.ProcessMessenger, quant mentalese.Relation, scopeSet []mentalese.Relation, rangeBindings mentalese.BindingSet, continueAfterEnough bool)  (mentalese.BindingSet, bool) {

	rangeVariable := quant.Arguments[mentalese.QuantRangeVariableIndex].TermValue
	scopeBindings := mentalese.NewBindingSet()
	groupedScopeBindings := []mentalese.BindingSet{}

	for _, rangeBinding := range rangeBindings.GetAll() {
		//singleScopeBindings := base.solver.SolveRelationSet(scopeSet, mentalese.InitBindingSet(rangeBinding))
		singleScopeBindings, loading := base.solveAsync(messenger, scopeSet, mentalese.InitBindingSet(rangeBinding))
		if loading {
			return mentalese.NewBindingSet(), true
		}

		if !singleScopeBindings.IsEmpty() {
			groupedScopeBindings = append(groupedScopeBindings, singleScopeBindings)
			scopeBindings.AddMultiple(singleScopeBindings)
		}

		value, found := rangeBinding.Get(rangeVariable)
		if found && value.IsId() {
			group := central.EntityReferenceGroup{central.CreateEntityReference(value.TermValue, value.TermSort) }
			base.anaphoraQueue.AddReferenceGroup(group)
		}

		rangeCount := rangeBindings.GetDistinctValueCount(rangeVariable)
		scopeCount := scopeBindings.GetDistinctValueCount(rangeVariable)

		ok, loading := base.tryQuantifier(messenger, quant, rangeCount, scopeCount, false)
		if loading {
			return mentalese.NewBindingSet(), true
		}
		if ok {
			if !continueAfterEnough {
				break
			}
		}
	}

	return scopeBindings, false
}

func (base *SystemSolverFunctionBase) tryQuantifier(messenger api.ProcessMessenger, quant mentalese.Relation, rangeCount int, scopeCount int, final bool) (bool, bool) {

	firstArgument := quant.Arguments[mentalese.QuantQuantifierIndex]

	// special case: the existential quantifier `some`
	if firstArgument.IsAtom() && firstArgument.TermValue == mentalese.AtomSome {
		if scopeCount == 0 {
			if base.log.Active() { base.log.AddDebug("Do/Find", "Quantifier Some mismatch: no results") }
			return false, false
		} else {
			return true, false
		}
	}

	// special case: the existential quantifier `none`
	if firstArgument.IsRelationSet() && len(firstArgument.TermValueRelationSet) == 0 {
		if final {
			return true, false
		} else {
			return false, false
		}
	}

	if !firstArgument.IsRelationSet() ||
		firstArgument.TermValueRelationSet[0].Predicate != mentalese.PredicateQuantifier {
		base.log.AddError("First argument of a `quant` must be a `quantifier`, but is " + firstArgument.String())
		return false, false
	}

	quantifier := firstArgument.TermValueRelationSet[0]

	scopeCountVariable := quantifier.Arguments[mentalese.QuantifierResultCountVariableIndex].TermValue
	rangeCountVariable := quantifier.Arguments[mentalese.QuantifierRangeCountVariableIndex].TermValue
	quantifierSet := quantifier.Arguments[mentalese.QuantifierSetIndex].TermValueRelationSet

	rangeVal := mentalese.NewTermString(strconv.Itoa(rangeCount))
	resultVal := mentalese.NewTermString(strconv.Itoa(scopeCount))

	b := mentalese.NewBinding()
	b.Set(rangeCountVariable, rangeVal)
	b.Set(scopeCountVariable, resultVal)

	//quantifierBindings := base.solver.SolveRelationSet(quantifierSet, mentalese.InitBindingSet(b))
	quantifierBindings, loading := base.solveAsync(messenger, quantifierSet, mentalese.InitBindingSet(b))
	if loading {
		return true, true
	}

	success := !quantifierBindings.IsEmpty()

	if !success {
		if base.log.Active() {
			base.log.AddDebug("Do/Find", "Quantifier mismatch")
			base.log.AddDebug("Do/Find", "  Range count: "+rangeCountVariable+" = "+strconv.Itoa(rangeCount))
			base.log.AddDebug("Do/Find", "  Scope count: "+scopeCountVariable+" = "+strconv.Itoa(scopeCount))
			base.log.AddDebug("Do/Find", "  Quantifier check: "+quantifierSet.String())
		}
	}

	return success, false
}

func (base *SystemSolverFunctionBase) quickAcceptabilityCheck(variable string, sort string, relations mentalese.RelationSet) bool {

	accepted := false

	for _, relation := range relations {
		for i, argument := range relation.Arguments {
			if argument.IsVariable() && argument.TermValue == variable {
				argumentEntityType := base.meta.GetSort(relation.Predicate, i)

				if argumentEntityType == "" || base.meta.MatchesSort(argumentEntityType, sort) {
					accepted = true
					break
				}
			}
		}
	}

	return accepted
}
