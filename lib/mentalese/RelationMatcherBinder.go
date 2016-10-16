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

// Returns a new relation, that has all variables bound to bindings
func (matcher *RelationMatcher) BindSingleRelationSingleBinding(relation Relation, binding Binding) Relation {

	boundRelation := Relation{}
	boundRelation.Predicate = relation.Predicate

	for _, argument := range relation.Arguments {

		arg := argument
		if argument.IsVariable() {
			newValue, found := binding[argument.TermValue]
			if found {
				arg = newValue
			}
		}

		boundRelation.Arguments = append(boundRelation.Arguments, arg)
	}

	return boundRelation
}

// Returns a new relation set, that has all variables bound to bindings
func (matcher *RelationMatcher) BindRelationSetSingleBinding(relations RelationSet, binding Binding) RelationSet {

	boundRelations := RelationSet{}

	for _, relation:= range relations {
		boundRelations = append(boundRelations, matcher.BindSingleRelationSingleBinding(relation, binding))
	}

	return boundRelations
}

// Returns new relation sets, that have all variables bound to bindings
func (matcher *RelationMatcher) BindRelationSetMultipleBindings(relations RelationSet, bindings []Binding) []RelationSet {

	boundRelationSets := []RelationSet{}

	for _, binding := range bindings {
		boundRelationSets = append(boundRelationSets, matcher.BindRelationSetSingleBinding(relations, binding))
	}

	return boundRelationSets
}