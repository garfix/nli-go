package function

import (
	"nli-go/lib/api"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"strconv"
)

// quant_check(quant() quant(), relationset)
func (base *SystemSolverFunctionBase) quantCheck(messenger api.ProcessMessenger, find mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
	if len(find.Arguments) != 2 {
		panic("quant_check(quants, scope) needs two arguments")
	}

	return base.solveQuantifiedRelations(messenger, find, binding, true, mentalese.TermList{})
}

// quant_foreach(quant() quant(), relationset)
func (base *SystemSolverFunctionBase) quantForeach(messenger api.ProcessMessenger, find mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
	if len(find.Arguments) != 2 {
		panic("quant_foreach(quants, scope) needs two arguments")
	}

	cursor := messenger.GetCursor()
	cursor.SetType(mentalese.FrameTypeLoop)

	return base.solveQuantifiedRelations(messenger, find, binding, false, mentalese.TermList{})
}

func (base *SystemSolverFunctionBase) quantOrderedList(messenger api.ProcessMessenger, quantList mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := quantList.BindSingle(binding)

	if !knowledge.Validate(bound, "rav", base.log) {
		return mentalese.NewBindingSet()
	}

	quant := bound.Arguments[0].TermValueRelationSet[0]
	orderFunction := bound.Arguments[1].TermValue
	listVariable := bound.Arguments[2].TermValue

	list := base.getQuantifiedEntities(messenger, quant, orderFunction, binding)

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

		orderedQuant.Arguments[mentalese.SeqFirstOperandIndex] = mentalese.NewTermRelationSet(base.quantOrderSingle(leftQuant, orderFunction))
		orderedQuant.Arguments[mentalese.SeqSecondOperandIndex] = mentalese.NewTermRelationSet(base.quantOrderSingle(rightQuant, orderFunction))
	}

	return mentalese.RelationSet{orderedQuant}
}

func (base *SystemSolverFunctionBase) getQuantifiedEntities(messenger api.ProcessMessenger, quant mentalese.Relation, orderFunction string, binding mentalese.Binding) mentalese.TermList {

	quantifiedEntities := mentalese.TermList{}

	if quant.Predicate == mentalese.PredicateOr {

		leftQuant := quant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
		rightQuant := quant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]

		leftEntities := base.getEntities(messenger, leftQuant, orderFunction, binding)
		rightEntities := base.getEntities(messenger, rightQuant, orderFunction, binding)
		combinedEntities := append(leftEntities, rightEntities...)
		uniqueEntities := unique(combinedEntities)
		orderedEntities := base.entityQuickSort(messenger, uniqueEntities, orderFunction)
		quantifiedEntities = base.applyQuantifierForOr(messenger, leftQuant, rightQuant, leftEntities, rightEntities, orderedEntities)

	} else if quant.Predicate == mentalese.PredicateAnd {

		leftQuant := quant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
		rightQuant := quant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]

		leftEntities := base.getEntities(messenger, leftQuant, orderFunction, binding)
		rightEntities := base.getEntities(messenger, rightQuant, orderFunction, binding)
		combinedEntities := append(leftEntities, rightEntities...)
		uniqueEntities := unique(combinedEntities)
		orderedEntities := base.entityQuickSort(messenger, uniqueEntities, orderFunction)
		quantifiedEntities, _ = base.applyQuantifierForAnd(messenger, leftQuant, rightQuant, leftEntities, rightEntities, orderedEntities)

	} else if quant.Predicate != mentalese.PredicateQuant {

		base.log.AddError("First argument of a `quant_list` must be a `quant`")
		return mentalese.TermList{}

	} else {

		entities := base.getEntities(messenger, quant, orderFunction, binding)
		orderedEntities := base.entityQuickSort(messenger, entities, orderFunction)
		quantifiedEntities = base.applyQuantifier(messenger, quant, orderedEntities)
	}

	return quantifiedEntities
}

func (base *SystemSolverFunctionBase) getEntities(messenger api.ProcessMessenger, quant mentalese.Relation, orderFunction string, binding mentalese.Binding) []mentalese.Term {

	if quant.Predicate != mentalese.PredicateQuant {
		return base.getQuantifiedEntities(messenger, quant, orderFunction, binding)
	}

	rangeSet := quant.Arguments[mentalese.QuantRangeSetIndex].TermValueRelationSet
	rangeVariable := quant.Arguments[mentalese.QuantRangeVariableIndex].TermValue
	rangeBindings := messenger.ExecuteChildStackFrame(rangeSet, mentalese.InitBindingSet(binding))
	return rangeBindings.GetIds(rangeVariable)
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
func (base *SystemSolverFunctionBase) applyQuantifierForOr(messenger api.ProcessMessenger, leftQuant mentalese.Relation, rightQuant mentalese.Relation, leftValues []mentalese.Term, rightValues []mentalese.Term, orderedValues []mentalese.Term) []mentalese.Term {

	leftScopeCount := 0
	rightScopeCount := 0
	selectedLeftIds := []mentalese.Term{}
	selectedRightIds := []mentalese.Term{}
	selectedIds := []mentalese.Term{}
	ok := false

	for i := 0; i < len(orderedValues); i++ {
		value := orderedValues[i]
		if containsId(leftValues, value) {
			leftScopeCount++
			selectedLeftIds = append(selectedLeftIds, value)
			if leftQuant.Predicate != mentalese.PredicateQuant {
				ok = leftScopeCount == len(leftValues)
			} else {
				ok = base.tryQuantifier(messenger, leftQuant, len(leftValues), leftScopeCount, true)
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
				ok = base.tryQuantifier(messenger, rightQuant, len(rightValues), rightScopeCount, true)
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
func (base *SystemSolverFunctionBase) applyQuantifierForAnd(messenger api.ProcessMessenger, leftQuant mentalese.Relation, rightQuant mentalese.Relation, leftValues []mentalese.Term, rightValues []mentalese.Term, orderedValues []mentalese.Term) ([]mentalese.Term, bool) {

	leftScopeCount := 0
	rightScopeCount := 0
	leftDone := false
	rightDone := false
	selectedIds := []mentalese.Term{}
	ok := false

	for i := 0; i < len(orderedValues); i++ {
		term := orderedValues[i]
		if !leftDone {
			if containsId(leftValues, term) {
				selectedIds = append(selectedIds, term)
				leftScopeCount++
				if leftQuant.Predicate != mentalese.PredicateQuant {
					ok = leftScopeCount == len(leftValues)
				} else {
					ok = base.tryQuantifier(messenger, leftQuant, len(leftValues), leftScopeCount, true)
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
					ok = base.tryQuantifier(messenger, rightQuant, len(rightValues), rightScopeCount, true)
				}
				if ok {
					rightDone = true
				}
			}
		}
		if leftDone && rightDone {
			break
		}
	}

	return selectedIds, false
}

func (base *SystemSolverFunctionBase) applyQuantifier(messenger api.ProcessMessenger, quant mentalese.Relation, rangeValues []mentalese.Term) []mentalese.Term {
	rangeCount := len(rangeValues)
	scopeCount := 0
	for i := 0; i <= rangeCount; i++ {
		ok := base.tryQuantifier(messenger, quant, rangeCount, i, true)
		if ok {
			scopeCount = i
			break
		}
	}

	return rangeValues[0:scopeCount]
}

func (base *SystemSolverFunctionBase) solveQuantifiedRelations(messenger api.ProcessMessenger, find mentalese.Relation, binding mentalese.Binding, continueAfterEnough bool, toMatch mentalese.TermList) mentalese.BindingSet {

	quants := find.Arguments[0].TermValueRelationSet
	scope := find.Arguments[1].TermValueRelationSet

	return base.solveQuants(messenger, quants[0], scope, binding, continueAfterEnough, toMatch)
}

func (base *SystemSolverFunctionBase) solveQuants(messenger api.ProcessMessenger, quant mentalese.Relation, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool, toMatch mentalese.TermList) mentalese.BindingSet {

	result := mentalese.NewBindingSet()

	if quant.Predicate == mentalese.PredicateXor {

		result = base.SolveXorQuant(messenger, quant, scopeSet, binding, continueAfterEnough, toMatch)

	} else if quant.Predicate == mentalese.PredicateOr {

		result = base.SolveOrQuant(messenger, quant, scopeSet, binding, continueAfterEnough, toMatch)

	} else if quant.Predicate == mentalese.PredicateAnd {

		result = base.SolveAndQuant(messenger, quant, scopeSet, binding, continueAfterEnough, toMatch)

	} else if quant.Predicate != mentalese.PredicateQuant {
		base.log.AddError("First argument of a `do` or `find` must contain only `quant`s")
		return mentalese.NewBindingSet()
	} else {

		result = base.solveSimpleQuant(messenger, quant, scopeSet, binding, continueAfterEnough, toMatch)

	}

	return result
}

func (base *SystemSolverFunctionBase) solveSimpleQuant(messenger api.ProcessMessenger, quant mentalese.Relation, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool, toMatch mentalese.TermList) mentalese.BindingSet {

	rangeVariable := quant.Arguments[mentalese.QuantRangeVariableIndex].TermValue
	rangeSet := quant.Arguments[mentalese.QuantRangeSetIndex].TermValueRelationSet

	setBindings := mentalese.NewBindingSet()
	if len(toMatch) > 0 {
		for _, element := range toMatch {
			b := binding.Copy()
			b.Set(rangeVariable, element)
			setBindings.Add(b)
		}
	} else {
		setBindings = mentalese.InitBindingSet(binding)
	}

	rangeBindings := messenger.ExecuteChildStackFrame(rangeSet, setBindings)
	scopeBindings := base.solveScope(messenger, quant, scopeSet, rangeBindings, continueAfterEnough)

	rangeCount := rangeBindings.GetDistinctValueCount(rangeVariable)
	scopeCount := scopeBindings.GetDistinctValueCount(rangeVariable)

	success := base.tryQuantifier(messenger, quant, rangeCount, scopeCount, true)

	if success {
		return scopeBindings
	} else {
		return mentalese.NewBindingSet()
	}
}

func (base *SystemSolverFunctionBase) SolveAndQuant(messenger api.ProcessMessenger, xorQuant mentalese.Relation, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool, toMatch mentalese.TermList) mentalese.BindingSet {

	leftQuant := xorQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
	rightQuant := xorQuant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]

	leftResultBindings := base.solveQuants(messenger, leftQuant, scopeSet, binding, continueAfterEnough, toMatch)

	resultBindings := mentalese.NewBindingSet()
	for _, leftResultBinding := range leftResultBindings.GetAll() {
		rightResultBindings := base.solveQuants(messenger, rightQuant, scopeSet, leftResultBinding, continueAfterEnough, toMatch)
		resultBindings.AddMultiple(rightResultBindings)
	}

	return resultBindings
}

func (base *SystemSolverFunctionBase) SolveOrQuant(messenger api.ProcessMessenger, orQuant mentalese.Relation, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool, toMatch mentalese.TermList) mentalese.BindingSet {
	leftQuant := orQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
	rightQuant := orQuant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]
	leftResultBindings := base.solveQuants(messenger, leftQuant, scopeSet, binding, continueAfterEnough, toMatch)
	rightResultBindings := base.solveQuants(messenger, rightQuant, scopeSet, binding, continueAfterEnough, toMatch)

	newBindings := leftResultBindings.Copy()
	newBindings.AddMultiple(rightResultBindings)
	return newBindings
}

func (base *SystemSolverFunctionBase) SolveXorQuant(messenger api.ProcessMessenger, xorQuant mentalese.Relation, scopeSet mentalese.RelationSet, binding mentalese.Binding, continueAfterEnough bool, toMatch mentalese.TermList) mentalese.BindingSet {
	leftQuant := xorQuant.Arguments[mentalese.SeqFirstOperandIndex].TermValueRelationSet[0]
	resultBindings := base.solveQuants(messenger, leftQuant, scopeSet, binding, continueAfterEnough, toMatch)
	if resultBindings.IsEmpty() {
		rightQuant := xorQuant.Arguments[mentalese.SeqSecondOperandIndex].TermValueRelationSet[0]
		resultBindings = base.solveQuants(messenger, rightQuant, scopeSet, binding, continueAfterEnough, toMatch)
	}

	return resultBindings
}

func (base *SystemSolverFunctionBase) solveScope(messenger api.ProcessMessenger, quant mentalese.Relation, scopeSet []mentalese.Relation, rangeBindings mentalese.BindingSet, continueAfterEnough bool) mentalese.BindingSet {

	rangeVariable := quant.Arguments[mentalese.QuantRangeVariableIndex].TermValue
	scopeBindings := mentalese.NewBindingSet()
	groupedScopeBindings := []mentalese.BindingSet{}

	for _, rangeBinding := range rangeBindings.GetAll() {

		singleScopeBindings := messenger.ExecuteChildStackFrame(scopeSet, mentalese.InitBindingSet(rangeBinding))

		if !singleScopeBindings.IsEmpty() {
			groupedScopeBindings = append(groupedScopeBindings, singleScopeBindings)
			scopeBindings.AddMultiple(singleScopeBindings)
		}

		rangeCount := rangeBindings.GetDistinctValueCount(rangeVariable)
		scopeCount := scopeBindings.GetDistinctValueCount(rangeVariable)

		ok := base.tryQuantifier(messenger, quant, rangeCount, scopeCount, false)
		if ok {
			if !continueAfterEnough {
				break
			}
		}
	}

	return scopeBindings
}

func (base *SystemSolverFunctionBase) tryQuantifier(messenger api.ProcessMessenger, quant mentalese.Relation, rangeCount int, scopeCount int, final bool) bool {

	firstArgument := quant.Arguments[mentalese.QuantQuantifierIndex]

	// special case: the existential quantifier `some`
	if firstArgument.IsAtom() && firstArgument.TermValue == mentalese.AtomSome {
		if scopeCount == 0 {
			if base.log.Active() {
				base.log.AddDebug("Do/Find", "Quantifier Some mismatch: no results")
			}
			return false
		} else {
			return true
		}
	}

	// special case: the existential quantifier `one`
	if firstArgument.IsAtom() && firstArgument.TermValue == mentalese.AtomOne {
		if scopeCount != 1 {
			if base.log.Active() {
				base.log.AddDebug("Do/Find", "Quantifier One mismatch: results not equal to 1")
			}
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

	b := mentalese.NewBinding()
	b.Set(rangeCountVariable, rangeVal)
	b.Set(scopeCountVariable, resultVal)

	quantifierBindings := messenger.ExecuteChildStackFrame(quantifierSet, mentalese.InitBindingSet(b))

	success := !quantifierBindings.IsEmpty()

	if !success {
		if base.log.Active() {
			base.log.AddDebug("Do/Find", "Quantifier mismatch")
			base.log.AddDebug("Do/Find", "  Range count: "+rangeCountVariable+" = "+strconv.Itoa(rangeCount))
			base.log.AddDebug("Do/Find", "  Scope count: "+scopeCountVariable+" = "+strconv.Itoa(scopeCount))
			base.log.AddDebug("Do/Find", "  Quantifier check: "+quantifierSet.String())
		}
	}

	return success
}

func (base *SystemSolverFunctionBase) quantMatch(messenger api.ProcessMessenger, quantMatch mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	if len(quantMatch.Arguments) != 2 {
		panic("quant match needs two arguments")
	}

	quant := quantMatch.Arguments[0].TermValueRelationSet[0]
	entities := mentalese.TermList{}

	if quantMatch.Arguments[1].IsVariable() {
		value, found := binding.Get(quantMatch.Arguments[1].TermValue)
		if found {
			entities = value.TermValueList
		} else {
			return mentalese.NewBindingSet()
		}
	} else {
		entities = quantMatch.Arguments[1].TermValueList
	}

	bindings := base.solveQuants(messenger, quant, mentalese.RelationSet{}, binding, true, entities)
	if bindings.IsEmpty() {
		return mentalese.NewBindingSet()
	} else {
		return mentalese.InitBindingSet(binding)
	}
}
