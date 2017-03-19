package mentalese

import "fmt"

type Term struct {
	TermType  int
	TermValue string
	TermValueRelationSet RelationSet
}

const Term_variable = 1
const Term_predicateAtom = 2
const Term_stringConstant = 3
const Term_number = 4
const Term_anonymousVariable = 5
const Term_regExp = 6
const Term_relationSet = 7

func (term Term) IsVariable() bool {
	return term.TermType == Term_variable
}

func (term Term) IsNumber() bool {
	return term.TermType == Term_number
}

func (term Term) IsRegExp() bool {
	return term.TermType == Term_regExp
}

func (term Term) IsAnonymousVariable() bool {
	return term.TermType == Term_anonymousVariable
}

func (term Term) IsRelationSet() bool {
	return term.TermType == Term_relationSet
}

func (term Term) Equals(otherTerm Term) bool {
	if term.TermType != otherTerm.TermType {
		return false
	}
	if term.TermType == Term_relationSet {
		return term.TermValueRelationSet.Equals(otherTerm.TermValueRelationSet)
	} else {
		return term.TermValue == otherTerm.TermValue
	}
}

func (term Term) AsKey() string {
	return fmt.Sprintf("%d/%s", term.TermType, term.TermValue)
}

func (term Term) String() string {

	s := ""

	switch term.TermType {
	case Term_variable:
		s = term.TermValue
	case Term_predicateAtom:
		s = term.TermValue
	case Term_stringConstant:
		s = "'" + term.TermValue + "'"
	case Term_regExp:
		s = "/" + term.TermValue + "/"
	case Term_number:
		s = term.TermValue
	case Term_anonymousVariable:
		s = "_"
	case Term_relationSet:
		s = term.TermValueRelationSet.String()
	default:
		s = "<unknown>"
	}
	return s
}
