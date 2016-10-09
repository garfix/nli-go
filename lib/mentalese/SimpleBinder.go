package mentalese


// Extends the binding with new variable bindings for the variables of subjectArgument
func (matcher *SimpleRelationMatcher) bindArgument(subjectArgument SimpleTerm, patternArgument SimpleTerm, binding SimpleBinding) (SimpleBinding, bool) {

	success := false

	if subjectArgument.IsAnonymousVariable() || patternArgument.IsAnonymousVariable() {

		// anonymous variables always match, but do not bind

		success = true

	} else if subjectArgument.IsVariable() {

		// variable

		value := SimpleTerm{}

		// does patternRelationArgument occur in boundVariables?
		value, match := binding[subjectArgument.String()]
		if match {
			// it does, use the bound variable
			if patternArgument.Equals(value) {
				success = true
			}
		} else {
			// it does not, just assign the actual argument
			binding[subjectArgument.String()] = patternArgument
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

	return binding, success
}

func (matcher *SimpleRelationMatcher) BindSingleRelationSingleBinding(relation SimpleRelation, binding SimpleBinding) SimpleRelation {

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

func (matcher *SimpleRelationMatcher) bindMultipleRelationsSingleBinding(relations SimpleRelationSet, binding SimpleBinding) SimpleRelationSet {

	for i, relation:= range relations {
		relations[i] = matcher.BindSingleRelationSingleBinding(relation, binding)
	}

	return relations
}

func (matcher *SimpleRelationMatcher) BindMultipleRelationsMultipleBindings(relations SimpleRelationSet, bindings []SimpleBinding) []SimpleRelationSet {

	relationSets := []SimpleRelationSet{}

	for _, binding := range bindings {
		relationSets = append(relationSets, matcher.bindMultipleRelationsSingleBinding(relations, binding))
	}

	return relationSets
}