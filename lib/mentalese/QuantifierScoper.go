package mentalese

import (
	"nli-go/lib/common"
	"sort"
)

// Turns a set of unscoped (or partially scoped) relations into scoped relations
// Like this
//
//     quant(S1, [ isa(S1, parent) ], D1, [ isa(D1, every) ], [
//         quant(O1, [ isa(O1, child) ], D2, [ isa(D2, 2) ], [
//             have_child(S1, O1)
//         ])
//     ])
//
// This part was mainly influenced by quantifier scoping of the Core Language Engine
// Specifically the article "Quantifier Scoping in the SRI Core Language Engine" by Douglas B. Moran

type QuantifierScoper struct {
	log *common.SystemLog
}

func NewQuantifierScoper(log *common.SystemLog) QuantifierScoper {
	return QuantifierScoper{log: log}
}

func (scoper QuantifierScoper) Scope(relations RelationSet) RelationSet {

	// separate quants from non-quants
	quants, nonQuants := scoper.collectQuants(relations)

	// add quantifiers
	quants, nonQuants = scoper.addQuantifiersAndRanges(quants, nonQuants)

	// sort the quants by hard constraints and preferences
	sort.Sort(QuantArray(quants))

	// nest the quants to create scopes
	scopedRelations := scoper.scopeQuants(quants, nonQuants)

	return scopedRelations
}

func (scoper QuantifierScoper) addQuantifiersAndRanges(quants QuantArray, nonQuants RelationSet) (QuantArray, RelationSet) {

	newQuants := QuantArray{}
	newNonQuants := nonQuants

	for _, quant := range quants {

		quantifierVar := quant.Arguments[QuantQuantifierVariableIndex].TermValue
		quantifier := newNonQuants.findRelationsStartingWithVariable(quantifierVar)
		newNonQuants = newNonQuants.RemoveRelations(quantifier)

		rangeVar := quant.Arguments[QuantRangeVariableIndex].TermValue
		rangeSet := newNonQuants.findRelationsStartingWithVariable(rangeVar)
		newNonQuants = newNonQuants.RemoveRelations(rangeSet)

		quant.Arguments[QuantQuantifierIndex] = NewRelationSet(quantifier)
		quant.Arguments[QuantRangeIndex] = NewRelationSet(rangeSet)

		newQuants = append(newQuants, quant)
	}

	return newQuants, newNonQuants
}

func (scoper QuantifierScoper) collectQuants(relations RelationSet) (QuantArray, RelationSet) {
	quants := QuantArray{}
	nonQuants := RelationSet{}
	for _, relation := range relations {
		if relation.Predicate == PredicateQuant {
			quants = append(quants, relation)
		} else {
			nonQuants = append(nonQuants, relation)
		}
	}
	return quants, nonQuants
}

func (scoper QuantifierScoper) scopeQuants(quants QuantArray, nonQuants RelationSet) RelationSet {

	newQuants, newNonQuants := scoper.scopeFirstQuant(quants, nonQuants)

	return append(newNonQuants, newQuants...)
}

func (scoper QuantifierScoper) scopeFirstQuant(quants QuantArray, nonQuants RelationSet) (QuantArray, RelationSet) {

	combineScopedSet := RelationSet{}

	if len(quants) == 0 {
		return quants, nonQuants
	}

	if len(quants) > 1 {
		// scope the rest of the quants first
		restQuants, restNonQuants := scoper.scopeFirstQuant(quants[1:], nonQuants)

		combineScopedSet = append(combineScopedSet, restQuants...)
		nonQuants = restNonQuants
	}

	quant := quants[0]

//	rangeVariable := quant.Arguments[QuantRangeVariableIndex].TermValue
	scopeVariable := quant.Arguments[QuantScopeIndex + 1].TermValue

	//rangeSet := nonQuants.findRelationsStartingWithVariable(rangeVariable)
	//nonQuants = nonQuants.RemoveRelations(rangeSet)
	scopeSet := nonQuants.findRelationsStartingWithVariable(scopeVariable)
	nonQuants = nonQuants.RemoveRelations(scopeSet)

	combineScopedSet = append(combineScopedSet, scopeSet...)

//	quant.Arguments[QuantRangeIndex] = NewRelationSet(rangeSet)
	quant.Arguments[QuantScopeIndex] = NewRelationSet(combineScopedSet)

	return QuantArray{quant}, nonQuants
}

func (scoper QuantifierScoper) replaceVariables(relations RelationSet, rangeScopeMap map[string]string) RelationSet {

	newRelations := RelationSet{}

	for _, relation := range relations {

		newRelation := relation

		if relation.Predicate == PredicateQuant {

			newScopeSet := scoper.replaceVariables(relation.Arguments[QuantScopeIndex].TermValueRelationSet, rangeScopeMap)

			for rangeVar, scopeVar := range rangeScopeMap {
				newScopeSet = scoper.replaceVariable(newScopeSet, scopeVar, rangeVar)
			}

			newRelation = Relation{
				Predicate: PredicateQuant,
				Arguments: []Term{
					relation.Arguments[QuantQuantifierVariableIndex],
					relation.Arguments[QuantQuantifierIndex],
					relation.Arguments[QuantRangeVariableIndex],
					relation.Arguments[QuantRangeIndex],
					NewRelationSet(newScopeSet),
				},
			}
		}

		newRelations = append(newRelations, newRelation)
	}

	return newRelations
}

func (scoper QuantifierScoper) replaceVariable(relations RelationSet, oldVar string, newVar string) RelationSet {

	result := RelationSet{}

	for _, relation := range relations {

		newRelation := relation.Copy()

		for i, argument := range newRelation.Arguments {
			if argument.IsVariable() && argument.TermValue == oldVar {
				newRelation.Arguments[i].TermValue = newVar
			}
		}

		result = append(result, newRelation)
	}

	return result
}
