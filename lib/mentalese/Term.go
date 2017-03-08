package mentalese

import "fmt"

type Term struct {
	TermType  int
	TermValue string
}

const Term_variable = 1
const Term_predicateAtom = 2
const Term_stringConstant = 3
const Term_number = 4
const Term_anonymousVariable = 5
const Term_regExp = 6

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

func (term Term) Equals(otherTerm Term) bool {
	return term.TermType == otherTerm.TermType && term.TermValue == otherTerm.TermValue
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
	default:
		s = "<unknown>"
	}
	return s
}
