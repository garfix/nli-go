package mentalese

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