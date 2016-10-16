package mentalese


// Extends the binding with new variable bindings for the variables of subjectArgument
func (matcher *RelationMatcher) BindTerm(subjectArgument Term, patternArgument Term, binding Binding) (Binding, bool) {

	success := false
	newBinding := Binding{}.Merge(binding)

	if subjectArgument.IsAnonymousVariable() || patternArgument.IsAnonymousVariable() {

		// anonymous variables always match, but do not bind

		success = true

	} else if subjectArgument.IsVariable() {

		// variable

		value := Term{}

		// does patternRelationArgument occur in boundVariables?
		value, match := newBinding[subjectArgument.String()]
		if match {
			// it does, use the bound variable
			if patternArgument.Equals(value) {
				success = true
			}
		} else {
			// it does not, just assign the actual argument
			newBinding[subjectArgument.String()] = patternArgument
			success = true
		}

	} else {

		// subject is atom, constant

		if patternArgument.IsVariable() {
			// note: no binding is made
			success = true
		} else if patternArgument.Equals(subjectArgument) {
			success = true
		}
	}

	return newBinding, success
}

func (matcher *RelationMatcher) BindSingleRelationSingleBinding(relation Relation, binding Binding) Relation {

	for i, argument := range relation.Arguments {

		if argument.IsVariable() {
			newValue, found := binding[argument.TermValue]
			if found {
				relation.Arguments[i] = newValue
			}
		}
	}

	return relation
}

func (matcher *RelationMatcher) BindMultipleRelationsSingleBinding(relations RelationSet, binding Binding) RelationSet {

	for i, relation:= range relations {
		relations[i] = matcher.BindSingleRelationSingleBinding(relation, binding)
	}

	return relations
}

func (matcher *RelationMatcher) BindMultipleRelationsMultipleBindings(relations RelationSet, bindings []Binding) []RelationSet {

	relationSets := []RelationSet{}

	for _, binding := range bindings {
		relationSets = append(relationSets, matcher.BindMultipleRelationsSingleBinding(relations, binding))
	}

	return relationSets
}