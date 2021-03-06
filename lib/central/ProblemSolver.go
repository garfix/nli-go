package central

import (
	"fmt"
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

// The problem solver takes a relation set and a set of bindings
// and returns a set of new bindings
// It uses knowledge bases to find these bindings
type ProblemSolver struct {
	index    			  *KnowledgeBaseIndex
	matcher               *RelationMatcher
	variableGenerator     *mentalese.VariableGenerator
	modifier              *FactBaseModifier
	dialogContext         *DialogContext
	log                   *common.SystemLog
}

func NewProblemSolver(matcher *RelationMatcher, dialogContext *DialogContext, log *common.SystemLog) *ProblemSolver {
	variableGenerator := mentalese.NewVariableGenerator()
	return &ProblemSolver{
		index: 			   NewProblemSolverIndex(),
		variableGenerator: variableGenerator,
		modifier:          NewFactBaseModifier(log, variableGenerator),
		matcher:           matcher,
		dialogContext:     dialogContext,
		log:               log,
	}
}

func (solver *ProblemSolver) AddFactBase(base api.FactBase) {
	solver.index.AddFactBase(base)
}

func (solver *ProblemSolver) AddFunctionBase(base api.FunctionBase) {
	solver.index.AddFunctionBase(base)
}

func (solver *ProblemSolver) AddRuleBase(base api.RuleBase) {
	solver.index.AddRuleBase(base)
}

func (solver *ProblemSolver) AddMultipleBindingBase(base api.MultiBindingBase) {
	solver.index.AddMultipleBindingBase(base)
}

func (solver *ProblemSolver) AddSolverFunctionBase(base api.SolverFunctionBase) {
	solver.index.AddSolverFunctionBase(base)
}

func (solver *ProblemSolver) ResetSession() {
	for _, factBase := range solver.index.factBases {
		switch v := factBase.(type) {
		case api.SessionBasedFactBase:
			v.ResetSession()
		}
	}
}

// set e.g. [ father(X, Y) father(Y, Z) ]
// bindings e.g. [{X: john, Z: jack} {}]
// return e.g. [
//  { X: john, Z: jack, Y: billy }
//  { X: john, Z: jack, Y: bob }
// ]
func (solver *ProblemSolver) SolveRelationSet(set mentalese.RelationSet, bindings mentalese.BindingSet) mentalese.BindingSet {

	newBindings := bindings
	for _, relation := range set {
		newBindings = solver.solveSingleRelationMultipleBindings(relation, newBindings)
		if newBindings.IsEmpty() {
			break
		}
	}

	return newBindings
}

// goal e.g. father(Y, Z)
// bindings e.g. {
//  { {X='john', Y='jack'} }
//  { {X='bob', Y='jonathan'} }
// }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver *ProblemSolver) solveSingleRelationMultipleBindings(relation mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {

	if solver.log.Active() { solver.log.StartDebug("Solve Relation", relation.String() + " " + fmt.Sprint(bindings)) }

	_, found := solver.index.knownPredicates[relation.Predicate]
	if !found {
		solver.log.AddError("Predicate not supported by any knowledge base: " + relation.Predicate)
		return mentalese.NewBindingSet()
	}

	newBindings, multiFound := solver.solveMultipleBindings(relation, bindings)

	if !multiFound {
		for _, binding := range bindings.GetAll() {
			newBindings.AddMultiple(solver.solveSingleRelationSingleBinding(relation, binding))
		}
	}

	if solver.log.Active() { solver.log.EndDebug("Solve Relation", relation.String() + ": " + fmt.Sprint(newBindings)) }

	return newBindings
}

func (solver *ProblemSolver) solveMultipleBindings(relation mentalese.Relation, bindings mentalese.BindingSet) (mentalese.BindingSet, bool) {

	newBindings := mentalese.NewBindingSet()
	multiFound := false

	functions, found := solver.index.multiBindingFunctions[relation.Predicate]
	if found {
		for _, function := range functions {
			newBindings = function(nil, relation, bindings)
			multiFound = true
		}
	}

	return newBindings, multiFound
}

// goalRelation e.g. father(Y, Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver *ProblemSolver) solveSingleRelationSingleBinding(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	relationVariables := relation.GetVariableNames()
	simpleBinding := binding.FilterVariablesByName(relationVariables)

	if solver.log.Active() { solver.log.StartDebug("Solve Simple Binding", relation.String() + " " + fmt.Sprint(simpleBinding)) }

	newBindings := mentalese.NewBindingSet()

	// go through all simple fact bases
	factBases, f4 := solver.index.factReadBases[relation.Predicate]
	if f4 {
		for _, factBase := range factBases {
			newBindings.AddMultiple(solver.FindFacts(factBase, relation, simpleBinding))
		}
	}

	// go through all rule bases
	bases, f3 := solver.index.ruleReadBases[relation.Predicate]
	if f3 {
		for _, base := range bases {
			newBindings.AddMultiple(solver.solveSingleRelationSingleBindingSingleRuleBase(relation, simpleBinding, base))
		}
	}

	// go through all simple function bases
	functions1, f1 := solver.index.simpleFunctions[relation.Predicate]
	if f1 {
		for _, function := range functions1 {
			resultBinding, success := function(relation, simpleBinding)
			if success {
				newBindings.Add(resultBinding)
			}
		}
	}

	// go through all solver functions
	functions2, f2 := solver.index.solverFunctions[relation.Predicate]
	if f2 {
		for _, function := range functions2 {
			newBindings.AddMultiple(function(nil, relation, simpleBinding))
		}
	}

	// do assert / retract
	newBindings.AddMultiple(solver.modifyKnowledgeBase(relation, simpleBinding))

	if solver.log.Active() { solver.log.EndDebug("Solve Simple Binding", relation.String() + ": " + fmt.Sprint(newBindings)) }

	// compose the result set
	completedBindings := mentalese.NewBindingSet()
	for _, newBinding := range newBindings.GetAll() {
		// remove temporary variables
		essentialResultBinding := newBinding.FilterVariablesByName(relationVariables)
		// combine the source binding with the clean results
		completedBinding := binding.Merge(essentialResultBinding)
		completedBindings.Add(completedBinding)
	}

	return completedBindings
}

// Creates bindings for the free variables in 'relations', by resolving them in factBase
func (solver *ProblemSolver) FindFacts(factBase api.FactBase, relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	dbBindings := mentalese.NewBindingSet()

	for _, ds2db := range factBase.GetReadMappings() {

		activeBinding, match := solver.matcher.MatchTwoRelations(relation, ds2db.Goal, mentalese.NewBinding())
		if !match { continue }

		activeBinding2, match2 := solver.matcher.MatchTwoRelations(ds2db.Goal, relation, mentalese.NewBinding())
		if !match2 { continue }

		dbRelations := ds2db.Pattern.ImportBinding(activeBinding2, solver.variableGenerator)

		localIdBinding := solver.replaceSharedIdsByLocalIds(binding, factBase)

		relevantBinding := localIdBinding.Select(dbRelations.GetVariableNames())
		newDbBindings := solver.solveMultipleRelationSingleFactBase(dbRelations, relevantBinding, factBase)

		for _, newDbBinding := range newDbBindings.GetAll() {

			dbBinding := activeBinding.Merge(newDbBinding)

			combinedBinding := localIdBinding.Merge(dbBinding.Select(relation.GetVariableNames()))
			sharedBinding := solver.replaceLocalIdBySharedId(combinedBinding, factBase)
			dbBindings.Add(sharedBinding)
		}
	}

	return dbBindings
}

func (solver *ProblemSolver) solveMultipleRelationSingleFactBase(relations []mentalese.Relation, binding mentalese.Binding, factBase api.FactBase) mentalese.BindingSet {

	sequenceBindings := mentalese.InitBindingSet(binding)

	for _, relation := range relations {
		sequenceBindings = solver.solveSingleRelationSingleFactBase(relation, sequenceBindings, factBase)
	}

	return sequenceBindings
}

func (solver *ProblemSolver) solveSingleRelationSingleFactBase(relation mentalese.Relation, bindings mentalese.BindingSet, factBase api.FactBase) mentalese.BindingSet {

	if solver.log.Active() { solver.log.StartDebug("Database" + " " + factBase.GetName(), relation.String() + " " + bindings.String()) }

	relationBindings, multiFound := solver.solveMultipleBindings(relation, bindings)
	resultBindings := mentalese.NewBindingSet()

	if !multiFound {

		for _, binding := range bindings.GetAll() {

			_, found := solver.index.simpleFunctions[relation.Predicate]
			if found {
				resultBindings = solver.solveSingleRelationSingleBinding(relation, binding)
			} else {
				resultBindings = factBase.MatchRelationToDatabase(relation, binding)
			}

			// found bindings must be extended with the bindings already present
			for _, resultBinding := range resultBindings.GetAll() {
				newRelationBinding := binding.Merge(resultBinding)
				relationBindings.Add(newRelationBinding)
			}
		}
	}

	if solver.log.Active() { solver.log.EndDebug("Database" + " " + factBase.GetName(), relationBindings.String()) }

	return relationBindings
}

func (solver *ProblemSolver) replaceSharedIdsByLocalIds(binding mentalese.Binding, factBase api.FactBase) mentalese.Binding {

	newBinding := mentalese.NewBinding()

	for key, value := range binding.GetAll() {
		newValue := value

		if value.IsId() {
			sharedId := value.TermValue
			sort := value.TermSort
			if sort != "" {
				localId := factBase.GetLocalId(sharedId, sort)
				if localId == "" {
					solver.log.AddError(fmt.Sprintf("Local id %s not found for %s in fact base %s", sharedId, sort, factBase.GetName()))
					return newBinding
				}
				newValue = mentalese.NewTermId(localId, sort)
			}
		}

		newBinding.Set(key, newValue)
	}

	return newBinding
}

func (solver *ProblemSolver) replaceLocalIdBySharedId(binding mentalese.Binding, factBase api.FactBase) mentalese.Binding {

	newBinding := mentalese.NewBinding()

	for key, value := range binding.GetAll() {
		newValue := value

		if value.IsId() {
			localId := value.TermValue
			sort := value.TermSort
			if sort != "" {
				sharedId := factBase.GetSharedId(localId, sort)
				if sharedId == "" {
					solver.log.AddError(fmt.Sprintf("Shared id %s not found for %s in fact base %s", localId, sort, factBase.GetName()))
					return newBinding
				}
				newValue = mentalese.NewTermId(sharedId, sort)
			}
		}

		newBinding.Set(key, newValue)
	}

	return newBinding
}

func (solver *ProblemSolver) modifyKnowledgeBase(relation mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	newBindings := mentalese.NewBindingSet()

	if len(relation.Arguments) == 0 { return newBindings }

	argument := relation.Arguments[0]

	if relation.Predicate == mentalese.PredicateAssert {
		if argument.IsRelationSet() {
			predicate := argument.TermValueRelationSet[0].Predicate
			factBases, found := solver.index.factWriteBases[predicate]
			if found {
				for _, factBase := range factBases {
					localIdBinding := solver.replaceSharedIdsByLocalIds(binding, factBase)
					boundRelation := relation.BindSingle(localIdBinding)
					singleRelation := boundRelation.Arguments[0].TermValueRelationSet[0]
					if singleRelation.IsBound() {
						solver.modifier.Assert(singleRelation, factBase)
					} else {
						solver.log.AddError("Cannot assert unbound relation " + singleRelation.String())
						return mentalese.NewBindingSet()
					}
					binding = solver.replaceLocalIdBySharedId(binding, factBase)
					newBindings.Add(binding)
				}
			} else {
				solver.log.AddError("Asserted relation not accepted by any fact base: " + predicate)
			}
		} else if argument.IsRule() {
			for _, ruleBase := range solver.index.ruleBases {
				rule := relation.Arguments[0].TermValueRule.BindSingle(binding)
				ruleBase.Assert(rule)
				solver.index.reindexRules()
				newBindings.Add(binding)
				//  only add the rule to a single rulebase
				break
			}
		} else if argument.IsList() {
			panic("assert not implemented for list")
		}
	} else if relation.Predicate == mentalese.PredicateRetract {
		if argument.IsRelationSet() {
			predicate := argument.TermValueRelationSet[0].Predicate
			factBases, found := solver.index.factWriteBases[predicate]
			if found {
				for _, factBase := range factBases {
					localIdBinding := solver.replaceSharedIdsByLocalIds(binding, factBase)
					boundRelation := relation.BindSingle(localIdBinding)
					solver.modifier.Retract(boundRelation.Arguments[0].TermValueRelationSet[0], factBase)
					binding = solver.replaceLocalIdBySharedId(binding, factBase)
					newBindings.Add(binding)
				}
			} else {
				solver.log.AddError("Retracted relation not accepted by any fact base: " + predicate)
			}
		}
	}

	return newBindings
}

// goalRelation e.g. father('jack', Z)
// binding e.g. { X='john', Y='jack' }
// return e.g. {
//  { {X='john', Y='jack', Z='joe'} }
//  { {X='bob', Y='jonathan', Z='bill'} }
// }
func (solver *ProblemSolver) solveSingleRelationSingleBindingSingleRuleBase(goalRelation mentalese.Relation, binding mentalese.Binding, ruleBase api.RuleBase) mentalese.BindingSet {

	inputVariables := goalRelation.GetVariableNames()

	goalBindings := mentalese.NewBindingSet()

	// match rules from the rule base to the goalRelation
	rules := ruleBase.GetRulesForRelation(goalRelation, binding)
	sourceSubgoalSets := []mentalese.RelationSet{}
	for _, rule := range rules {
		aBinding, _ := solver.matcher.MatchTwoRelations(goalRelation, rule.Goal, binding)
		bBinding, _ := solver.matcher.MatchTwoRelations(rule.Goal, goalRelation, mentalese.NewBinding())
		boundRule := rule.BindSingle(bBinding)
		boundRule = boundRule.InstantiateUnboundVariables(aBinding, solver.variableGenerator)
		sourceSubgoalSets = append(sourceSubgoalSets, boundRule.Pattern)
	}

	for _, sourceSubgoalSet := range sourceSubgoalSets {

		scopedBinding := mentalese.NewBinding().Merge(binding)
		subgoalResultBindings := mentalese.InitBindingSet(scopedBinding)

		for _, subGoal := range sourceSubgoalSet {

			subgoalResultBindings = solver.SolveRelationSet([]mentalese.Relation{subGoal}, subgoalResultBindings)
			if subgoalResultBindings.IsEmpty() {
				break
			}
		}

		for _, subgoalResultBinding := range subgoalResultBindings.GetAll() {

			// filter out the input variables
			filteredBinding := subgoalResultBinding.FilterVariablesByName(inputVariables)

			// make sure all variables of the original binding are present
			goalBinding := scopedBinding.Merge(filteredBinding)

			goalBindings.Add(goalBinding)
		}
	}

	return goalBindings
}
