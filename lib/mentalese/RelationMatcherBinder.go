package mentalese


// Extends the binding with new variable bindings for the variables of subjectArgument
func (matcher *RelationMatcher) BindTerm(subjectArgument Term, patternArgument Term, binding Binding) (Binding, bool) {

	success := false

	if subjectArgument.IsAnonymousVariable() || patternArgument.IsAnonymousVariable() {

		// anonymous variables always match, but do not bind

		// A, _
		// _, A
		return binding, true

	} else if subjectArgument.IsVariable() {

		// subject is variable

		value, match := binding[subjectArgument.String()]
		if match {

			// A, 13 {A:13}
			if patternArgument.Equals(value) {
				success = true
			}
			return binding, success

		} else {

			// A, 13, {B:7} => {B:7, A:13}
			newBinding := binding.Copy()
			newBinding[subjectArgument.String()] = patternArgument
			return newBinding, true
		}

	} else {

		// subject is atom, constant

		if patternArgument.IsVariable() {
			// 13, V
			success = true
		} else if patternArgument.Equals(subjectArgument) {
			// 13, 13
			// female, female
			// 'Jack', 'Jack'
			success = true
		}

		return binding, success
	}
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