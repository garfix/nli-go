package mentalese

type RelationSet []Relation

func (set RelationSet) Copy() RelationSet {

	copiedSet := RelationSet{}

	for _, relation := range set {
		copiedSet = append(copiedSet, relation)
	}

	return copiedSet
}

func (set RelationSet) Merge(newSet RelationSet) RelationSet {

	mergedSet := set.Copy()

	for _, newRelation := range newSet {

		found := false

		for _, existingRelation := range set {
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

func (set RelationSet) String() string {

	s, sep := "", ""

	for _, relation := range set {
		s += sep + relation.String()
		sep = " "
	}

	return "[" + s + "]";
}