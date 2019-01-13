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

	// turn new_quantifier into quantifier
	newRelations := scoper.fromQuantifierToTemp(relations)

	// collect all quantifications
	quantifications, nonQuantifications := scoper.collectQuantifications(newRelations)

	// sort the quantifications by hard constraints and preferences
	sort.Sort(QuantArray(quantifications))

	// nest the quantifications to create scopes
	scopedRelations := scoper.scopeQuants(quantifications)

	// fill in the other relations at the outermost position where they still are scoped.
	scoper.addNonQuantifications(&scopedRelations, len(quantifications), nonQuantifications)

	return scopedRelations
}

func (scoper QuantifierScoper) fromQuantifierToTemp(relations RelationSet) RelationSet {

	newRelations := RelationSet{}
	rangeRelations := RelationSet{}
	quantifierRelations := RelationSet{}

	workingSet := relations.Copy()

	for len(workingSet) > 0 {

		relation := workingSet[0]
		workingSet = workingSet[1:]

		newRelation := relation

		if relation.Predicate == Predicate_Quantification {

			quantificationVar := relation.Arguments[0]
			quantifierVar := relation.Arguments[1]
			rangeVar := relation.Arguments[2]

			rangeRelations, workingSet = scoper.extractRelationsWithVariable(workingSet, rangeVar.TermValue)
			quantifierRelations, workingSet = scoper.extractRelationsWithVariable(workingSet, quantifierVar.TermValue)

			workingSet = scoper.replaceVariable(workingSet, quantificationVar.TermValue, rangeVar.TermValue)

			newRelation = Relation{
				Predicate: Predicate_Quant,
				Arguments: []Term{
					relation.Arguments[2],
					NewRelationSet(rangeRelations),
					relation.Arguments[1],
					NewRelationSet(quantifierRelations),
					NewRelationSet(RelationSet{}),
				},
			}
		}

		newRelations = append(newRelations, newRelation)
	}

	return newRelations
}

func (scoper QuantifierScoper) extractRelationsWithVariable(relations RelationSet, variable string ) (RelationSet, RelationSet) {

	result := RelationSet{}
	remainder := RelationSet{}

	for _, relation := range relations {

		found := false

		for _, argument := range relation.Arguments {
			if argument.IsVariable() && argument.TermValue == variable {
				found = true
			}
		}

		if found {
			result = append(result, relation)
		} else {
			remainder = append(remainder, relation)
		}
	}

	return result, remainder
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

func (scoper QuantifierScoper) collectQuantifications(relations RelationSet) (QuantArray, RelationSet) {
	quantifications := QuantArray{}
	nonQuantifications := RelationSet{}
	for _, relation := range relations {
		if relation.Predicate == Predicate_Quant {
			quantifications = append(quantifications, relation)
		} else {
			nonQuantifications = append(nonQuantifications, relation)
		}
	}
	return quantifications, nonQuantifications
}

func (scoper QuantifierScoper) scopeQuants(quants QuantArray) RelationSet {

	scope := RelationSet{}

	for i := len(quants) - 1; i >= 0; i-- {

		quant := quants[i]
		quant.Arguments[Quantification_ScopeIndex] = NewRelationSet(scope)

		scope = RelationSet{quant}
	}

	return scope
}

func (scoper QuantifierScoper) addNonQuantifications(scopedRelations *RelationSet, depth int, nonQuantifications RelationSet) {

	for _, nonQuantification := range nonQuantifications {

		scope := scopedRelations
		nonQuantificationScope := scope

		for d := 0; d < depth; d++ {

			scopedRelation := (*scope)[0]
			rangeVariable := scopedRelation.Arguments[Quantification_RangeVariableIndex]

			scope = &scopedRelation.Arguments[Quantification_ScopeIndex].TermValueRelationSet

			if scoper.variableMatches(nonQuantification, rangeVariable) {
				nonQuantificationScope = scope
			} else if scoper.someRelationVariableMatches(nonQuantification, scope) {
				nonQuantificationScope = scope
			}
		}

		*nonQuantificationScope = append(*nonQuantificationScope, nonQuantification)
	}
}

func (scoper QuantifierScoper) variableMatches(relation Relation, variable Term) bool {
	match := false

	for _, argument := range relation.Arguments {
		if argument.Equals(variable) {
			match = true
			break
		}
	}

	return match
}

// if range variable is R5 and scope is
//		have_child(R5, E6)
// and needle is
//      number_of(E6, 4)
// then it should be added to the scope
func (scoper QuantifierScoper) someRelationVariableMatches(needle Relation, hayStack *RelationSet) bool {
	match := false

	for _, argument1 := range needle.Arguments {
		if argument1.IsVariable() {
			for _, straw := range *hayStack {
				for _, argument2 := range straw.Arguments {
					if argument2.IsVariable() && argument2.TermValue == argument1.TermValue {
						match = true
						break
					}
				}
			}
		}
	}

	return match
}
