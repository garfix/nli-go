package mentalese

import "strconv"

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

// Removes all relations from set whose predicates match that of any of newSet
func (set RelationSet) RemoveMatchingPredicates(newSet RelationSet) RelationSet {

	resultSet := RelationSet{}

	for _, relation := range set {
		found := false
		for _, newRelation := range newSet {
			if relation.Predicate == newRelation.Predicate {
				found = true
				break
			}
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

	return set.BindRelationSetSingleBinding(importBinding)
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

func ResetVariables() {
	i = 0
}

func createVariable() Term {
	i++
	return NewVariable("X" + strconv.Itoa(i))
}








// Returns a new relation set, that has all variables bound to bindings
func (relations RelationSet) BindRelationSetSingleBinding(binding Binding) RelationSet {

	boundRelations := RelationSet{}

	for _, relation := range relations {
		boundRelations = append(boundRelations, relation.BindSingleRelationSingleBinding(binding))
	}

	return boundRelations
}

// Returns new relation sets, that have all variables bound to bindings
func (set RelationSet) BindRelationSetMultipleBindings(bindings []Binding) []RelationSet {

	boundRelationSets := []RelationSet{}

	for _, binding := range bindings {
		boundRelationSets = append(boundRelationSets, set.BindRelationSetSingleBinding(binding))
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

		if relation.Predicate == Predicate_Quant || relation.Predicate == Predicate_Quantification {
			// unscope the relation sets
			for i, argument := range relation.Arguments {
				if argument.IsRelationSet() {

					scopedSet := relationCopy.Arguments[i].TermValueRelationSet
					relationCopy.Arguments[i].TermValueRelationSet = RelationSet{}

					// recurse into the scope
					unscoped = append(unscoped, scopedSet.UnScope()...)
				}
			}
		}

		unscoped = append(unscoped, relationCopy)
	}

	return unscoped
}

//func (set RelationSet) UnmarshalJSON(b []byte) error {
//
//	var raw string
//
//	var parser importer.InternalGrammarParser
//
//	err := json.Unmarshal(b, &raw)
//	if err != nil {
//		return err
//	}
//
//	relationSet := parser.CreateRelationSet(raw)
//	parseResult := parser.GetLastParseResult()
//	if !parseResult.Ok {
//		return errors.New(parseResult.String())
//	}
//
//	for _, relation := range relationSet {
//		set = append(set, relation)
//	}
//
//	return nil
//}