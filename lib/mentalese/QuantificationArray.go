package mentalese

type QuantificationArray RelationSet

func (s QuantificationArray) Len() int {
	return len(s)
}

func (s QuantificationArray) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s QuantificationArray) Less(i, j int) bool {

	less := false

	first := s[i]
	firstQuantifier := first.Arguments[Quantification_QuantifierIndex]

	second := s[j]
	secondQuantifier := second.Arguments[Quantification_QuantifierIndex]

	firstQuantifierSimple := ""
	secondQuantifierSimple := ""

	// for now, we're just doing `all`, and `some`
	if len(firstQuantifier.TermValueRelationSet) == 1 {
		// isa(X, all)
		firstQuantifierSimple = firstQuantifier.TermValueRelationSet[0].Arguments[1].TermValue
	}
	if len(secondQuantifier.TermValueRelationSet) == 1 {
		secondQuantifierSimple = secondQuantifier.TermValueRelationSet[0].Arguments[1].TermValue
	}

	if firstQuantifierSimple == "all" && secondQuantifierSimple == "all" {
		less = i < j
	} else if firstQuantifierSimple == "all" {
		less = true
	} else if secondQuantifierSimple == "all" {
		less = false
	} else if false { // interrogative determiner
		less = true
	} else {
		// by default, reading order is order of preference
		less = i < j
	}

	return less
}
