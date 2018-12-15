package mentalese

// Extends the Binding with new variable bindings for the variables of subjectArgument
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

	} else if subjectArgument.IsRelationSet() {

		newBinding := binding.Copy()

		if patternArgument.IsVariable() {
			// [ isa(E, very) ], V
			success = true

		} else if patternArgument.IsRelationSet() {

			subSetBindingins, ok := matcher.MatchSequenceToSet(subjectArgument.TermValueRelationSet, patternArgument.TermValueRelationSet, newBinding)

			if ok {
				newBinding = subSetBindingins[0]
				success = true
			}
		}

		return newBinding, success

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
