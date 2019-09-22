package mentalese

import "strconv"

// An array of relations
type RelationSet []Relation

func (set RelationSet) Copy() RelationSet {

	copiedSet := RelationSet{}

	for _, relation := range set {
		copiedSet = append(copiedSet, relation.Copy())
	}

	return copiedSet
}

func (set RelationSet) IsEmpty() bool {
	return len(set) == 0
}

func (set RelationSet) Equals(newSet RelationSet) bool {

	if len(set) != len(newSet) {
		return false
	}

	for i, newRelation := range newSet {
		if !newRelation.Equals(set[i]) {
			return false
		}
	}

	return true
}

func (set RelationSet) Merge(newSet RelationSet) RelationSet {

	mergedSet := set.Copy()

	for _, newRelation := range newSet {

		found := false

		for _, existingRelation := range mergedSet {
			if newRelation.Equals(existingRelation) {
				found = true
			}
		}

		if !found {
			mergedSet = append(mergedSet, newRelation)
		}
	}

	return mergedSet
}

func (set RelationSet) Contains(needle Relation) bool {

	for _, relation := range set {
		if relation.Equals(needle) {
			return true
		}
	}

	return false
}

func (set RelationSet) RemoveDuplicates() RelationSet {

	resultSet := RelationSet{}

	for _, relation := range set {
		if !resultSet.Contains(relation) {
			resultSet = append(resultSet, relation)
		}
	}

	return resultSet
}

func (set RelationSet) RemoveRelations(remove RelationSet) RelationSet {

	resultSet := RelationSet{}

	for _, relation := range set {
		if !remove.Contains(relation) {
			resultSet = append(resultSet, relation)
		}
	}

	return resultSet
}

// Recursively removes all relations from set whose predicates match predicate
func (set RelationSet) RemoveMatchingPredicate(predicate string) RelationSet {

	resultSet := RelationSet{}

	for _, relation := range set {
		found := false

		for a, relationArgument := range relation.Arguments {
			if relationArgument.IsRelationSet() {
				relation = relation.Copy()
				relation.Arguments[a].TermValueRelationSet = relationArgument.TermValueRelationSet.RemoveMatchingPredicate(predicate)
			}
		}

		if relation.Predicate == predicate {
			found = true
		}

		if !found {
			resultSet = append(resultSet, relation)
		}
	}

	return resultSet
}

// set contains variables A, B, and C
// binding: A: X, C: Z
// resulting set contains X (for A), Z (for C) and X22 (a new variable)
func (set RelationSet) ImportBinding(binding Binding) RelationSet {

	importBinding := binding.Copy()

	// find all variables
	variables := findVariables(set)

	// extend binding with extra variables in set
	for _, variable  := range variables {
		_, found := importBinding[variable]
		if !found {
			importBinding[variable] = createVariable()
		}
	}

	// replace variables in set

	return set.BindSingle(importBinding)
}

func findVariables(set RelationSet) []string {

	var variables []string

	for _, relation := range set {
		for _, argument := range relation.Arguments {
			if argument.IsVariable() {
				variables = append(variables, argument.TermValue)
			} else if argument.IsRelationSet() {
				variables = append(variables, findVariables(argument.TermValueRelationSet)...)
			}
		}
	}

	return variables
}

var i = 0

func createVariable() Term {
	i++
	return NewVariable("Gen" + strconv.Itoa(i))
}








// Returns a new relation set, that has all variables bound to bindings
func (relations RelationSet) BindSingle(binding Binding) RelationSet {

	boundRelations := RelationSet{}

	for _, relation := range relations {
		boundRelations = append(boundRelations, relation.BindSingleRelationSingleBinding(binding))
	}

	return boundRelations
}

// Returns new relation sets, that have all variables bound to bindings
func (set RelationSet) BindRelationSetMultipleBindings(bindings Bindings) []RelationSet {

	boundRelationSets := []RelationSet{}

	for _, binding := range bindings {
		boundRelationSets = append(boundRelationSets, set.BindSingle(binding))
	}

	return boundRelationSets
}

func (set RelationSet) String() string {

	s, sep := "", ""

	for _, relation := range set {
		s += sep + relation.String()
		sep = " "
	}

	return "[" + s + "]"
}

func (set RelationSet) UnScope() RelationSet {

	unscoped := RelationSet{}

	for _, relation := range set {

		relationCopy := relation.Copy()

		// unscope the relation sets
		for i, argument := range relation.Arguments {
			if argument.IsRelationSet() {

				scopedSet := relationCopy.Arguments[i].TermValueRelationSet
				relationCopy.Arguments[i].TermValueRelationSet = RelationSet{}

				// recurse into the scope
				unscoped = append(unscoped, scopedSet.UnScope()...)
			}
		}

		unscoped = append(unscoped, relationCopy)
	}

	return unscoped
}

// Returns all relations with variable as argument; those relations have other variables, find all relations with those as well
func (set RelationSet) findRelationsStartingWithVariable(variable string) RelationSet {

	foundVariables := map[string]bool{ variable: true }
	markedRelationIndexes := map[int]bool{}
	awaitingVariables := []string{ variable }

	// process a stack with variables
	for len(awaitingVariables) != 0 {
		// pop
		var activeVariable = awaitingVariables[0]
		awaitingVariables = awaitingVariables[1:]

		// check all relations for this variable
		for r, relation := range set {

			// quick skip
			_, processedBefore := markedRelationIndexes[r]
			if processedBefore {
				continue
			}

			if relation.RelationUsesVariable(activeVariable) {
				// mark this relation
				markedRelationIndexes[r] = true
				// add all variables of this relation
				for _, argument := range relation.Arguments {
					if argument.TermType == TermVariable {
						var someVar = argument.TermValue
						var _, variableAlreadyFound = foundVariables[someVar]
						if !variableAlreadyFound {
							foundVariables[someVar] = true
							awaitingVariables = append(awaitingVariables, someVar)
						}
					}
				}
			}
		}
	}

	// create result set
	var resultSet RelationSet
	for i := range markedRelationIndexes {
		resultSet = append(resultSet, set[i])
	}

	return resultSet
}

func (set RelationSet) ConvertVariablesToConstants() RelationSet {
	newSet := RelationSet{}

	for _, relation := range set {
		newSet = append(newSet, relation.ConvertVariablesToConstants())
	}

	return newSet
}